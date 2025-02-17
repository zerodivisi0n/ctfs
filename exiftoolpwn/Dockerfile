FROM golang:1.24-alpine3.21 AS builder

ARG TARGETOS TARGETARCH

WORKDIR /app

COPY . /app

RUN go build .

FROM alpine:3.21

ENV PORT=3000

LABEL org.opencontainers.image.source=https://github.com/zerodivisi0n/exiftoolpwn

RUN apk --no-cache add \
            curl \
            util-linux \
            py3-pip \
            exiftool \
            libreoffice-common \
            libreoffice-writer \
            font-droid-nonlatin \
            font-droid \
            font-dejavu \
            font-freefont \
            font-liberation \
        && curl -Ls https://raw.githubusercontent.com/dagwieers/unoconv/master/unoconv -o /usr/local/bin/unoconv \
        && chmod +x /usr/local/bin/unoconv

COPY --from=builder /app/exiftoolpwn /usr/local/bin  

ENTRYPOINT ["/usr/local/bin/exiftoolpwn"]