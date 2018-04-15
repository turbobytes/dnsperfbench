package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var resolvers arrayFlags

var (
	raw              = flag.Bool("r", false, "Output raw mode")
	defaultResolvers = []string{"8.8.8.8", "1.1.1.1", "9.9.9.9", "208.67.222.222", "199.85.126.20", "185.228.168.168", "8.26.56.26"}
	resolverNames    = map[string]string{
		"8.8.8.8":         "Google",
		"1.1.1.1":         "Cloudflare",
		"9.9.9.9":         "Quad9",
		"208.67.222.222":  "OpenDNS",
		"199.85.126.20":   "Norton",
		"185.228.168.168": "Clean Browsing",
		"8.26.56.26":      "Comodo",
	}
	//All answers must match these
	expectedanswers = map[string]struct{}{
		"138.197.54.54": struct{}{},
		"138.197.53.4":  struct{}{},
	}
	//Duration to signal fail
	failDuration = time.Second * 10
	hostnamesHIT = []string{"fixed.turbobytes.net.", "fixed2.turbobytes.net."}
	auths        = map[string]string{
		"NS1":         "tbrum3.com.",
		"Google":      "tbrum4.com.",
		"AWS Route53": "tbrum5.com.",
		"DNSimple":    "tbrum14.com.",
		"GoDaddy":     "tbrum2.com.",
		"Akamai":      "tbrum9.com.",
		"Dyn":         "tbrum10.com.",
		"CloudFlare":  "tbrum8.com.",
		"EasyDNS":     "tbrum16.com.",
		"Ultradns":    "tbrum22.com.",
		"Azure":       "tbrum25.com.",
	}
	authSl    []string
	ratelimit = make(chan struct{}, 5) //Max number of dns queries in flight at a time
)

const (
	testrep = 15 //Number of times to repeat each test
)

func appendIfMissing(src []string, new string) []string {
	for _, ele := range src {
		if ele == new {
			return src
		}
	}
	return append(src, new)
}

func init() {
	var tmp arrayFlags
	flag.Var(&tmp, "resolver", "Additional resolvers to test. default="+strings.Join(defaultResolvers, ", "))
	flag.Parse()
	resolvers = defaultResolvers
	for _, res := range tmp {
		resolvers = appendIfMissing(resolvers, res)
	}
	rand.Seed(time.Now().Unix())
	authSl = make([]string, 0)
	for auth := range auths {
		authSl = append(authSl, auth)
	}
	sort.Strings(authSl)

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func testresolver(hostname, resolver string) *time.Duration {
	//Add to ratelimit, block until a slot is available
	ratelimit <- struct{}{}
	//Remove from rate limit when done
	defer func() { <-ratelimit }()
	m := new(dns.Msg)
	m.Id = dns.Id()
	m.RecursionDesired = true
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeA)
	c := new(dns.Client)
	//Life is too short to wait for DNS...
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in, rtt, err := c.ExchangeContext(ctx, m, resolver+":53")
	if err != nil {
		return nil
	}
	//Validate response
	//Expect only one answer
	if len(in.Answer) != 1 {
		return nil
	}
	arec, ok := in.Answer[0].(*dns.A)
	if !ok {
		return nil
	}
	_, ok = expectedanswers[arec.A.String()]
	if !ok {
		return nil
	}
	//rtt = rtt.Truncate(time.Millisecond / 4)
	return &rtt
}

func runtests(host, res string, rndSuffix bool) resolverResults {
	//Actual test...
	vals := make([]time.Duration, 0)
	fails := 0
	for i := 0; i < testrep; i++ {
		hostname := host
		if rndSuffix {
			hostname = randStringRunes(15) + "." + host
		}
		rtt := testresolver(hostname, res)
		if rtt == nil {
			fails++
		} else {
			vals = append(vals, *rtt)
		}
	}
	//Print summary
	//fmt.Printf("Failures: %v of 5\n", fails)
	//fmt.Printf("Timings: %v\n", vals)
	validVals := make([]float64, len(vals))
	for i, val := range vals {
		validVals[i] = float64(val)
	}
	median, _ := stats.Median(validVals)
	mean, _ := stats.Mean(validVals)
	return resolverResults{mean: time.Duration(mean), median: time.Duration(median), failratio: float64(fails) / testrep}
}

//SummaryResolver stores score for individual resolver
type SummaryResolver struct {
	Res   string
	Score float64
}

//Summary enables sorting slice of SummaryResolver
type Summary []SummaryResolver

