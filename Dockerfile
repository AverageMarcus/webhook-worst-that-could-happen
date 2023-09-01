FROM golang:1.20-alpine as builder
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD pkg /app/pkg
ADD scenarios /app/scenarios
RUN mkdir -p /out
RUN for cmd in ./scenarios/*; do cmd=${cmd#./scenarios/}; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /out/$cmd scenarios/$cmd/cmd/main.go; done

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /out/ ./
CMD ["./fork-bomb"]
