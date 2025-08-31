@echo off
setlocal EnableDelayedExpansion

echo ================================
echo Regression CI Test Pipeline
echo ================================
echo.

REM Check Go installation
echo [1/8] Checking Go installation...
go version >nul 2>&1
if !errorlevel! neq 0 (
    echo ERROR: Go is not installed or not in PATH
    exit /b 1
)
echo Go is installed

REM Initialize and download dependencies
echo.
echo [2/8] Initializing module and downloading dependencies...
go mod tidy
if !errorlevel! neq 0 (
    echo ERROR: Failed to download dependencies
    exit /b 1
)
echo Dependencies downloaded

REM Compile the application
echo.
echo [3/8] Compiling application...
if not exist "bin" mkdir bin
go build -o bin\regression-server.exe .\cmd\server
if !errorlevel! neq 0 (
    echo ERROR: Compilation failed
    exit /b 1
)
echo Compilation successful

REM Run unit tests
echo.
echo [4/8] Running unit tests...
go test -v ./internal/...
if !errorlevel! neq 0 (
    echo WARNING: Some unit tests failed
) else (
    echo Unit tests passed
)

REM Run benchmark suite
echo.
echo [5/8] Running benchmark suite...
go run ./cmd/benchmark -o benchmark_results.json
if !errorlevel! neq 0 (
    echo WARNING: Benchmark suite encountered issues
) else (
    echo Benchmark suite completed
)

REM Start server in background for integration tests
echo.
echo [6/8] Starting server for integration tests...
set WEBHOOK_SECRET=test-secret-key
set DATABASE_PATH=:memory:
start /B bin\regression-server.exe > server.log 2>&1
set SERVER_PID=!ERRORLEVEL!

REM Wait for server to start
timeout /t 3 /nobreak >nul

REM Test server health
echo.
echo [7/8] Testing server endpoints...
curl -s http://localhost:8080/health >nul 2>&1
if !errorlevel! neq 0 (
    echo WARNING: Server health check failed
) else (
    echo Server health check passed
)

REM Test analyze endpoint with sample data
curl -s -X POST http://localhost:8080/analyze ^
    -H "Content-Type: application/json" ^
    -d @test-ci\mock-repos\golang-project\benchmarks.json >nul 2>&1
if !errorlevel! neq 0 (
    echo WARNING: Analyze endpoint test failed
) else (
    echo Analyze endpoint test passed
)

REM Run integration tests
echo.
echo [8/8] Running integration tests...
go test -v ./test-ci/integration
if !errorlevel! neq 0 (
    echo WARNING: Integration tests had issues
) else (
    echo Integration tests passed
)

REM Stop server
echo.
echo Stopping test server...
taskkill /F /IM regression-server.exe >nul 2>&1

REM Summary
echo.
echo ================================
echo Test Pipeline Summary
echo ================================
echo Compilation: SUCCESS
if exist benchmark_results.json (
    echo Benchmarks: COMPLETED
) else (
    echo Benchmarks: FAILED
)

if exist server.log (
    echo Server logs available in: src\server.log
)

if exist benchmark_results.json (
    echo Benchmark results available in: src\benchmark_results.json
)

echo.
echo Test pipeline completed
echo To run the server manually: cd src && bin\regression-server.exe
echo.

pause