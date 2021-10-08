#!/bin/bash
cd ../ || exit
go env -w CGO_ENABLED=0
go env -w GOOS=linux
go env -w GOARCH=arm64
go build -ldflags '-w -s' -gcflags '-l' -a -o esh-linux-arm64
upx -9 esh-linux-arm64
mv ./esh-linux-arm64 ./pkg/esh-linux-arm64
chmod 777 .pkg/esh-linux-arm64
go env -w CGO_ENABLED=1
go env -w GOOS=linux
go env -w GOARCH=amd64
cd ./builders/ || exit
echo "esh-linux-arm64 success"
