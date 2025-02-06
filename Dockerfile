FROM golang:1.23 as build

ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN cd app && go build -o /build/yttg -ldflags "-s -w"

FROM debian:bookworm-slim

COPY --from=build /build/yttg /srv/yttg
COPY app/web/templates /srv/app/web/templates
COPY app/web/static /srv/app/web/static

RUN apt update && \
    apt install -y --no-install-recommends \
        ca-certificates \
        ffmpeg \
        python3 \
        python3-pip && \
    rm -rf /var/lib/apt/lists/*
RUN pip3 install --break-system-packages --no-cache-dir --no-deps -U yt-dlp

WORKDIR /srv
CMD ["/srv/yttg"]
