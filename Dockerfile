FROM golang:1.20-alpine as builder
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD pkg /app/pkg
ADD scenarios /app/scenarios
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o fork-bomb scenarios/fork-bomb/cmd/main.go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o data-overload scenarios/data-overload/cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/fork-bomb ./fork-bomb
COPY --from=builder /app/data-overload ./data-overload
CMD ["./fork-bomb"]
