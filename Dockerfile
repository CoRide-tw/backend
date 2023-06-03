FROM golang:1.20.4
RUN mkdir /src
WORKDIR /src
COPY . /src/
RUN go build -o ./bin/coride
EXPOSE 8080
CMD ["./bin/coride"]