FROM golang:1.10

WORKDIR /go/src/app
COPY . .

ENV GOPATH /go/src/app
ENV GOBIN  /go/src/app/bin

RUN go get -d -v ./... && go install -v cmd/alien-invasion/main.go

CMD ["bin/main"]
