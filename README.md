# dnsperfbench
DNS Performance Benchmarker


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

Assume the latest tag is `v0.1.1`

```
#One-off test
curl -Lo /tmp/dnsbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.1/dnsperfbench-linux && chmod +x /tmp/dnsbench && dnsbench #Linux
curl -Lo /tmp/dnsbench https://github.com/turbobytes/dnsperfbench/releases/download/v0.1.1/dnsperfbench-osx && chmod +x /tmp/dnsbench && dnsbench #OSX
```

To have it permanently available store the binary somewhere permanent

## Usage

Arguments

- `-resolver IP` Specify multiple times to test additional resolvers. Might be useful for comparing your ISP provided servers against public resolvers.
- `-r` Print the output in a machine readable format
- `-version` Print the version and exit
