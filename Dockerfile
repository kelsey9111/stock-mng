FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

COPY go.mod go.sum ./
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crm_server ./cmd/server


FROM alpine:latest  

WORKDIR /


COPY --from=builder /app/config /config


COPY --from=builder /app/crm_server /crm_server


RUN chmod +x /crm_server


CMD ["/crm_server", "config/local.yaml"]