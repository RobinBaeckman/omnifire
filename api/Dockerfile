FROM alpine:3.17

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY api api
COPY storage/migration storage/migration

ENTRYPOINT ["sh", "-c", "./api"]
