package httpbench

import (
	"context"
	"log"
	"time"

	"github.com/miekg/dns"
)

//Get the list of possible IPs a hostname resolves to using specific DNS server
func resolve(hostname, ip string) []string {
	ips := make([]string, 0)

	for i := 0; i < 5; i++ {
		m := new(dns.Msg)
		m.Id = dns.Id()
		m.RecursionDesired = true
		m.SetQuestion(dns.Fqdn(hostname), dns.TypeA)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		c := new(dns.Client)
		in, _, err := c.ExchangeContext(ctx, m, ip+":53")
		if err == nil {
			if len(in.Answer) < 1 {
				log.Println("Error no answer")
			} else {
				for _, rec := range in.Answer {
					arec, ok := rec.(*dns.A)
					if ok {
						ips = appendIfMissing(ips, arec.A.String())
					}
				}
			}
		}
	}
	return ips
}
