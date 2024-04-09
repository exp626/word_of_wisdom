FROM golang:1.21.3

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY cmd/server/ /cmd/server/
COPY pkg/ /pkg/

RUN go build -o /server /cmd/server/*.go

EXPOSE 8080

CMD ["/server"]
