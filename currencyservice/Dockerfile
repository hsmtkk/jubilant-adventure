FROM golang:1.19 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY env/ /app/env/
COPY currencyservice/ /app/currencyservice/
WORKDIR /app/currencyservice
RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /app/currencyservice/currencyservice /usr/local/bin/currencyservice
ENTRYPOINT ["/usr/local/bin/currencyservice"]
