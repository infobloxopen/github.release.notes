# build the release-notes binary
FROM golang:1.16.3 AS builder
LABEL stage=server-intermediate
WORKDIR /go/src/github.com/infobloxopen/github.release.notes

COPY . .
RUN go build -mod=vendor -o bin/release-notes ./release-notes

# copy the release-notes binary from builder stage; run the release-notes binary
FROM alpine:latest AS runner
WORKDIR /bin

# Go programs require libc
RUN mkdir -p /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder /go/src/github.com/infobloxopen/github.release.notes/bin/release-notes .
COPY /templates/* ./


ENTRYPOINT ["release-notes"]
