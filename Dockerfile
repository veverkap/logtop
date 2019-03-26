FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -v -t -d ./...
RUN go build -o goreader ./reader/reader.go
CMD ["goreader"]
