FROM golang:1.20.4

WORKDIR /usr/src/app

# install all dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy source code
COPY . .
RUN go build -v -o /usr/local/bin/app/coride-backend ./cmd/server.go

CMD ["/usr/local/bin/app/coride-backend"]
