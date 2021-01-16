discord-cloud-status-bot
========================

여러 클라우드 서비스의 상태 대시보드를 주기적으로 체크하여, 변경점이 있을 때 알림을 주는 Discord 봇 입니다.

## Developments

- Go 1.15
- Go modules

## How to run

### Step 1. 디스코드 앱 및 봇 생성

https://discord.com/developers/applications

### Step 2. 봇 설정 파일을 마련합니다.

```bash
cp configs/app.yaml.example configs/app.yaml
```

이 후 app.yaml 내의 `bot.token`을 입력합니다.

### Step 3. 봇 실행

다음과 같이 로컬에서 실행합니다:

```bash
$ make init
$ make run
```

Docker를 사용한다면:

```bash
$ make docker-build
$ make docker-run
```

### Step 4. 서버에 봇 초대

(INSERT_HERE에 자신의 앱 클라이언트 ID로 대체)

https://discord.com/oauth2/authorize?client_id=INSERT_HERE&scope=bot&permissions=19456

### Step 5. 채널에 메세지 보내보기

Example: `!cs test`

이 후, 봇 로그가 출력된다면 성공
