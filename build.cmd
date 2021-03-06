SET CGO_ENABLED=0
SET BINARY_NAME=jetbrainser
SET OUTDIR=bin
SET SRCPATH=jetbrainser/src/app

SET GOARCH=amd64
SET GOOS=linux
go build -tags console -o "%OUTDIR%/%BINARY_NAME%-linux-x64" "%SRCPATH%"

SET GOARCH=amd64
SET GOOS=windows
go build -tags console -o "%OUTDIR%/%BINARY_NAME%-win-x64.exe" "%SRCPATH%"

SET GOARCH=arm64
SET GOOS=linux
go build -tags console -o "%OUTDIR%/%BINARY_NAME%-linux-arm64" "%SRCPATH%"

SET GOARCH=arm64
SET GOOS=windows
go build -tags console -o "%OUTDIR%/%BINARY_NAME%-win-arm64.exe" "%SRCPATH%"

SET GOARCH=arm64
SET GOOS=darwin
go build -tags console -o "%OUTDIR%/%BINARY_NAME%-osx-arm64" "%SRCPATH%"

SET GOARCH=amd64
SET GOOS=darwin
go build -tags console -o "%OUTDIR%/%BINARY_NAME%-osx-amd64" "%SRCPATH%"