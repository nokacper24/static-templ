@echo off

REM Check if the .version file has been updated and needs to be committed
git diff --cached --name-only | findstr /R /C:"^\.version$" >nul
if %ERRORLEVEL% EQU 0 (
    git commit --amend --no-edit
)
