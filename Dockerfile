FROM golang:1.6
COPY . /go/src/github.com/QuentinPerez/c14-cli
WORKDIR /go/src/github.com/QuentinPerez/c14-cli
RUN go install -v ./cmd/c14
ENTRYPOINT ["c14"]
CMD ["help"]
