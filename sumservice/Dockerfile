FROM golang:1.19 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY sumservice/ /app/sumservice/
WORKDIR /app/sumservice
RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /app/sumservice/sumservice /usr/local/bin/sumservice
ENTRYPOINT ["/usr/local/bin/sumservice"]
