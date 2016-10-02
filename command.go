package main

import (
	"fmt"
	"os/exec"
)

func ExecCommand(CmdString string) (string, error) {
	//Cmd := exec.Command(CmdString)
	Cmd := exec.Command("REG", `query "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" /v "AllowDevelopmentWithoutDevLicense"`)
	CmdOut, err := Cmd.Output()
	fmt.Printf("[DEBUG]: %s", CmdOut)
	if err != nil {
		return "", err
	}
	return string(CmdOut[:]), nil
}
