# 使用方式

## cmd

``` docker
wsl -l -v

wsl -d Ubuntu-24.04 --cd "D:\History\Git\MemoNest\docker\"
wsl -t Ubuntu-24.04

sudo docker-compose up -d
sudo docker-compose down
```

``` golang
wsl -d Ubuntu-24.04 --cd "D:\History\Git\MemoNest\work\"
go run ./cmd/.
```