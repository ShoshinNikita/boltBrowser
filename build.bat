@ECHO off
set GOPATH=E:\\Desktop\\Repository\\boltbrowser
ECHO Errors:
C:\\Users\\Admin\\go\\bin\\megacheck converters db web
go vet converters db web

SET /p command=Should continue(y/n): 
IF "%command%" == "y" (
	:: Build for windows
	go build -o program/boltBrowser.exe src/main.go

	:: Build for linux
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o program/boltBrowser src/main.go
)
