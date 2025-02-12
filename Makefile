BINARY_NAME=jetbrainser
OUTDIR=bin
SRCPATH=jetbrainser/src/app
VERSION=0.0.9
BUILD_ID=$(shell date +"%Y%m%d%H%M%S")

build:
	rm -f bin/*
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-${BUILD_ID}-linux-x64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-${BUILD_ID}-win-x64.exe" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-${BUILD_ID}-linux-arm64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-${BUILD_ID}-win-arm64.exe" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-${BUILD_ID}-osx-arm64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-${BUILD_ID}-osx-amd64" "${SRCPATH}"
	chmod +x bin/*

buildgui-win:
	go run github.com/fyne-io/fyne-cross@latest windows -arch=amd64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id=com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-${BUILD_ID}-win-amd64.exe" ./src/app
	go run github.com/fyne-io/fyne-cross@latest windows -arch=arm64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id=com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-${BUILD_ID}-win-arm64.exe" ./src/app

buildgui-linux-amd64:
	go run github.com/fyne-io/fyne-cross@latest linux -image=fyne-cross-custom:linux -arch=amd64 -tags gui -app-version="${VERSION}" -app-build="${BUILD_ID}" --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-${BUILD_ID}-linux-amd64" ./src/app

buildgui-linux-arm64:
	go run github.com/fyne-io/fyne-cross@latest linux -image=fyne-cross-custom:linux -arch=arm64 -tags gui -app-version="${VERSION}" -app-build="${BUILD_ID}" --icon src/app/Icon.png -output "jetbrainser-${VERSION}-${BUILD_ID}-linux-arm64" ./src/app

buildgui-osx:
	go run github.com/fyne-io/fyne-cross@latest darwin -arch=amd64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-${BUILD_ID}-amd64" ./src/app
	go run github.com/fyne-io/fyne-cross@latest darwin -arch=arm64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-${BUILD_ID}-arm64" ./src/app

build-non-macos: clean build buildgui-win buildgui-linux-amd64


clean:
	go clean