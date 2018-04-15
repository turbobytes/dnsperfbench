docker:
	mkdir -p bin/
	CGO_ENABLED=0 go build -o bin/dnsperfbench main.go
	docker build -t turbobytes/dnsperfbench .
	docker push turbobytes/dnsperfbench

release:
	#Make release assets
	CGO_ENABLED=0 go build -o bin/linux/dnsperfbench main.go
	CGO_ENABLED=0 GOOS=darwin go build -o bin/osx/dnsperfbench main.go
	CGO_ENABLED=0 GOOS=windows go build -o bin/windows/dnsperfbench.exe main.go
