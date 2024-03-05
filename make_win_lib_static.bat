@if "%~1"=="" @echo BHD v1.0, build script v1.0
@echo Building Win64 executable
@set GOHOSTARCH=amd64
@set GOHOSTOS=windows
@go build  -o ../build/libbhd.lib -buildmode=c-archive -ldflags="-extldflags=-static"
@if ERRORLEVEL 0 @echo Done, output:../build/bhd_backend.lib