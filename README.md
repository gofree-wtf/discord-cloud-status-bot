discord-cloud-status-bot
========================

여러 클라우드 서비스의 상태 대시보드를 주기적으로 체크하여, 변경점이 있을 때 알림을 주는 Discord 봇 입니다.

## Developments

- Go 1.15
- Go modules

## Run

```bash
$ make init
$ make run
```

Docker를 사용한다면:

```bash
$ make docker-build
$ make docker-run
```
