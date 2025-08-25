FROM golang:alpine AS builder
WORKDIR /work
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -tags netgo \
  -ldflags="-extldflags=-static" \
  -o github-forgejo-backup \
  && chmod +x github-forgejo-backup

FROM scratch
WORKDIR /work
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /work/github-forgejo-backup /usr/local/bin/github-forgejo-backup
CMD [ "github-forgejo-backup" ]
