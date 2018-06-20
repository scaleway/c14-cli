FROM golang:1.6
COPY . /go/src/github.com/scaleway/c14-cli
WORKDIR /go/src/github.com/scaleway/c14-cli
RUN go install -v ./cmd/c14
ENTRYPOINT ["c14"]
CMD ["help"]
