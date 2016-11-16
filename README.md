# WSAOSC
[![Go Report Card](https://goreportcard.com/badge/github.com/AOSC-Dev/WSAOSC)](https://goreportcard.com/report/github.com/AOSC-Dev/WSAOSC) [![Build status](https://ci.appveyor.com/api/projects/status/7rvcxtpagguel03n)](https://ci.appveyor.com/project/LER0ever/wsaosc)
(WIP) AOSC OS on WSL

## Feature
- Auto-Detect System Capability
- Auto-fix some of the options(like dev-mode)
- Multi-thread downloading the rootfs tarball
- Automatically using ubuntu rootfs as bootstrap to untar
- Two-Click Install :)

## Get Started
1. Grab a WSAOSC binary from CI: [WSAOSC.zip](https://ci.appveyor.com/api/projects/LER0ever/WSAOSC/artifacts/WSAOSC.zip) 
2. Enable Windows Subsystem for Linux and Developer Mode(restart needed)
3. Double Click on WSAOSC.exe

#### Default Install (latest AOSC container tarball)
4. Click Detect
5. If everything goes well, Click Install
6. After Installation, run bash.exe to start AOSC!

#### Customized Install 
4. Hold `Alt` while click Detect
5. If everything goes well, continue to 6.
6. Choose a tarball in the combo box or type third-party tarbal URL.
7. Click Install and it will handle the rest.

Currently functioning tarball: 
- AOSC (default)
- Archlinux, Deepin, Manjaro (you may need to do a unpack-and-repack stuff)

## Compile Manually
##### under Unix or WSL
```bash
# Assume Golang environment is set up properly
curl https://glide.sh/get | sh # install glide package manager
go get github.com/akavel/rsrc # download rsrc to embed manifest
git clone https://github.com/AOSC-Dev/WSAOSC
cd WSAOSC
glide install
rsrc -manifest WSAOSC.exe.manifest -ico aosc.ico WSAOSC.syso
env GOOS=windows GOARCH=amd64 go build
```
##### under Windows (64 bit only)
1. Set up go dev environment properly
2. Download glide from [glide.sh](https://glide.sh) and put it into PATH
3. ```git clone https://github.com/AOSC-Dev/WSAOSC.git # to %GOPATH%/src/github.com/AOSC-Dev/WSAOSC.git```
4. Enter the WSAOSC directory
5. ```glide install```
6. ```rsrc -manifest WSAOSC.exe.manifest -ico AOSC.ico WSAOSC.syso```
7. **go build**

**The compilation can go under Unix, but will produce Windows/amd64 binary only**

## License
GNU General Public License V3
