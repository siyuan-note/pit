echo 'Building Kernel'

go version
set GO111MODULE=on
set GOPROXY=https://goproxy.io

set GOOS=windows
set GOARCH=amd64
go build -v -o "updater-win.exe" -ldflags "-s -w -H=windowsgui"

set GOOS=darwin
set GOARCH=amd64
go build -v -o "updater-darwin" -ldflags "-s -w"

set GOOS=linux
set GOARCH=amd64
go build -v -o "updater-linux" -ldflags "-s -w"
