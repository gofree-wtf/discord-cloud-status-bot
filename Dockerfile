FROM golang:1.16rc1
LABEL maintainer="gofree.wtf@gmail.com"

ENV LOG_LEVEL=info
ENV LOG_FORMAT=console
ENV LOG_SESSION_LEVEL=warn

ENV BOT_TOKEN=INSERT_HERE
ENV BOT_COMMAND_PREFIX=!cs
ENV BOT_TIMEZONE=Asia/Seoul

ENV API_PORT=8080
ENV API_SELF_HEALTHCHECK_ENABLED=true
ENV API_SELF_HEALTHCHECK_URL=INSERT_HERE
ENV API_SELF_HEALTHCHECK_PERIOD_MINUTES=5

WORKDIR /discord-cloud-status-bot

ADD . /discord-cloud-status-bot/

RUN make init && \
    make build

ENTRYPOINT ["bin/discord-cloud-status-bot"]
