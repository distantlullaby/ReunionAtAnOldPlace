@echo off
chcp 65001 >nul
REM 一键启动后端(8080)与前端(5173) —— Windows
setlocal
set GOROOT=C:\mine\code\tool\Google\go\go1.22.1\sdk
set PATH=%GOROOT%\bin;%PATH%

start "memorylink-backend" cmd /k "cd /d %~dp0backend && go run ./cmd/server"
timeout /t 2 >nul
start "memorylink-frontend" cmd /k "cd /d %~dp0frontend && npm run dev"

echo 后端: http://localhost:8080
echo 前端: http://localhost:5173
endlocal
