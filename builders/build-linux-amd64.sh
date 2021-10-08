#!/bin/bash
cd ../ || exit
go env -w CGO_ENABLED=0
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -ldflags '-w -s' -gcflags '-l' -a -o esh-linux-amd64
upx -9 esh-linux-amd64
mv ./esh-linux-amd64 ./pkg/esh-linux-amd64
chmod 777 pkg/esh-linux-amd64
go env -w CGO_ENABLED=1
go env -w GOOS=linux
go env -w GOARCH=amd64
cd ./builders/ || exit
echo "esh-linux-amd64 success"
