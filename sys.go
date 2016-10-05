package main

import (
	//	"bytes"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	//"fmt"
	"log"
	"strconv"
	//"github.com/apex/log"

	"golang.org/x/sys/windows/registry"
)

func DetectDevMode() bool {
	DevKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock`, registry.QUERY_VALUE)
	if err != nil {
		//log.Fatalf("ERROR Querying Registry: %s\n", err)
		ErrMsg("ERROR Querying Registry", "ERROR Querying Registry: %s", err)
		return false
	}
	defer DevKey.Close()
	DevMode, _, err := DevKey.GetIntegerValue("AllowDevelopmentWithoutDevLicense")
	if err != nil {
		//log.Fatalf("ERROR Querying Value: %s\n", err)
		ErrMsg("ERROR Querying Value", "ERROR Querying Value: %s", err)
		return false
	}
	log.Printf("REG: AllowDevelopmentWithoutDevLicense: %d\n", DevMode)
	return DevMode != 0
}

func GetWindow10Version() int64 {
	OSVersionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		//log.Fatalf("ERROR Reading Windows 10 VERSION: %s", err)
		ErrMsg("ERROR Reading Windows 10 Version", "ERROR Reading Windows 10 Version: %s", err)
	}
	defer OSVersionKey.Close()
	OSVersion, _, err := OSVersionKey.GetStringValue("CurrentBuildNumber")
	if err != nil {
		//log.Fatalf("ERROR Reading CurrentBuildNumver: %s", err)
		ErrMsg("ERROR Reading CurrentBuildNumver", "ERROR Reading CurrentBuildNumver: %s", err)
	}
	log.Printf("REG: CurrentBuildNumber: %s", OSVersion)
	OSVerInt, _ := strconv.ParseInt(OSVersion, 10, 0)
	return OSVerInt
}

func GetWindows10Edition() string {
	OSEditionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("ERROR Reading Windows 10 Edition: %s", err)
	}
	defer OSEditionKey.Close()
	OSEdition, _, err := OSEditionKey.GetStringValue("EditionID")
	if err != nil {
		log.Fatalf("ERROR Reading EditionID: %s", err)
	}
	log.Printf("REG: EditionID: %s", OSEdition)
	return OSEdition
}

func DetectUAC() {
	// UAC added to manifest
	log.Printf("UAC Access: Granted")
}

func GetGoArch() string {
	log.Printf("Runtime.GOARCH: %s", runtime.GOARCH)
	return runtime.GOARCH
}

func EnableDevMode() {
	DevKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock`, registry.SET_VALUE)
	if err != nil {
		log.Fatalf("ERROR Setting Registry: %s\n", err)
	}
	defer DevKey.Close()
	err = DevKey.SetDWordValue("AllowDevelopmentWithoutDevLicense", 1)
	if err != nil {
		log.Fatalf("ERROR Setting Value: %s\n", err)
	}
}

func DetectInstalledWSL() bool {
	if _, err := os.Stat(path.Join(os.Getenv("localappdata"), "lxss/sha256")); err == nil {
		log.Printf("ERROR: Already found an existing install of LXSS")
		return true
	}
	log.Printf("Existing install of LXSS not found")
	return false
}

func DetectInstalledRootfs() bool {
	if _, err := os.Stat(path.Join(os.Getenv("localappdata"), "lxss/root/.bashrc")); err == nil {
		//WarnMsg("ERROR: Install MS Basic RootFS Failed")
		log.Printf("RootFS found under lxss")
		return true
	}
	log.Printf("ERROR: Cannot download MS Basic RootFS")
	return false
}

func MSPathToWSL(path string) string {
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, ":", "", -1)
	path = strings.Replace(path, "C", "c", 1) //TODO: a little dirty to assume C
	path = "/mnt/" + path
	return path
}

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
