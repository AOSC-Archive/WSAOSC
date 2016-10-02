package main

import (
	"os"
	"path"
	"runtime"
	//"fmt"
	"log"
	"strconv"
	//"github.com/apex/log"

	"golang.org/x/sys/windows/registry"
)

func DetectDevMode() bool {
	DevKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("ERROR Querying Registry: %s\n", err)
		return false
	}
	defer DevKey.Close()
	DevMode, _, err := DevKey.GetIntegerValue("AllowDevelopmentWithoutDevLicense")
	if err != nil {
		log.Fatalf("ERROR Querying Value: %s\n", err)
		return false
	}
	log.Printf("REG: AllowDevelopmentWithoutDevLicense: %d\n", DevMode)
	return DevMode != 0
}

func GetWindow10Version() int64 {
	OSVersionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("ERROR Reading Windows 10 VERSION: %s", err)
	}
	defer OSVersionKey.Close()
	OSVersion, _, err := OSVersionKey.GetStringValue("CurrentBuildNumber")
	if err != nil {
		log.Fatalf("ERROR Reading CurrentBuildNumver: %s", err)
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
