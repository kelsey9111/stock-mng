FROM golang:alpine AS builder  

WORKDIR /build

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crm.com ./cmd/server

FROM alpine  

WORKDIR /

COPY ./config /config

COPY --from=builder /build/crm.com /
RUN chmod +x /crm.com

CMD ["/crm.com", "config/local.yaml"]
