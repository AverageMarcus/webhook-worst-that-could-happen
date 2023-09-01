FROM golang:1.20 as builder
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD pkg /app/pkg
ADD scenarios /app/scenarios
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fork-bomb scenarios/fork-bomb/cmd/main.go

FROM alpine:3.18
RUN apk update && apk --no-cache  add ca-certificates
WORKDIR /app
COPY --from=builder /app/fork-bomb ./fork-bomb
CMD ["./fork-bomb"]
