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

``` ip
ip -4 addr

http://172.31.235.34:8080/
```

``` ConEmu
%windir%\system32\wsl.exe -cur_console:t:Linux -d Ubuntu-24.04 --cd "D:\History\Git\MemoNest\docker\"
%windir%\system32\wsl.exe -cur_console:t:Golang -d Ubuntu-24.04 --cd "D:\History\Git\MemoNest\work\"
```


層級 | 責任
-- | --
domain/repo | 資料存取抽象（DB、Cache、API）
domain/service | 業務邏輯抽象（加密、驗證、Session、通知）