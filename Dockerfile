FROM golang:1.19-alpine as builder
WORKDIR /app/
ENV GOPROXY=https://goproxy.cn
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o jrebel

FROM alpine as runner
ENV PORT=8877
WORKDIR /app
EXPOSE $PORT
COPY --from=builder /app/jrebel /app
CMD /app/jrebel -h 0.0.0.0 -p ${PORT}
