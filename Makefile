BINARY_NAME=jetbrainser
OUTDIR=bin
SRCPATH=jetbrainser/src/app

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags console -o "${OUTDIR}/${BINARY_NAME}-linux-x64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -tags console -o "${OUTDIR}/${BINARY_NAME}-win-x64.exe" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -tags console -o "${OUTDIR}/${BINARY_NAME}-linux-arm64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -tags console -o "${OUTDIR}/${BINARY_NAME}-win-arm64.exe" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -tags console -o "${OUTDIR}/${BINARY_NAME}-osx-arm64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -tags console -o "${OUTDIR}/${BINARY_NAME}-osx-amd64" "${SRCPATH}"
	chmod +x bin/*

buildgui:
	rm -rf fyne-cross/dist
	go run github.com/fyne-io/fyne-cross linux -arch=amd64,arm64 -tags gui -output jetbrainser ./src/app
	go run github.com/fyne-io/fyne-cross windows -arch=amd64 -tags gui -output jetbrainser.exe ./src/app
	go run github.com/fyne-io/fyne-cross darwin -app-id com.jetbrainser.app -arch=amd64,arm64 -tags gui -output jetbrainser ./src/app

clean:
	go clean
	# rm -f "${OUTDIR}/${BINARY_NAME}-linux_64"
