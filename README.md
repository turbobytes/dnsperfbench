# dnsperfbench
DNS Performance Benchmarker

dnsperfbench compares the performance of popular public DNS resolvers from your computer. It for each resolver it tests cache hit (basically round trip latency), and cache miss against various authoritative managed DNS providers. Each test is repeated 15 times. Once the tests are finished, it presents a summary.

For example with one of my ISPs

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

and from another ISP at same location

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

While Cloudflare happened to be the fastest in both cases, the performance of OpenDNS was different between ISPs.

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

Visit the age for the [latest release](https://github.com/turbobytes/dnsperfbench/releases/latest). Then copy download link from there.

Assume the latest tag is `v0.1.2`

```
#One-off test
curl -Lo /tmp/dnsperfbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.2/dnsperfbench-linux && chmod +x /tmp/dnsperfbench && /tmp/dnsperfbench #Linux
curl -Lo /tmp/dnsperfbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.2/dnsperfbench-osx && chmod +x /tmp/dnsperfbench && /tmp/dnsperfbench #OSX
```

To have it permanently available store the binary somewhere permanent

## Usage

Arguments

- `-resolver IP` Specify multiple times to test additional resolvers. Might be useful for comparing your ISP provided servers against public resolvers. IPv6 goes in [square brackets]
- `-r` Print the output in a machine readable format
- `-version` Print the version and exit
