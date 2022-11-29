FROM alpine:3.17

RUN apk add --no-cache ffmpeg

COPY ./go-cqhttp_linux_amd64/go-cqhttp /bin/go-cqhttp

EXPOSE 80

WORKDIR /gocq
CMD ["go-cqhttp", "-faststart"]