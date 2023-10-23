#!/bin/bash
echo "compiling..."

GOOS=js GOMAXPROCS=1 GOARCH=wasm go build -tags=ebitensinglethread -pgo=auto -trimpath -gcflags=all="-B" -ldflags="-s -w" -o main.wasm

echo "compressing..."
rm main.wasm.gz
gzip -9 main.wasm

cp main.wasm.gz ../goMMOServ/www/
echo "copied..."

scp -P 5313 main.wasm.gz dist@gosnake.go-game.net:~/goMMOServ/www/