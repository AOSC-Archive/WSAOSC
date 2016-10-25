package main

import (
	//	"bytes"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// DetectDevMode : return true if Dev Mode is enable.
func DetectDevMode() bool {
	DevKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock`, registry.QUERY_VALUE)
	if err != nil {
		ErrMsg("ERROR Querying Registry", err)
		return false
	}
	defer DevKey.Close()
	DevMode, _, err := DevKey.GetIntegerValue("AllowDevelopmentWithoutDevLicense")
	if err != nil {
		ErrMsg("ERROR Querying Value", err)
		return false
	}
	log.Printf("REG: AllowDevelopmentWithoutDevLicense: %d\n", DevMode)
	return DevMode != 0
}

// GetWindow10Version : return the build number of windows 10.
func GetWindow10Version() int64 {
	OSVersionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		ErrMsg("ERROR Reading Windows 10 Version", err)
	}
	defer OSVersionKey.Close()
	OSVersion, _, err := OSVersionKey.GetStringValue("CurrentBuildNumber")
	if err != nil {
		ErrMsg("ERROR Reading CurrentBuildNumver", err)
	}
	log.Printf("REG: CurrentBuildNumber: %s", OSVersion)
	OSVerInt, _ := strconv.ParseInt(OSVersion, 10, 0)
	return OSVerInt
}

// GetWindows10Edition : return windows 10 edition: professional, home basic, enterprise.
func GetWindows10Edition() string {
	OSEditionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		ErrMsg("ERROR Reading Windows 10 Edition", err)
	}
	defer OSEditionKey.Close()
	OSEdition, _, err := OSEditionKey.GetStringValue("EditionID")
	if err != nil {
		ErrMsg("ERROR Reading EditionID", err)
	}
	log.Printf("REG: EditionID: %s", OSEdition)
	return OSEdition
}

// DetectUAC : To be removed
func DetectUAC() {
	log.Printf("UAC Access: Granted")
}

// GetGoArch : return the archetecture of the host machine
func GetGoArch() string {
	log.Printf("Runtime.GOARCH: %s", runtime.GOARCH)
	return runtime.GOARCH
}

// EnableDevMode : Enable Devmode by setting that balue to 1
func EnableDevMode() {
	DevKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock`, registry.SET_VALUE)
	if err != nil {
		ErrMsg("ERROR Setting Registry", err)
	}
	defer DevKey.Close()
	err = DevKey.SetDWordValue("AllowDevelopmentWithoutDevLicense", 1)
	if err != nil {
		ErrMsg("ERROR Setting Value", err)
	}
}

// DetectInstalledWSL : return true if found an previous install
func DetectInstalledWSL() bool {
	if _, err := os.Stat(path.Join(os.Getenv("localappdata"), "lxss/sha256")); err == nil {
		log.Printf("Warning: Already found an existing install of LXSS")
		return true
	}
	log.Printf("Existing install of LXSS not found")
	return false
}

// DetectInstalledRootfs : return true if basic rootfs is successfully installed
func DetectInstalledRootfs() bool {
	if _, err := os.Stat(path.Join(os.Getenv("localappdata"), "lxss/root/.bashrc")); err == nil {
		log.Printf("RootFS found under lxss")
		return true
	}
	log.Printf("ERROR: Cannot download MS Basic RootFS")
	return false
}

// MSPathToWSL : Convert the windows style path to WSL compatible path.
func MSPathToWSL(path string) string {
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, ":", "", -1)
	path = strings.Replace(path, "C", "c", 1) //TODO: a little dirty to assume C
	path = "/mnt/" + path
	return path
}

// Powershell : Run Powershell command
func Powershell(Ps string) string {

	cmd := exec.Command("Powershell.exe",
		"-NoProfile",
		"-NoLogo",
		"-ExecutionPolicy",
		"Bypass",
		"-Command", Ps)

	var Out bytes.Buffer
	var Stderr bytes.Buffer

	cmd.Stdin = os.Stdin
	cmd.Stdout = &Out
	cmd.Stderr = &Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + Stderr.String())
	}

	r := Out.String()
	return r
}

// DetectExecInPath : return if the executable is in PATH
func DetectExecInPath(executable string) bool {
	_, err := exec.LookPath(executable)
	if err != nil {
		return false
	}
	return true
}
