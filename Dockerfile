FROM golang:1.25-alpine AS builder
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN \
  apk --no-cache --update add git make && \
  make build

FROM alpine AS runner
WORKDIR /app
EXPOSE 8877
COPY --from=builder /app/bin/jrebel /app/jrebel
ENTRYPOINT ["/app/jrebel"]
CMD ["-h", "0.0.0.0", "-p", "8877"]
