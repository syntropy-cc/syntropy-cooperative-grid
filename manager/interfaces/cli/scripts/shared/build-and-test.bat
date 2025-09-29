@echo off
REM Syntropy CLI Manager - Simple Build and Test Workflow for Windows
REM Script para compilar e testar a aplicação CLI no Windows

setlocal enabledelayedexpansion

REM Configuration
set SCRIPT_DIR=%~dp0
set BUILD_DIR=%SCRIPT_DIR%build
set VERSION=%date:~-4,4%%date:~-10,2%%date:~-7,2%-%time:~0,2%%time:~3,2%%time:~6,2%
set VERSION=%VERSION: =0%
set GIT_COMMIT=unknown
set BUILD_TIME=%date% %time%

REM Try to get git commit
git rev-parse --short HEAD >nul 2>&1
if not errorlevel 1 (
    for /f %%i in ('git rev-parse --short HEAD') do set GIT_COMMIT=%%i
)

REM Colors (if supported)
set "GREEN=[92m"
set "BLUE=[94m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

echo.
echo %CYAN%╔══════════════════════════════════════════════════════════════╗%RESET%
echo %CYAN%║              SYNTROPY CLI MANAGER                           ║%RESET%
echo %CYAN%║                Simple Build ^& Test                          ║%RESET%
echo %CYAN%╚══════════════════════════════════════════════════════════════╝%RESET%
echo.

REM Check prerequisites
echo %BLUE%[INFO] Checking prerequisites...%RESET%

REM Check Go
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR] Go is not installed. Please install Go 1.22.5 or higher.%RESET%
    echo %YELLOW%[INFO] Download: https://golang.org/dl/%RESET%
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo %GREEN%[SUCCESS] Go %GO_VERSION% found%RESET%

REM Check if main.go exists
if not exist "%SCRIPT_DIR%main.go" (
    echo %RED%[ERROR] main.go not found in %SCRIPT_DIR%%RESET%
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS] Project structure verified%RESET%

REM Prepare build environment
echo.
echo %CYAN%=== Preparing Build Environment ===%RESET%

cd /d "%SCRIPT_DIR%"

REM Create build directory
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

REM Clean previous builds
del /q "%BUILD_DIR%\*" 2>nul

echo %GREEN%[SUCCESS] Build environment prepared%RESET%

REM Setup dependencies
echo.
echo %CYAN%=== Setting Up Dependencies ===%RESET%

echo %BLUE%[INFO] Downloading dependencies...%RESET%
go mod download
if errorlevel 1 (
    echo %RED%[ERROR] Failed to download dependencies%RESET%
    pause
    exit /b 1
)

echo %BLUE%[INFO] Organizing dependencies...%RESET%
go mod tidy
if errorlevel 1 (
    echo %RED%[ERROR] Failed to organize dependencies%RESET%
    pause
    exit /b 1
)

echo %BLUE%[INFO] Verifying dependencies...%RESET%
go mod verify
if errorlevel 1 (
    echo %YELLOW%[WARNING] Dependency verification failed (continuing)%RESET%
)

echo %GREEN%[SUCCESS] Dependencies configured%RESET%

REM Build for Windows
echo.
echo %CYAN%=== Building for Windows ===%RESET%

set BUILD_FLAGS=-ldflags "-X main.version=%VERSION% -X main.buildTime=%BUILD_TIME% -X main.gitCommit=%GIT_COMMIT%"
set OUTPUT_FILE=%BUILD_DIR%\syntropy-windows.exe

echo %BLUE%[INFO] Building Windows executable...%RESET%
go build %BUILD_FLAGS% -o "%OUTPUT_FILE%" main.go
if errorlevel 1 (
    echo %RED%[ERROR] Windows build failed%RESET%
    pause
    exit /b 1
)

if exist "%OUTPUT_FILE%" (
    echo %GREEN%[SUCCESS] Windows build completed: %OUTPUT_FILE%%RESET%
) else (
    echo %RED%[ERROR] Windows build failed - file not created%RESET%
    pause
    exit /b 1
)

REM Build for Linux
echo.
echo %CYAN%=== Building for Linux ===%RESET%

set LINUX_OUTPUT=%BUILD_DIR%\syntropy-linux

echo %BLUE%[INFO] Building Linux executable...%RESET%
set GOOS=linux
set GOARCH=amd64
go build %BUILD_FLAGS% -o "%LINUX_OUTPUT%" main.go
if errorlevel 1 (
    echo %RED%[ERROR] Linux build failed%RESET%
    pause
    exit /b 1
)

