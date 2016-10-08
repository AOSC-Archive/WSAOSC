package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	//"path"
	//"time"
)

// PurgeLx : completely remote the lxss install
func PurgeLx() {
	purgelx := exec.Command("lxrun", "/uninstall", "/full", "/y")
	err := purgelx.Run()
	if err != nil {
		log.Fatalf("ERROR running lxrun: %s", err)
	}

	Prepare2()
}

// InstallLx : Install the basic ubuntu lxss rootfs
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

// LxCmd : Directly executes the bash command
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

// ExtractBaseTarbal : Extract tarbal with permission info and move rootdir outside
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
	log.Printf("Patching Display and Dbus ...")
	PatchDbus()
	PatchDisplay()
	Install4()
}

// PatchDbus : Make dbus work under WSAOSC, note that this hack won't work on other distro
func PatchDbus() {
	LxCmd(`echo "<listen>tcp:host=localhost,port=0</listen>" >> /etc/dbus-1/session.conf`)
	UpdateInstallProgress(87)
}

// PatchDisplay : Make xorg display work for VcXsrv
func PatchDisplay() {
	LxCmd(`echo "export DISPLAY=:0.0" >> ~/.bashrc`)
	UpdateInstallProgress(95)
}
