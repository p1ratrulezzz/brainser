BINARY_NAME=jetbrainser
OUTDIR=bin
SRCPATH=jetbrainser/src/app

build:
	#GOARCH=amd64 GOOS=darwin go build -o "${OUTDIR}/${BINARY_NAME}-osx64" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o "${OUTDIR}/${BINARY_NAME}-linux-x64" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=amd64 GOOS=windows go build -o "${OUTDIR}/${BINARY_NAME}-win-x64.exe" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=arm64 GOOS=linux go build -o "${OUTDIR}/${BINARY_NAME}-linux-arm64" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=arm64 GOOS=windows go build -o "${OUTDIR}/${BINARY_NAME}-win-arm64.exe" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin go build -o "${OUTDIR}/${BINARY_NAME}-osx-arm64" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -o "${OUTDIR}/${BINARY_NAME}-osx-amd64" "${SRCPATH}"
	chmod +x bin/*

clean:
	go clean
	# rm -f "${OUTDIR}/${BINARY_NAME}-linux_64"
