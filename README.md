# WHU saves all webhooks and allow pull it anytime

## Usage

Download: `go get -u github.com/zhuharev/whu`

Change dirrectory `cd /path/to/github.com/zhuharev/whu`

Run: `PORT=2019 DB_PATH=db.storm go run who.go`

Create an webhook data:

```curl
curl -X POST \
  http://localhost:2019/1 \
  -d '[{"allo":"da"},{}]'
  ```

Receive updates:

```curl
curl -X GET \
  'http://localhost:2019/1/updates?offset=3'
```