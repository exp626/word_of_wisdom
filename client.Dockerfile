FROM golang:1.21.3

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY cmd/client/ /cmd/client/
COPY pkg/ /pkg/

RUN go build -o /client /cmd/client/*.go

CMD ["/client"]
