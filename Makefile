VERSION=v0.1.2
GOVERSION=$(shell go version)
LDFLAGS='-X main.versionString=${VERSION} -X "main.goVersionString=${GOVERSION}"'

docker:
	mkdir -p bin/
	CGO_ENABLED=0 go build -ldflags ${LDFLAGS} -o bin/dnsperfbench main.go
	docker build -t turbobytes/dnsperfbench .
	docker tag turbobytes/dnsperfbench turbobytes/dnsperfbench:$(VERSION)
	docker push turbobytes/dnsperfbench:latest
	docker push turbobytes/dnsperfbench:$(VERSION)

release:
	#Make release assets
	#echo ${LDFLAGS}
	CGO_ENABLED=0 go build -ldflags ${LDFLAGS} -o bin/dnsperfbench-linux main.go
	CGO_ENABLED=0 GOOS=darwin go build -ldflags ${LDFLAGS} -o bin/dnsperfbench-osx main.go
	CGO_ENABLED=0 GOOS=windows go build -ldflags ${LDFLAGS} -o bin/dnsperfbench.exe main.go
