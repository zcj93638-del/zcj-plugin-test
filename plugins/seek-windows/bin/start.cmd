@echo off
setlocal
set "arch=%PROCESSOR_ARCHITECTURE%"
if defined PROCESSOR_ARCHITEW6432 set "arch=%PROCESSOR_ARCHITEW6432%"
if /i "%arch%"=="ARM64" goto arm64
if /i not "%arch%"=="AMD64" exit /b 1
"%~dp0seek-mcp-windows-amd64.exe" %*
exit /b %errorlevel%
:arm64
"%~dp0seek-mcp-windows-arm64.exe" %*
exit /b %errorlevel%
