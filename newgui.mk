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
	CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w -H=windowsgui" -o bin/binary-windows-arm64.exe ./src/gui
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w -H=windowsgui" -o bin/binary-windows-amd64.exe ./src/gui

linux-arm-newgui:
	CGO_ENABLED=1 \
	GOARCH=arm64 \
	GOOS=linux \
	go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w" -o bin/binary-linux-arm64 ./src/gui

linux-amd64-newgui:
	CGO_ENABLED=1 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -tags "${BUILD_TAGS}" -trimpath -ldflags "-s -w" -o bin/binary-linux-amd64 ./src/gui

linux-docker-newgui:
	docker compose build
	docker run --rm --user "$(id -u):$(id -g)" --platform linux/arm64 -v ".:/app" -w /app test-linux:arm64 make linux-arm-newgui
	docker run --rm --user "$(id -u):$(id -g)" --platform linux/amd64 -v ".:/app" -w /app test-linux:amd64 make linux-amd64-newgui