FROM golang:1.20 as builder
WORKDIR /app
ADD src /app/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM alpine:3.18
RUN apk update && apk --no-cache  add ca-certificates
WORKDIR /app
COPY --from=builder /app/main ./main
CMD ["./main"]