func (a Summary) Len() int           { return len(a) }
func (a Summary) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Summary) Less(i, j int) bool { return a[i].Score < a[j].Score }

type resolverResults struct {
	mean      time.Duration
	median    time.Duration
	failratio float64
}

type recursiveResults map[string]resolverResults

func getms(dur time.Duration) float64 {
	return float64(dur) / float64(time.Millisecond)
}

func (res recursiveResults) Print(resolver, name string) {
	if *raw {
		result := res["ResolverHit"]
		fmt.Printf("Raw\t%s\tResolverHit\t%.2f\t%.2f\t%.2f\n", resolver, getms(result.mean), getms(result.median), result.failratio*100)
		for _, auth := range authSl {
			result := res[auth]
			fmt.Printf("Raw\t%s\t%s\t%.2f\t%.2f\t%.2f\n", resolver, auth, getms(result.mean), getms(result.median), result.failratio*100)
		}
	} else {
		fmt.Printf("========== %s (%s) ===========\n", resolver, name)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Auth", "Mean", "Median", "Fail"})
		result := res["ResolverHit"]
		table.Append([]string{"ResolverHit", result.mean.Round(time.Millisecond).String(), result.median.Round(time.Millisecond).String(), fmt.Sprintf("%.2f%%", result.failratio*100)})
		for _, auth := range authSl {
			result := res[auth]
			table.Append([]string{auth, result.mean.Round(time.Millisecond).String(), result.median.Round(time.Millisecond).String(), fmt.Sprintf("%.2f%%", result.failratio*100)})
		}
		table.Render()
	}
}

func (res recursiveResults) Score() float64 {
	result := res["ResolverHit"]
	score := 5 * (float64(result.mean/time.Millisecond) + float64(result.median/time.Millisecond) + result.failratio*testrep*float64(failDuration/time.Millisecond))
	for _, auth := range authSl {
		result := res[auth]
		score += float64(result.mean/time.Millisecond) + float64(result.median/time.Millisecond) + result.failratio*testrep*float64(failDuration/time.Millisecond)
	}
	return score
}

func testrecursive(res string) recursiveResults {
	results := make(map[string]resolverResults)
	hithost := hostnamesHIT[rand.Intn(len(hostnamesHIT))]
	//Prime the caches... ignoring results
	for i := 0; i < 5; i++ {
		testresolver(hithost, res)
	}
	results["ResolverHit"] = runtests(hithost, res, false)

	//Perform the auths
	for _, auth := range authSl {
		host := auths[auth]
		results[auth] = runtests(host, res, true)
	}
	return results
}

type resultoutput struct {
	recursive string
	result    recursiveResults
}

func main() {
	resscore := make(map[string]float64)
	results := make(map[string]recursiveResults)
	resultschan := make(chan resultoutput, 1)
	//Fire off tests
	for _, res := range resolvers {
		//Stagger the start of tests
		time.Sleep(time.Millisecond * 50)
		log.Println("Issuing tests for ", res)
		go func(recursive string) {
			resultschan <- resultoutput{recursive: recursive, result: testrecursive(recursive)}
		}(res)
	}
	//Gather results
	log.Println("Waiting to gather results ")
	for i := range resolvers {
		result := <-resultschan
		log.Printf("[%v/%v] Got results for %s\n", i+1, len(resolvers), result.recursive)
		results[result.recursive] = result.result
	}
	for _, res := range resolvers {
		name := resolverNames[res]
		if name == "" {
			name = "Unknown"
		}
		result := results[res]
		result.Print(res, name)
		resscore[res] = result.Score()
	}
	//Make slice
	var summary Summary = make([]SummaryResolver, 0)
	for k, v := range resscore {
		summary = append(summary, SummaryResolver{k, v})
	}
	sort.Sort(summary)
	if !*raw {
		fmt.Printf("========== Summary ===========\n")
		fmt.Println("Scores (lower is better)")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Resolver", "Performance Score"})

	for _, sum := range summary {
		name := resolverNames[sum.Res]
		if name == "" {
			name = "Unknown"
		}
		table.Append([]string{fmt.Sprintf("%s (%s)", sum.Res, name), fmt.Sprintf("%.0f", sum.Score)})
		if *raw {
			fmt.Printf("Score\t%s\t%.0f\n", sum.Res, sum.Score)
		}
		//log.Println(sum.Res, sum.Score)
	}
	if *raw {
		fmt.Printf("Recommendation\t%s\n", summary[0].Res)
	} else {
		table.Render()
		fmt.Printf("You should probably use %s as your default resolver\n", summary[0].Res)
	}
}
