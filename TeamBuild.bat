set GOOS=windows
set GOARCH=amd64
go build  -o output\TeamServer\TeamServer.exe -ldflags="-w -s" .\cmd\Teamserver\server.go

REM call build.bat


bash -c "bash buildTeam.sh"