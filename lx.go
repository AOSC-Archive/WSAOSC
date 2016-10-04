package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	//"path"
	//"time"
)

func PurgeLx() {
	purgelx := exec.Command("lxrun", "/uninstall", "/full", "/y")
	//var out bytes.Buffer
	//urgelx.Stdout = &out
	err := purgelx.Run()
	if err != nil {
		log.Fatalf("ERROR running lxrun: %s", err)
	}

	Prepare2()
}

func InstallLx() {
	installlx := exec.Command("lxrun", "/install", "/y")
	var out bytes.Buffer
	installlx.Stdout = &out
	err := installlx.Run()
	if err != nil {
		log.Fatalf("ERROR running lxrun: %s", err)
	}
	Install2()
	//log.Printf("lxrun stdout: %s", out.String())
}

func TestLx() {
	p := exec.Command("bash.exe", "-c", "bash -version")
	p.Stdin = os.Stdin
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	e := p.Run()
	if e != nil {
		fmt.Println(e)
	}
}
func LxCmd(Cmd string) {
	p := exec.Command("bash.exe", "-c", Cmd)
	p.Stdin = os.Stdin
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	e := p.Run()
	if e != nil {
		log.Printf(e.Error())
	}
}

func ExtractBaseTarbal() {
	UpdateInstallProgress(0)
	LxCmd("bash -version")
	log.Printf("Start Extracting AOSC Base RootFS ...")
	LXDIR := path.Join(os.Getenv("localappdata"), "lxss")

	//XTar := exec.Command("cd", path.Join(os.Getenv("localappdata"), "lxss/"), "&", LxCmd("mkdir -p rootfs-aosc && mv aosc.tar.xz rootfs-aosc && cd rootfs-aosc && tar -xvf aosc.tar.xz"))
	LxCmd("mv " + MSPathToWSL(path.Join(LXDIR, "aosc.tar.xz")) + " /root")
	UpdateInstallProgress(20)
	LxCmd("cd /root && " +
		"rm -f .bashrc && " +
		"tar -xvpf aosc.tar.xz && " +
		"exit")

	UpdateInstallProgress(40)
	log.Printf("Moving Home folder ...")
	psout := Powershell("cd " + LXDIR + "; " +
		"move rootfs rootfs-ubuntu" + "; " +
		"move root rootfs" + "; " +
		"move rootfs\\root root")
	log.Printf("PS> Output %s", psout)
	UpdateInstallProgress(80)
	log.Printf("Everything Complete, Continuing ...")
	Install4()
}
