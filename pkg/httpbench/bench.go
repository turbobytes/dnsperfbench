package httpbench

import (
	"log"
	"net/url"
	"sort"
)

func appendIfMissing(src []string, new string) []string {
	for _, ele := range src {
		if ele == new {
			return src
		}
	}
	return append(src, new)
}

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
	ipmap := make(map[string]string)
	allipmap := make(map[string][]string)
	log.Println("Resolving")
	for _, server := range resolvers {
		ips := resolve(u.Hostname(), server)
		//Map popularity of each ip
		for _, ip := range ips {
			allipmap[ip] = append(allipmap[ip], server)
		}
	}
	//Popularity contest
	for i := len(resolvers); i != 0; i-- {
		for ip, servers := range allipmap {
			if len(servers) == i {
				for _, server := range servers {
					if ipmap[server] == "" {
						ipmap[server] = ip
					}
				}
			}
		}
	}
	iplist := make([]string, 0)
	for _, v := range ipmap {
		iplist = appendIfMissing(iplist, v)
	}
	ipres := make(map[string]*ConInfo)
	log.Println("Issuing HTTP(s) tests")
	for _, ip := range iplist {
		results := make([]*ConInfo, 10)
		for i := 0; i < 10; i++ {
			ci, err := testoverhttp(u, ip)
			//log.Println(ip, ci, err)
			if err != nil {
				//ipresall[ip] = append(ipresall[ip], nil)
				log.Printf("%v, Using IP: %v, Reported by %v\n", err, ip, allipmap[ip])
				results[i] = nil
			} else {
				results[i] = ci
				//ipresall[ip] = append(ipresall[ip], ci)
			}
		}
		ipres[ip] = medianconinfo(results)
	}
	//Compute "median" run per ip
	var finalresult Results = make([]Result, 0)
	for server, ip := range ipmap {
		if ipres[ip] != nil {
			//finalresult[server] = ipres[ip]
			finalresult = append(finalresult, Result{
				Server: server,
				CI:     ipres[ip],
			})
		}
	}
	sort.Sort(finalresult)
	return finalresult
}
