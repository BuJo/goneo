#
# build stage
#
FROM golang:1.17-alpine AS builder

RUN apk add git

COPY . /build
WORKDIR /build
ARG release
RUN CGO_ENABLED=0 go build -o goneo -ldflags "-s -w -extldflags '-static' -X main.buildversion=${release:-$(git describe --abbrev=0 --tags)-$(git rev-list -1 --abbrev-commit HEAD)}" -tags timetzdata ./cmd/goneo

#
# runtime image
#
FROM scratch AS runtime

WORKDIR /app

COPY --from=builder /build/goneo .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 7474

ENTRYPOINT ["/app/goneo"]
CMD ["-size", "universe"]
