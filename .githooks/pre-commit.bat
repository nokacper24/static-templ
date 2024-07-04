@echo off
SETLOCAL

REM Read the current version from the .version file
set versionFile=.version
for /f "tokens=*" %%i in (%versionFile%) do set currentVersion=%%i

REM Split the version number into its components
for /f "tokens=1-3 delims=." %%a in ("%currentVersion%") do (
    set major=%%a
    set minor=%%b
    set patch=%%c
)

REM Increment the patch version number
set /a patch=patch+1

REM Combine the version parts back into a version string
set newVersion=%major%.%minor%.%patch%

REM Write the new version to the .version file
echo %newVersion% > %versionFile%

REM Stage the updated .version file
git add %versionFile%

ECHO Updated version to %newVersion% and staged %versionFile% for commit.

ENDLOCAL
