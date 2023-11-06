#!/bin/bash
echo "compiling..."

GOOS=js GOARCH=wasm go build -pgo=auto -trimpath -gcflags=all="-B" -ldflags="-s -w" -o input.wasm
echo "optimizing..."
wasm-opt --enable-bulk-memory -O0 -o main.wasm input.wasm
rm input.wasm

echo "compressing..."
rm main.wasm.gz
gzip -9 main.wasm

cp main.wasm.gz ../goMMOServ/www/
echo "copied..."

scp -P 5313 main.wasm.gz dist@gosnake.go-game.net:~/goMMOServ/www/
rm main.wasm.gz
