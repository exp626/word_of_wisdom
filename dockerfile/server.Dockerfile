FROM golang:1.21

WORKDIR /app

COPY go.mod /
RUN go mod download

COPY cmd/server/ /cmd/server/
COPY pkg/ /pkg/

RUN go build -o /server /cmd/server/*.go

EXPOSE 8080

CMD ["/server"]
