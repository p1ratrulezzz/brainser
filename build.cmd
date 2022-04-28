SET BINARY_NAME=jetbrainser
SET OUTDIR=bin
SET SRCPATH=jetbrainser/src/app

SET GOARCH=amd64
SET GOOS=linux
go build -o "%OUTDIR%/%BINARY_NAME%-linux-x64" "%SRCPATH%"

SET GOARCH=amd64
SET GOOS=windows
go build -o "%OUTDIR%/%BINARY_NAME%-win-x64.exe" "%SRCPATH%"

SET GOARCH=arm64
SET GOOS=linux
go build -o "%OUTDIR%/%BINARY_NAME%-linux-arm64" "%SRCPATH%"

SET GOARCH=arm64
SET GOOS=windows
go build -o "%OUTDIR%/%BINARY_NAME%-win-arm64.exe" "%SRCPATH%"