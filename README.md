# WSAOSC-ng
AOSC OS on WSL (Windows Subsystem for Linux)

# Get Started
1. Enable WSL. (https://aka.ms/wslinstall).
2. Grab a WSAOSC binary from [releases](https://github.com/AOSC-Dev/WSAOSC/releases).
3. Grab AOSC OS kernel-less tarball. (https://aosc.io/os-download, we recommend the "Container" variant).
4. Decompress .tar.xz to get .tar file. (Using 7-Zip or xz command).
5. Rename .tar file to `install.tar.gz` and move it to aosc-os.exe path.
6. Run aosc-os.exe. AOSC OS will install in the same path.

# Compile Manually
## Windows
Open `WSAOSC.sln` in Visual Studio and compile.

or

Open VS command line tools, cd to WSAOSC directory and run: `msbuild WSAOSC.sln /p:Configuration=Release`

(change `Release` to `Debug` if you need debug version.)

## Linux
Install [mingw-w64](http://mingw-w64.org), cd to WSAOSC directory and run: `x86_64-w64-mingw32-gcc main.c -lole32 -O3`

# Known issue
- I/O performance is not ideal. Disabling Windows Defender and Windows Indexing Service for the rootfs directory may mitigate this issue.

# License
GNU General Public License V3
