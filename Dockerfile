FROM golang:1.15
LABEL maintainer="gofree.wtf@gmail.com"

WORKDIR /discord-cloud-status-bot

ADD . /discord-cloud-status-bot/

RUN make init && \
    make build

ENTRYPOINT ["bin/discord-cloud-status-bot"]
