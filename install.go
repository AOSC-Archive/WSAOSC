package main

import (
	"log"
	"time"
	//	"sync"
)

func Install() {
	Install1()
}

func Install1() {
	progCurr.SetValue(0)
	progTotal.SetValue(20)
	log.Printf("Installing Basic Lx RootFS\n")
	go InstallLx()

}

func Install2() {
	if DetectInstalledRootfs() == false {
		WarnMsg("Oops", "ERROR: Install MS Basic RootFS Failed! Exiting ...")
		log.Fatalf("Exit due to Basic RootFS Installation Error")
	} else {
		log.Printf("Basic Lx RootFS Successfully Installed")
		progCurr.SetValue(100)
		progTotal.SetValue(30)
	}
	log.Printf("Downloading AOSC OS Base Tarbal ...")

	//var wg sync.WaitGroup
	//wg.Add(1)
	go Download()
	/*
		log.Printf("Waiting for Download ...")
		<-dldone
		log.Printf("Downloading Finished ...")
	*/
}

func Install3() {
	log.Printf("This should start install3() ...")
	time.Sleep(5 * time.Second)
	go ExtractBaseTarbal()
}

func Install4() {
	InfoMsg("AOSC OS successfully installed!", "Your AOSC OS 4.0 on WSL is ready to roll.\n"+
		"Now open cmd and type bash to have a try.")
	UpdateInstallProgress(100)
}

func UpdateDownloadProgress(progress int) {
	progCurr.SetValue(progress)
	progTotal.SetValue(20 + progress*5/10)
}

func UpdateInstallProgress(progress int) {
	progCurr.SetValue(progress)
	progTotal.SetValue(70 + progress*3/10)
}
