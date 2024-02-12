FROM golang:1.22-alpine as builder
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN \
  apk --no-cache --update add git make && \
  make build

FROM alpine as runner
ENV PORT 8877
WORKDIR /app
EXPOSE $PORT
COPY --from=builder /app/bin/jrebel /app/jrebel
CMD /app/jrebel -h 0.0.0.0 -p ${PORT}
