@if "%~1"=="" @echo BHD v1.0, build script v1.0
@echo Building Win64 executable
@set GOHOSTARCH=amd64
@set GOHOSTOS=windows
@gomobile bind -target="android/arm64,android/amd64" -o ../build/libbhd.aar
@if ERRORLEVEL 0 @echo Done, output:../build/bhd_backend.dll