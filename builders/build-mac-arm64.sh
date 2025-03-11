#!/bin/bash

CGO_ENABLED_ORI=`go env CGO_ENABLED`
GOOS_ORI=`go env GOOS`
GOARCH_ORI=`go env GOARCH`

cd ../ || exit
go env -w CGO_ENABLED=0
go env -w GOOS=darwin
go env -w GOARCH=arm64
go build -ldflags '-w -s' -gcflags '-l' -a -o pkg/esh-mac-arm64
upx -9 pkg/esh-mac-arm64
chmod 777 pkg/esh-mac-arm64
go env -w CGO_ENABLED=$CGO_ENABLED_ORI
go env -w GOOS=$GOOS_ORI
go env -w GOARCH=$GOARCH_ORI
cd ./builders/ || exit
echo "esh-mac-arm64 success"