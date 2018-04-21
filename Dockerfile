FROM golang:1.10.1-alpine

RUN apk update && apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/turbobytes/dnsperfbench

COPY . .

RUN dep ensure
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' github.com/turbobytes/dnsperfbench

FROM scratch

COPY --from=0 /go/bin/dnsperfbench /bin/dnsperfbench
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/bin/dnsperfbench" ]
