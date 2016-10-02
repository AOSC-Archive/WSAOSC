package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
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
	/*
		purgelx = exec.Command("rd", "/q /s", path.Join(os.Getenv("localappdata"), "lxss/"))
		err = purgelx.Run()
		if err != nil {
			log.Fatalf("ERROR running rd: %s", err)
		}
	*/
	//log.Printf("lxrun stdout: %s", out.String())
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

func LxCmd(Cmd string) string {
	Cmd = "\"" + Cmd + "\""
	return "bash -c " + Cmd
	/*
		LxCmd := exec.Command("cd", path.Join(os.Getenv("localappdata"), "lxss"), "&", "bash -c", Cmd)
		err := LxCmd.Run()
		if err != nil {
			log.Fatalf("LxCmd Execution Err: %s", err)
		}*/
}

func ExtractBaseTarbal() {
	log.Printf("Start Extracting AOSC Base RootFS ...")
	//XTar := exec.Command("cd", path.Join(os.Getenv("localappdata"), "lxss/"), "&", LxCmd("mkdir -p rootfs-aosc && mv aosc.tar.xz rootfs-aosc && cd rootfs-aosc && tar -xvf aosc.tar.xz"))
	XTar := exec.Command(LxCmd("cd /root &&" +
		"rm -f .bashrc &&" +
		"tar -xvf aosc.tar.xz &&" +
		""))
	err := XTar.Run()
	if err != nil {
		log.Fatalf("Error while trying to extract AOSC Tarbal")
	}
	UpdateInstallProgress(40)
	log.Printf("Moving Home folder ...")
	XPostTar := exec.Command("cd", path.Join(os.Getenv("localappdata"), "lxss/"), "&",
		"move rootfs rootfs-ubuntu", "&",
		"move root rootfs", "&",
		"copy rootfs\\root root")
	err = XPostTar.Run()
	if err != nil {
		log.Fatalf("Error while trying to moving rootfs")
	}
	UpdateInstallProgress(80)
	Install4()
}
