# GiTTY
[![Go Report Card](https://goreportcard.com/badge/github.com/LetUsFlow/GiTTY)](https://goreportcard.com/report/github.com/LetUsFlow/GiTTY)

A worse version of the KiTTY SSH client with fewer features but no known security vulnerabilities.

Uses [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss) for terminal UI.

## Building

### Linux
```
go build
```

### Windows
```
windres -o resource.syso resource.rc
go build
```
The `windres` command is optional, but without it, the final executable won't have the GiTTY icon. `windres` is included in `mingw-w64-binutils`.

## License
GPL-3.0
