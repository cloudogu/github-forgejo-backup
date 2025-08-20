FROM golang:alpine AS builder
WORKDIR /work
COPY . .
RUN apk update && apk add --no-cache ca-certificates git
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X main.GitTag=$(git tag -l --sort=-v:refname)" -tags netgo -o github-forgejo-backup && chmod +x github-forgejo-backup

FROM alpine:3
WORKDIR /work
COPY --from=builder /work/github-forgejo-backup /usr/local/bin/github-forgejo-backup
CMD [ "github-forgejo-backup" ]
