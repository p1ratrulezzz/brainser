APPID=com.jetbrainser.gui
APPNAME=Jetbrainser GUI
VERSION=$(shell git describe --tags --always)
BUILD_TAGS=guinew
LD_FLAGS=-X main.Version=${VERSION} -X main.BuildNumber=${BUILD_ID}
UID=$(shell id -u)
GID=$(shell id -g)

clean-newgui:
	rm -rf bin/osx
	rm -f bin/*
macos-newgui-arm:
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

macos-newgui-amd64:
	rm -rf bin/osx
	rm -rf dist/darwin
	mkdir -p ./dist/darwin
	mkdir -p ./bin/osx
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
	go build -tags "${BUILD_TAGS}" -buildvcs=false -trimpath -ldflags "-s -w ${LD_FLAGS}" -o bin/${BINARY_NAME}-linux-arm64 "${SRCPATH}"

linux-amd64-newgui:
	CGO_ENABLED=1 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -tags "${BUILD_TAGS}" -buildvcs=false -trimpath -ldflags "-s -w ${LD_FLAGS}" -o bin/${BINARY_NAME}-linux-amd64 "${SRCPATH}"

linux-docker-newgui:
	docker buildx build --platform linux/arm64 -t brainser-linuxgo:arm64 --output type=docker --load ./docker/linux-new
	docker buildx build --platform linux/amd64 -t brainser-linuxgo:amd64 --output type=docker --load ./docker/linux-new

	docker run --rm --user root --platform linux/amd64 -v ".:/app" -w /app brainser-linuxgo:amd64 bash -c "make linux-amd64-newgui ; chown -R ${UID}:${GID} /app/bin"
	docker run --rm --user root --platform linux/arm64 -v ".:/app" -w /app brainser-linuxgo:arm64 bash -c "make linux-arm-newgui ; chown -R ${UID}:${GID} /app/bin"
