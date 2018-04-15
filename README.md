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

### 
