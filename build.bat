echo 'Building'

go version
set GO111MODULE=on
set GOPROXY=https://goproxy.io

set GOOS=windows
set GOARCH=amd64
go build -v -o "pit-win.exe" -ldflags "-s -w -H=windowsgui"

set GOOS=darwin
set GOARCH=amd64
go build -v -o "pit-darwin" -ldflags "-s -w"

set GOOS=linux
set GOARCH=amd64
go build -v -o "pit-linux" -ldflags "-s -w"
