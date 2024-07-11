@echo off

REM Get the commit message file
set commit_msg_file=%1

REM Read the commit message from the file
setlocal enabledelayedexpansion
set "commit_message="
for /f "usebackq tokens=* delims=" %%A in ("%commit_msg_file%") do (
    set "commit_message=%%A"
    goto :endfor
)
:endfor

REM Check if the commit message matches the pattern "chore: release v<version>"
echo !commit_message! | findstr /R /C:"^chore: release v[0-9]\+\.[0-9]\+\.[0-9]\+$" >nul
if %ERRORLEVEL% NEQ 0 (
    echo "Commit message does not match 'chore: release <version>'. Skipping version update."
    exit /b 0
)

REM Extract the version from the commit message
for /f "tokens=3 delims= " %%A in ("!commit_message!") do set "new_version=%%A"

REM Write the new version to the .version file
echo %new_version% > .version

echo.
echo Bump version to %new_version%.
echo.

REM Stage the updated .version file
git add .version

endlocal
