package httpbench

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"github.com/montanaflynn/stats"
)

//ConInfo holds timing data for http test
type ConInfo struct {
	DNS      time.Duration
	Connect  time.Duration
	SSL      time.Duration
	TTFB     time.Duration
	Total    time.Duration
	Transfer time.Duration
	//Transfer    time.Duration No Transfer time because we don't consume body
	Addr string
}

type conTrack struct {
	DNSStart             time.Time
	DNSDone              time.Time
	ConnectStart         map[string]time.Time
	ConnectDone          map[string]time.Time
	Addr                 string
	WroteRequest         time.Time
	GotFirstResponseByte time.Time
}

func (ct *conTrack) getConInfo() *ConInfo {
	ci := &ConInfo{
		Addr: ct.Addr,
	}
	if ct.GotFirstResponseByte.After(ct.WroteRequest) {
		ci.TTFB = ct.GotFirstResponseByte.Sub(ct.WroteRequest)
	}
	if ct.DNSDone.After(ct.DNSStart) {
		ci.DNS = ct.DNSDone.Sub(ct.DNSStart)
	}
	if ct.Addr == "" && len(ct.ConnectStart) > 0 { //If no addr(cause FAIL) but map has key(s) use any
		for ct.Addr = range ct.ConnectStart {
			//log.Println(ct.Addr)
		}
	}
	cs := ct.ConnectStart[ct.Addr]
	cd, ok := ct.ConnectDone[ct.Addr]
	if !ok {
		cd = time.Now()
	}
	if cd.After(cs) {
		ci.Connect = cd.Sub(cs)
	}
	if ct.WroteRequest.After(cd) {
		ci.SSL = ct.WroteRequest.Sub(cd)
	}
	ci.Total = ci.DNS + ci.Connect + ci.SSL + ci.TTFB
	return ci
}

func testoverhttp(u *url.URL, resolver string) (*ConInfo, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req = req.WithContext(ctx)

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
			//Patch in our hacked resolver which always returns same ip for everything
			Resolver: &net.Resolver{
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{}
					return d.DialContext(ctx, "udp", resolver+":53")
				},
			},
		}).DialContext,
		MaxIdleConns:          100,              //Irrelevant
		IdleConnTimeout:       90 * time.Second, //Irrelevant
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: time.Second * 20,
	}
	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, //Since we now use high-level client we must stop redirects.
	}

	//Initialize connection tracker
	ct := &conTrack{
		ConnectStart: make(map[string]time.Time),
		ConnectDone:  make(map[string]time.Time),
	}
	//Initialize httptrace
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			ct.Addr = connInfo.Conn.RemoteAddr().String()
			//log.Println(ct.Addr)
		},
		DNSStart: func(ds httptrace.DNSStartInfo) {
			ct.DNSStart = time.Now()
		},
		DNSDone: func(dd httptrace.DNSDoneInfo) {
			ct.DNSDone = time.Now()
		},
		ConnectStart: func(network, addr string) {
			ct.ConnectStart[addr] = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			ct.ConnectDone[addr] = time.Now()
		},
		GotFirstResponseByte: func() {
			ct.GotFirstResponseByte = time.Now()
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			ct.WroteRequest = time.Now()
		},
	}
	//Wrap trace into req
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	ti := ct.getConInfo()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got status: %v %v", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//Compute transfer time after consuming body
	ti.Transfer = time.Since(ct.GotFirstResponseByte)
	//Recompute Total for our purpose
	ti.Total = ti.DNS + ti.Connect + ti.SSL + ti.TTFB + ti.Transfer
	return ti, nil
}

func medianconinfo(results []*ConInfo) *ConInfo {
	validVals := make([]float64, 0)
	//Store the individual total times of successful result
	for _, ci := range results {
		if ci != nil {
			validVals = append(validVals, float64(ci.Total))
		}
	}
	if len(validVals) == 0 {
		//No valid tests
		return nil
	}
	median, _ := stats.Median(validVals)
	//Find the closest matching result to the median
	var bestmatch *ConInfo
	bestmatchdur := float64(time.Hour)
	for _, ci := range results {
		if ci != nil {
			if math.Abs(float64(ci.Total-time.Duration(median))) < bestmatchdur {
				bestmatchdur = math.Abs(float64(ci.Total - time.Duration(median)))
				bestmatch = ci
			}
		}
	}
	return bestmatch
}
