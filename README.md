# WSAOSC
(WIP) AOSC OS on WSL

## Feature
- Auto-Detect System Capability
- Auto-fix some of the options(like dev-mode)
- Multi-thread downloading the rootfs tarbal
- Automatically using ubuntu rootfs as bootstrap to untar
- Two-Click Install :)

## Get Started
1. Grab a WSAOSC binary
2. Double Click on WSAOSC.exe
3. Click Detect
4. If everything goes well, Click Install
5. After Installation, run bash.exe to start AOSC!

## Compile
##### Using HydroDev
```
HydroDev # or HydroDev build
```

##### Old-fashioned Way
```
curl https://glide.sh/get | sh
glide init && glide install
env GOOS=windows GOARCH=amd64 go build
```

**These two works for both Unix and Windows, but will only produce Windows/amd64 binary**

## License
Gnu Public License V3
