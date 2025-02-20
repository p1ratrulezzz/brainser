BINARY_NAME=jetbrainser
OUTDIR=bin
SRCPATH=jetbrainser/src/app
VERSION=0.0.11
BUILD_ID=$(shell date +"%Y%m%d%H%M%S")

docker-build:
	docker buildx bake -f docker-compose.yml --load

build:
	rm -f bin/*
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-linux-x64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-win-x64.exe" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-linux-arm64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-win-arm64.exe" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-osx-arm64" "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -tags console -o "${OUTDIR}/${BINARY_NAME}-${VERSION}-osx-amd64" "${SRCPATH}"
	chmod +x bin/*

buildgui-win:
	go run github.com/fyne-io/fyne-cross@latest windows -arch=amd64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id=com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-win-amd64.exe" ./src/app
	go run github.com/fyne-io/fyne-cross@latest windows -arch=arm64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id=com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-win-arm64.exe" ./src/app

buildgui-linux-amd64:
	go run github.com/fyne-io/fyne-cross@latest linux -image=fyne-cross-custom:linux-amd64 -arch=amd64 -tags gui -app-version="${VERSION}" -app-build="${BUILD_ID}" --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-linux-amd64" ./src/app

buildgui-linux-arm64:
	go run github.com/fyne-io/fyne-cross@latest linux -image=fyne-cross-custom:linux-arm64 -arch=arm64 -tags gui -app-version="${VERSION}" -app-build="${BUILD_ID}" --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-linux-arm64" ./src/app

buildgui-osx:
	go run github.com/fyne-io/fyne-cross@latest darwin -arch=amd64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-amd64" ./src/app
	go run github.com/fyne-io/fyne-cross@latest darwin -arch=arm64 -app-version="${VERSION}" -app-build="${BUILD_ID}" -app-id com.jetbrainser.app -tags gui --icon src/app/Icon.png -output "jetbrainser-gui-${VERSION}-arm64" ./src/app

build-non-macos: clean build buildgui-win buildgui-linux-amd64

clean:
	go clean