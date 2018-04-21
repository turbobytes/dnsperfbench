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
