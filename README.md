# dnsperfbench
DNS Performance Benchmarker

dnsperfbench is a command line tool to compare performance of popular public DNS resolvers from your computer, including Google Public DNS, OpenDNS and Cloudflare's 1.1.1.1.

You may optionally specify IP addresses of additional resolvers to include in the benchmark, for example the IP address of your ISP's resolver.

For each resolver, dnsperfbench first runs runs a few tests for a specific FQDN, to ensure the resolver has the response in cache.
Next, dnsperfbench tests cache hit performance (basically round trip latency; we call this ResolverHit) and cache miss performance against various major authoritative DNS providers.
Each test is repeated 15 times.

Once all the tests have finished, dnsperfbench shows the results per resolver (median and mean RTT and % tests that failed) and computes an overall performance score per resolver. The lower the score, the better was the resolver's performance.

Example of the overall Summary:

```
========== Summary ===========
Scores (lower is better)
+--------------------------------+-------------------+
|            RESOLVER            | PERFORMANCE SCORE |
+--------------------------------+-------------------+
| 1.1.1.1 (Cloudflare)           |              1983 |
| 208.67.222.222 (OpenDNS)       |              3175 |
| 199.85.126.20 (Norton)         |              3846 |
| 8.8.8.8 (Google)               |              6275 |
| 185.228.168.168 (Clean         |             12825 |
| Browsing)                      |                   |
| 9.9.9.9 (Quad9)                |             14677 |
| 8.26.56.26 (Comodo)            |            200494 |
+--------------------------------+-------------------+
You should probably use 1.1.1.1 as your default resolver
```

... and the Summary of running dnsperfbench at the same location but from another ISP:

```
========== Summary ===========
Scores (lower is better)
+--------------------------------+-------------------+
|            RESOLVER            | PERFORMANCE SCORE |
+--------------------------------+-------------------+
| 1.1.1.1 (Cloudflare)           |              1720 |
| 199.85.126.20 (Norton)         |              3031 |
| 8.8.8.8 (Google)               |              5580 |
| 185.228.168.168 (Clean         |              7205 |
| Browsing)                      |                   |
| 208.67.222.222 (OpenDNS)       |             12415 |
| 9.9.9.9 (Quad9)                |             62491 |
| 8.26.56.26 (Comodo)            |            131040 |
+--------------------------------+-------------------+
You should probably use 1.1.1.1 as your default resolver
```

While Cloudflare happened to be the best in both cases, performance of OpenDNS was quite different.

The performance score is calculated using the following formula:
`5 * (ResolverHit mean + ResolverHit median) + ( for each auth: (auth mean + auth median) )`
A failed test is treated as the test taking 10 seconds.

## httptest

New functionality `-httptest` benchmarks the performance of a given http(s) endpoint using the answer returned by various resolvers.

Example

