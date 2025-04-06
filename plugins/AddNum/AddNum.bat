@echo off

set URI=%1
echo Current URI: %URI%

set PLUGIN_NAME=AddNum
set PLUGIN_PATH=.\%PLUGIN_NAME%.exe

if not exist "%PLUGIN_PATH%" (
    echo Error: Plugin "%PLUGIN_NAME%" not found at %PLUGIN_PATH%.
    exit /b 1
)

echo Starting plugin %PLUGIN_NAME%...
call "%PLUGIN_PATH%" %URI%

if %ERRORLEVEL% neq 0 (
    echo Start failed with error code %ERRORLEVEL%.
    exit /b %ERRORLEVEL%
)

echo Plugin started successfully.
exit /b 0