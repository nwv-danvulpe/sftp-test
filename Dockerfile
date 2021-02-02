FROM golang:1.15-alpine as builder

WORKDIR /workspace

COPY . .

RUN go build -o sftp-test && \
    chmod +x sftp-test

FROM alpine:latest

COPY --from=builder /workspace/sftp-test /usr/bin/sftp-test

ENTRYPOINT [ "/usr/bin/sftp-test" ]