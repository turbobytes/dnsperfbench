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
docker run --rm -it turbobytes/dnsperfbench dnsperfbench
```

### Download a release

View the page of the [latest release](https://github.com/turbobytes/dnsperfbench/releases/latest). Then copy download link from there.

Assume the latest tag is `v0.1.3`

```
#One-off test
curl -Lo /tmp/dnsperfbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.2/dnsperfbench-linux && chmod +x /tmp/dnsperfbench && /tmp/dnsperfbench #Linux
curl -Lo /tmp/dnsperfbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.2/dnsperfbench-osx && chmod +x /tmp/dnsperfbench && /tmp/dnsperfbench #OSX
```

To have it permanently available store the binary somewhere permanent.

## Usage

Arguments

- `-resolver IP` Specify multiple times to test additional resolvers. Might be useful for comparing your ISP provided resolver against public resolvers. IPv6 goes in [square brackets]
- `-r` Print the output in a machine readable format
- `-version` Print the version and exit