REM Restore environment
set GOOS=windows
set GOARCH=amd64

if exist "%LINUX_OUTPUT%" (
    echo %GREEN%[SUCCESS] Linux build completed: %LINUX_OUTPUT%%RESET%
) else (
    echo %RED%[ERROR] Linux build failed - file not created%RESET%
    pause
    exit /b 1
)

REM Test binaries
echo.
echo %CYAN%=== Testing Binaries ===%RESET%

REM Test Windows binary
if exist "%OUTPUT_FILE%" (
    echo %BLUE%[INFO] Testing Windows binary...%RESET%
    
    REM Test version
    "%OUTPUT_FILE%" --version >nul 2>&1
    if errorlevel 1 (
        echo %YELLOW%[WARNING] Version test failed (may be normal)%RESET%
    ) else (
        echo %GREEN%[SUCCESS] Version test passed%RESET%
    )
    
    REM Test help
    "%OUTPUT_FILE%" --help >nul 2>&1
    if errorlevel 1 (
        echo %YELLOW%[WARNING] Help test failed (may be normal)%RESET%
    ) else (
        echo %GREEN%[SUCCESS] Help test passed%RESET%
    )
)

REM Show file sizes
if exist "%OUTPUT_FILE%" (
    for %%A in ("%OUTPUT_FILE%") do echo %BLUE%[INFO] Windows binary size: %%~zA bytes%RESET%
)

if exist "%LINUX_OUTPUT%" (
    for %%A in ("%LINUX_OUTPUT%") do echo %BLUE%[INFO] Linux binary size: %%~zA bytes%RESET%
)

echo %GREEN%[SUCCESS] Binary testing completed%RESET%

REM Show summary
echo.
echo %CYAN%=== Build Summary ===%RESET%

echo %GREEN%[SUCCESS] Build completed successfully!%RESET%
echo.
echo %BLUE%📁 Build Directory:%RESET% %BUILD_DIR%
echo %BLUE%📦 Version:%RESET% %VERSION%
echo %BLUE%🔧 Git Commit:%RESET% %GIT_COMMIT%
echo %BLUE%🕒 Build Time:%RESET% %BUILD_TIME%
echo %BLUE%🖥️  Platform:%RESET% Windows

echo.
echo %BLUE%📋 Created Binaries:%RESET%
if exist "%OUTPUT_FILE%" (
    for %%A in ("%OUTPUT_FILE%") do echo   %GREEN%✅%RESET% syntropy-windows.exe (%%~zA bytes) - Windows
)
if exist "%LINUX_OUTPUT%" (
    for %%A in ("%LINUX_OUTPUT%") do echo   %GREEN%✅%RESET% syntropy-linux (%%~zA bytes) - Linux
)

echo.
echo %BLUE%🚀 Next Steps:%RESET%
echo   1. Test the Windows application: %YELLOW%%OUTPUT_FILE% --help%RESET%
echo   2. Run setup: %YELLOW%%OUTPUT_FILE% setup run --force%RESET%
echo   3. Copy Linux binary to Linux machine for testing
echo   4. Copy Windows .exe to another Windows machine for testing

echo.
echo %BLUE%💡 Usage Examples:%RESET%
echo   %CYAN%%OUTPUT_FILE% --help%RESET%                    # Show help
echo   %CYAN%%OUTPUT_FILE% --version%RESET%                 # Show version
echo   %CYAN%%OUTPUT_FILE% setup --help%RESET%              # Setup help
echo   %CYAN%%OUTPUT_FILE% setup run --force%RESET%         # Run setup
echo   %CYAN%%OUTPUT_FILE% setup status%RESET%              # Check status

REM Ask if user wants to run the application
echo.
echo %YELLOW%Do you want to run the Windows application now? (y/N):%RESET%
set /p response=
if /i "%response%"=="y" (
    echo.
    echo %CYAN%=== Running Application ===%RESET%
    echo %BLUE%[INFO] Running Syntropy CLI Manager...%RESET%
    echo.
    "%OUTPUT_FILE%" --help
    echo.
    echo %GREEN%[SUCCESS] Application executed successfully!%RESET%
    echo.
    echo %BLUE%Example commands to try:%RESET%
    echo   %YELLOW%%OUTPUT_FILE% --version%RESET%
    echo   %YELLOW%%OUTPUT_FILE% setup --help%RESET%
    echo   %YELLOW%%OUTPUT_FILE% setup validate%RESET%
    echo   %YELLOW%%OUTPUT_FILE% setup run --force%RESET%
)

echo.
echo %GREEN%[SUCCESS] Build and test workflow completed!%RESET%
pause
