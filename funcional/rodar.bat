@echo off
cd /d "%~dp0"
runghc Main.hs %*
exit /b %ERRORLEVEL%

