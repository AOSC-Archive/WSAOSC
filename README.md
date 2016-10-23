# WSAOSC
[![Go Report Card](https://goreportcard.com/badge/github.com/LER0ever/WSAOSC)](https://goreportcard.com/report/github.com/LER0ever/WSAOSC)  
(WIP) AOSC OS on WSL

## Feature
- Auto-Detect System Capability
- Auto-fix some of the options(like dev-mode)
- Multi-thread downloading the rootfs tarball
- Automatically using ubuntu rootfs as bootstrap to untar
- Two-Click Install :)

## Get Started
1. Grab a WSAOSC binary
2. Enable Windows Subsystem for Linux and Developer Mode(restart needed)
3. Double Click on WSAOSC.exe
4. Click Detect
5. If everything goes well, Click Install
6. After Installation, run bash.exe to start AOSC!

## Compile
##### Using HydroDev under Unix
```bash
# Run under bash
HydroDev # or HydroDev build
```

**Though the compilation can go under Unix, but will produce Windows/amd64 binary only**

## License
Gnu Public License V3
