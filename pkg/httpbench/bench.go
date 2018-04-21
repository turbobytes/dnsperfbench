package httpbench

import (
	"log"
	"net/url"
	"sort"
)

//Result of individual server
type Result struct {
	Server string
	CI     *ConInfo
}

//Results list of Result, to make it sortable
type Results []Result

func (a Results) Len() int           { return len(a) }
func (a Results) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Results) Less(i, j int) bool { return a[i].CI.Total < a[j].CI.Total }

//TestOverHTTP tests a url over HTTP after resolving against various resolvers
func TestOverHTTP(u *url.URL, resolvers []string) Results {
	var finalresult Results = make([]Result, 0)
	for _, resolver := range resolvers {
		results := make([]*ConInfo, 0)
		//Issue 10 tests
		for i := 0; i < 10; i++ {
			ci, err := testoverhttp(u, resolver)
			//log.Println(ip, ci, err)
			if err != nil {
				//ipresall[ip] = append(ipresall[ip], nil)
				log.Printf("%v, Using Resolver: %v\n", err, resolver)
			} else {
				results = append(results, ci)
				//ipresall[ip] = append(ipresall[ip], ci)
			}
		}
		if len(results) > 2 {
			//Atleast 3 tests worked... find median...
			finalresult = append(finalresult, Result{
				Server: resolver,
				CI:     medianconinfo(results),
			})

		}
	}
	sort.Sort(finalresult)
	return finalresult
}
