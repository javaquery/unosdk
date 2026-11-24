@echo off
rem unosdk shell wrapper for cmd.exe
rem Runs the real binary (core\unosdk.exe) and then applies any environment
rem changes it requested to THIS shell, so no terminal restart is needed.
rem
rem From another .bat file, call it with:  call unosdk switch java openjdk 21

set "UNOSDK_ENV_FILE=%TEMP%\unosdk_%RANDOM%%RANDOM%.cmd"
"%~dp0core\unosdk.exe" %*
set "UNOSDK_RC=%ERRORLEVEL%"
if exist "%UNOSDK_ENV_FILE%" (
    call "%UNOSDK_ENV_FILE%"
    del "%UNOSDK_ENV_FILE%" >nul 2>&1
)
set "UNOSDK_ENV_FILE=" & set "UNOSDK_RC=" & exit /b %UNOSDK_RC%