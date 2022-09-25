#!/bin/bash

CGO_ENABLED_ORI=`go env CGO_ENABLED`
GOOS_ORI=`go env GOOS`
GOARCH_ORI=`go env GOARCH`

cd ../ || exit
go env -w CGO_ENABLED=0
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -ldflags '-w -s' -gcflags '-l' -a -o pkg/esh.exe
chmod 777 pkg/esh.exe
go env -w CGO_ENABLED=$CGO_ENABLED_ORI
go env -w GOOS=$GOOS_ORI
go env -w GOARCH=$GOARCH_ORI
cd ./builders/ || exit
echo "ztest.exe build success"