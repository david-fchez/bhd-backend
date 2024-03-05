# BHD backend project

The backend project is written in GoLang, it is used as a dependency for the BHD game.

## Build instructions
### Windows
Run the [make_win_dll_dynamic.bat](./make_win_dll_dynamic.bat) or in your terminal manually enter:
```bash
go build  -o ../build/libbhd.dll -buildmode=c-shared
```

### Linux
Run the command:
```bash
o build -o build/libbhd.so -buildmode=c-shared
```

### macOS
Run the command:
```bash
o build -o build/libbhd.dylib -buildmode=c-shared
```

### Android
`go-mobile` needs to be installed (`go get golang.org/x/mobile/cmd/gomobile`).
Run the [make_android_aar.bat](./make_android_aar.bat) or in your terminal manually enter:
```bash
gomobile bind -target="android/arm64,android/amd64" -o ../build/libbhd.aar
```
