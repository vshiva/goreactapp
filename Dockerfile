FROM reactgo/builder-ci:latest as builder
ARG LD_FLAGS

COPY ./ /go/src/github.com/vshiva/goreactapp/
WORKDIR /go/src/github.com/vshiva/goreactapp/

RUN cd web && yarn install && yarn build

RUN go get -u github.com/golang/dep/cmd/dep && dep ensure && \
    go get -u github.com/go-bindata/go-bindata/... && go generate ./web && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "${LD_FLAGS}" -o bin/goreactapp ./cmd

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/vshiva/goreactapp/bin/goreactapp /usr/bin
ENTRYPOINT ["/usr/bin/goreactapp"]
CMD ["-h"]
