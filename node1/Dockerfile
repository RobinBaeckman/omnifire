FROM alpine:3.17

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY node1 node1
COPY storage/migration storage/migration

ENTRYPOINT ["sh", "-c", "./node1"]
