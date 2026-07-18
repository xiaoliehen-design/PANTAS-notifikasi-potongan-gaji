FROM golang:1.26.5-alpine3.24 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/pantas ./cmd/server

FROM alpine:3.24
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S pantas \
    && adduser -S -G pantas -H -s /sbin/nologin pantas
ENV TZ=Asia/Jakarta PORT=10000
WORKDIR /app
COPY --from=build /out/pantas /app/pantas
USER pantas
EXPOSE 10000
ENTRYPOINT ["/app/pantas"]
