#!/bin/bash

CGO_ENABLED_ORI=`go env CGO_ENABLED`
GOOS_ORI=`go env GOOS`
GOARCH_ORI=`go env GOARCH`

cd ../ || exit
## THe libvirt-go package is a CGo binding to the native libvirt platform library.As such it is not possible to disable CGO when building it
go env -w CGO_ENABLED=0
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -ldflags '-w -s' -gcflags '-l' -a -o pkg/esh.exe
#upx -9 ../tem/ztest
chmod 777 pkg/esh.exe
go env -w CGO_ENABLED=$CGO_ENABLED_ORI
go env -w GOOS=$GOOS_ORI
go env -w GOARCH=$GOARCH_ORI
cd ./builders/ || exit
echo "ztest.exe build success"