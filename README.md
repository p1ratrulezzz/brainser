# brainser

## Requirements for building
1) Docker
2) Git
3) Go 1.23 or higher
4) resources.zip file (not included here due to legal reasons, should be found somewhere else)

## Building console variant

0) Clone the project
1) Extract resources.zip (resources_enc folder) to src/app/ so it will be src/app/resources_enc
2) make build - to build console version

## Building the gui variant (advanced users)

### Simple build of gui variant:
```
CGO_ENABLED=1 go build -tags gui -o brainser-binary "jetbrainser/src/app"
```

### Fyne build

Gui version with fyne is built in a more complexs way. Cross compiling from linux only works for windows and linux. 
Compiling for macos is only available from  macos, xcode is required. 
Check other targets in makefile and try building it
You have to run docker-compose build from the project root in order to ensure local docker images are created (for linux).

Compiling from windows to other platforms and to windows is a hell though possible. msys32 can do the job, build by simple running go build with CGO_ENABLED=1 the recommended way is to use linux.

Example of building linux-x64 verison (ensure docker is running)
```
make buildgui-linux-amd64
```

## Run GUI variant without gui

```
./binary-amd64 -nogui
```

This will run the console variant


## Limitations
Some versions of this app might not work in virtual machines such as VMWare and Parallels Desktop due to lack of support for OpenGL. Use console variant for virtual machines.


# Donate

## Crypto
**Bitcoin**: 14QQYoqNsL9eVwgLRi5QgAVpgBU1jjpQ3f<br>
**Litecoin**: LUGX7Utf8zGTckDnHMBpGEftpZyVuXJPgz<br>
**ETH(ERC20/BEP20)** 0xfa97d136abf1c83b0cccc9ac3d7e82a451c87abd<br>
**BNB(BEP20)**: 0xfa97d136abf1c83b0cccc9ac3d7e82a451c87abd<br>
**TRX/TRC20**: TTWpRXAZJdz1AazKQRzsaAfHRHW6q2veSJ<br>
