#Build server
FROM golang:1.23 as builder
RUN mkdir -p /go/src/folder
WORKDIR /go/src/folder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./cmd/app/

#Run server
FROM alpine:latest
WORKDIR /
COPY --from=builder /go/src/folder/app .
RUN chmod 0777 app
EXPOSE 9090
CMD ["/app"]