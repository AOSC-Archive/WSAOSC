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

## Compile Manually
##### under Unix or WSL
```bash
# Assume Golang environment is setup perfectly
curl https://glide.sh/get | sh # install glide package manager
go get github.com/akavel/rsrc # download rsrc to embed manifest
git clone https://github.com/AOSC-Dev/WSAOSC
cd WSAOSC
glide install
rsrc -manifest WSAOSC.exe.manifest -ico aosc.ico WSAOSC.syso
env GOOS=windows GOARCH=amd64 go build
```

**Though the compilation can go under Unix, but will produce Windows/amd64 binary only**

## License
Gnu Public License V3
