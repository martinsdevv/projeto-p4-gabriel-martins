@echo off
cd /d "%~dp0"
javac -encoding UTF-8 -d out -sourcepath src src\matchmaking\Main.java
if errorlevel 1 exit /b 1
java -cp out matchmaking.Main %*
exit /b %ERRORLEVEL%
