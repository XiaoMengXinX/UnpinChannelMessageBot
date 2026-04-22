FROM golang:1.21 AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends git

WORKDIR /app

ARG REPO_URL=https://github.com/XiaoMengXinX/UnpinChannelMessageBot.git

RUN git clone --depth 1 $REPO_URL . 

RUN go mod tidy && \
    CGO_ENABLED=0 \
    go build \
        -a \
        -ldflags="-s -w" \
        -trimpath \
        -installsuffix cgo \
        -o UnpinBot .

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/UnpinBot .

USER nonroot:nonroot

CMD ["./UnpinBot"]
