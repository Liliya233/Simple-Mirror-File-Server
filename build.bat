SET CGO_ENABLED=0

SET GOOS=windows
go build -ldflags="-w -s" -o ./build/fileserver.exe
SET GOOS=linux
go build -ldflags="-w -s" -o  ./build/fileserver
