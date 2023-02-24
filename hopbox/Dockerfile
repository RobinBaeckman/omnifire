FROM alpine:3.17

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY hopbox hopbox
COPY storage/migration storage/migration
COPY env/sample.config env/config

ENTRYPOINT ["sh", "-c", "./hopbox"]
