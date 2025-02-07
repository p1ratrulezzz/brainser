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

buildgui-win:
	#rm -rf fyne-cross/dist
	go run github.com/fyne-io/fyne-cross@latest windows -arch=amd64,arm64 -app-id=com.jetbrainser.app -tags gui -output jetbrainser.exe ./src/app
	#go run github.com/fyne-io/fyne-cross@latest darwin --macosx-sdk-path=./SDKs/MacOSX12.3.sdk -arch=amd64 -app-id com.jetbrainser.app -tags gui -output jetbrainser ./src/app

buildgui-linux-x64:
	#go run github.com/fyne-io/fyne-cross@latest linux -arch=amd64 -tags gui -output jetbrainser ./src/app
	GOARCH=amd64 GOOS=linux go build -tags gui -o "${OUTDIR}/${BINARY_NAME}-gui-linux-x64" "${SRCPATH}"
	CGO_ENABLED=1 GOARCH=arm64 GOOS=linux go build -tags gui -o "${OUTDIR}/${BINARY_NAME}-gui-linux-arm" "${SRCPATH}"

buildgui-osx:
	#go run github.com/fyne-io/fyne-cross@latest darwin -arch=amd64,arm64 -app-id com.jetbrainser.app -tags gui -output jetbrainser ./src/app
	go run github.com/fyne-io/fyne-cross@latest darwin -arch=amd64,arm64 -app-id com.jetbrainser.app -tags gui --icon src/app/Icon.png -output jetbrainser ./src/app

clean:
	go clean
	# rm -f "${OUTDIR}/${BINARY_NAME}-linux_64"
