FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 \
    go build \
        -a \
        -ldflags="-s -w" \
        -trimpath \
        -installsuffix cgo \
        -o UnpinBot .

FROM gcr.io/distroless/static:nonroot

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/UnpinBot /app/
WORKDIR /app

CMD ["./UnpinBot"]
