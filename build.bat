@ECHO off
set GOPATH=E:\\Desktop\\Repository\\boltbrowser-github
ECHO Errors:
C:\\Users\\Admin\\go\\bin\\megacheck converters db web
go vet converters db web

SET /p command=Should continue(y/n): 
IF "%command%" == "y" (
	:: Build for windows
	go build -o program/boltBrowser_v1.0.exe src/main.go

	:: Build for linux
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o program/boltBrowser_v1.0 src/main.go
)
