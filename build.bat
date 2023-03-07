cd cmd\Client
astilectron-bundler.exe
pause
REM bash -c "cp -r output -t ../../"
REM bash -c "rm -rf output"
REM cd ..\..\
REM mkdir output\TeamServer
REM set GOOS=windows
REM go build -o output\TeamServer\TeamServer_Win.exe cmd\TeamServer\main.go

REM set GOOS=linux
REM go build -o output\TeamServer\TeamServer_Lin cmd\TeamServer\main.go

REM set GOOS=darwin
REM go build -o output\TeamServer\TeamServer_Mac cmd\TeamServer\main.go

REM Agent Obfuscation and Build Process
bash -c "cp -r cmd backupAgent"
REM Backed up source files
REM Starting Obfuscation of Agent
cd cmd\Agent
call obf-and-build.bat
cd ..\Agent.Light
call obf-and-build.bat
REM bash -c "bash build.sh"
