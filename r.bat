@echo off

setlocal

if exist r.bat goto ok
echo r.bat must be run from its folder
goto end

: ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

::gofmt -w src

go install server
go install client

:end
echo finished