```
$ dnsperfbench -resolver 115.178.58.10 -httptest https://turbobytes.akamaized.net/static/rum/100kb-image.jpg
2018/04/21 20:09:47 Resolving
2018/04/21 20:10:00 Issuing HTTP(s) tests
+-------------------------------------+---------------------+-------+---------+-------+-------+----------+--------+
|              RESOLVER               |       REMOTE        |  DNS  | CONNECT |  TLS  | TTFB  | TRANSFER | TOTAL  |
+-------------------------------------+---------------------+-------+---------+-------+-------+----------+--------+
| 115.178.58.10 (Unknown)             | 49.231.112.9:443    | 26ms  | 23ms    | 83ms  | 25ms  | 40ms     | 197ms  |
| 9.9.9.9 (Quad9)                     | 202.183.253.10:443  | 56ms  | 25ms    | 89ms  | 26ms  | 43ms     | 238ms  |
| 208.67.222.222 (OpenDNS)            | 49.231.112.33:443   | 71ms  | 23ms    | 83ms  | 25ms  | 40ms     | 242ms  |
| 8.8.8.8 (Google)                    | 49.231.112.9:443    | 110ms | 24ms    | 84ms  | 24ms  | 40ms     | 281ms  |
| [2001:4860:4860::8888] (Google)     | 49.231.112.33:443   | 111ms | 24ms    | 83ms  | 24ms  | 41ms     | 282ms  |
| [2620:0:ccc::2] (OpenDNS)           | 49.231.112.33:443   | 156ms | 24ms    | 82ms  | 24ms  | 40ms     | 326ms  |
| [2606:4700:4700::1111] (Cloudflare) | 23.49.60.192:443    | 53ms  | 53ms    | 114ms | 53ms  | 69ms     | 342ms  |
| 1.1.1.1 (Cloudflare)                | 23.49.60.192:443    | 54ms  | 56ms    | 120ms | 57ms  | 72ms     | 358ms  |
| 199.85.126.20 (Norton)              | 23.52.171.99:443    | 89ms  | 62ms    | 130ms | 60ms  | 85ms     | 426ms  |
| 180.76.76.76 (Baidu)                | 23.2.16.27:443      | 129ms | 86ms    | 180ms | 86ms  | 123ms    | 605ms  |
| [2620:fe::fe] (Quad9)               | 23.215.102.26:443   | 27ms  | 223ms   | 456ms | 224ms | 317ms    | 1.247s |
| [2a0d:2a00:1::] (Clean Browsing)    | 23.219.38.67:443    | 350ms | 228ms   | 465ms | 229ms | 278ms    | 1.55s  |
| 185.228.168.168 (Clean Browsing)    | 23.219.38.48:443    | 382ms | 227ms   | 463ms | 227ms | 292ms    | 1.591s |
| 119.29.29.29 (DNSPod)               | 205.197.140.136:443 | 474ms | 236ms   | 481ms | 244ms | 335ms    | 1.772s |
| 114.114.114.114 (114dns)            | 23.215.104.203:443  | 382ms | 283ms   | 573ms | 284ms | 401ms    | 1.924s |
| 8.26.56.26 (Comodo)                 | 104.86.110.154:443  | 219ms | 280ms   | 566ms | 279ms | 648ms    | 1.993s |
+-------------------------------------+---------------------+-------+---------+-------+-------+----------+--------+
```

Behind the scenes, the tool fetches `https://turbobytes.akamaized.net/static/rum/100kb-image.jpg` 10 times for each resolver, using the IP returned by the individual resolver, and reports the result with the median total time.

In the above example, we can see that the best performing IP were returned by my ISP resolver and ECS compatible public resolvers.

This test is bandwidth intensive. In the above example, it ran 10 tests each against 16 resolvers downloading 100KB each time. A total of ~16 MB + TLS handshake overhead.

IPv6 is supported by default.

## Installation

### From source

```
go get -u github.com/turbobytes/dnsperfbench
go install github.com/turbobytes/dnsperfbench
#Run using
$GOPATH/bin/dnsperfbench  #If $GOPATH/bin is not in PATH
or
dnsperfbench #If $GOPATH/bin is in PATH
```

Binary will be located at `$GOPATH/bin/dnsperfbench`

### Docker

```
docker pull turbobytes/dnsperfbench
#Run using
docker run --rm -it turbobytes/dnsperfbench
```

You can also set flags with the Docker image:

```
docker run --rm -it turbobytes/dnsperfbench -workers 2
```

### Download a release

View the page of the [latest release](https://github.com/turbobytes/dnsperfbench/releases/latest). Then copy download link from there.

Assume the latest tag is `v0.1.6`

```
#One-off test
curl -Lo /tmp/dnsperfbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.6/dnsperfbench-linux && chmod +x /tmp/dnsperfbench && /tmp/dnsperfbench #Linux
curl -Lo /tmp/dnsperfbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.6/dnsperfbench-osx && chmod +x /tmp/dnsperfbench && /tmp/dnsperfbench #OSX
```

To have it permanently available store the binary somewhere permanent.

## Usage

```
$ dnsperfbench --help
Usage of dnsperfbench:
-httptest string
    Specify a URL to test including protocol (http or https)
-queries int
    Limit the number of DNS queries in-flight at a time (default 5)
-r	Output raw mode
-resolver value
    Additional resolvers to test. default=199.85.126.20, 185.228.168.168, [2001:4860:4860::8888], [2620:fe::fe], 9.9.9.9, [2606:4700:4700::1111], [2a0d:2a00:1::], 208.67.222.222, 8.26.56.26, 8.8.8.8, 1.1.1.1, 114.114.114.114, 119.29.29.29, 180.76.76.76, [2620:0:ccc::2]
-version
    Print version and exit
-workers int
    Number of tests to run at once (default 15)
```
