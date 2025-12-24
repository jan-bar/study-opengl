@echo off
set CGO_ENABLED=1
set GOARCH=amd64
set GOOS=windows
set GOAMD64=v3
if exist "%~1\" (
	echo Build %1% ...
	go build -C %1% -ldflags "-s -w -buildid=" -buildvcs=false -trimpath -o ../../main.exe
) else (
	echo Skip Build ...
)
.\main.exe
