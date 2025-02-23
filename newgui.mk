APPID=com.jetbrainser.gui
APPNAME=Jetbrainser GUI
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TAGS=guinew
LD_FLAGS=-X main.Version=\"${VERSION}\" -X main.BuildNumber=\"${BUILD_ID}\"

clean-newgui:
	rm -rf bin/osx
	rm -f bin/*
macos-newgui:
	rm -rf bin/osx
	rm -rf dist/darwin
	mkdir -p ./dist/darwin
	mkdir -p ./bin/osx
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w ${LD_FLAGS}" -o bin/${BINARY_NAME}-darwin-arm64 "${SRCPATH}"
	cp bin/"${BINARY_NAME}-darwin-arm64" bin/osx/app
	go run ./tools/macapp.go \
    		-assets ./bin/osx \
    		-bin app \
    		-icon ./Icon.png \
    		-identifier "${APPID}" \
    		-name "${APPNAME}-arm64" \
    		-o ./dist/darwin

	rm -r bin/osx/*
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w ${LD_FLAGS}" -o bin/${BINARY_NAME}-darwin-amd64 "${SRCPATH}"
	cp "bin/${BINARY_NAME}-darwin-amd64" bin/osx/app
	go run ./tools/macapp.go \
			-assets ./bin/osx \
			-bin app \
			-icon ./Icon.png \
			-identifier "${APPID}" \
			-name "${APPNAME}-amd64" \
			-o ./dist/darwin
	rm -rf bin/osx

windows-newgui:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w -H=windowsgui ${LD_FLAGS}" -o bin/${BINARY_NAME}-windows-arm64.exe "${SRCPATH}"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w -H=windowsgui ${LD_FLAGS}" -o bin/${BINARY_NAME}-windows-amd64.exe "${SRCPATH}"

linux-arm-newgui:
	CGO_ENABLED=1 \
	GOARCH=arm64 \
	GOOS=linux \
	go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w ${LD_FLAGS}" -o bin/${BINARY_NAME}-linux-arm64 "${SRCPATH}"

linux-amd64-newgui:
	CGO_ENABLED=1 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w ${LD_FLAGS}" -o bin/${BINARY_NAME}-linux-amd64 "${SRCPATH}"

linux-docker-newgui:
	docker compose build linuxgo-arm linuxgo-amd64
	docker run --rm --user "$(id -u):$(id -g)" --platform linux/arm64 -v ".:/app" -w /app brainser-linuxgo:arm64 make linux-arm-newgui
	docker run --rm --user "$(id -u):$(id -g)" --platform linux/amd64 -v ".:/app" -w /app brainser-linuxgo:amd64 make linux-amd64-newgui