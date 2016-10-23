package main

import "log"

//	"sync"

// Install:
func Install() {
	Install1()
}

// Install1: Install ubuntu LxRootfs
func Install1() {
	progCurr.SetValue(0)
	progTotal.SetValue(20)
	log.Printf("Installing Basic Lx RootFS\n")
	go InstallLx()

}

// Install2: Download AOSC Tarball
func Install2() {
	if DetectInstalledRootfs() == false {
		WarnMsg("Oops", "ERROR: Install MS Basic RootFS Failed! Exiting ...")
		log.Fatalf("Exit due to Basic RootFS Installation Error")
	} else {
		log.Printf("Basic Lx RootFS Successfully Installed")
		progCurr.SetValue(100)
		progTotal.SetValue(30)
	}
	log.Printf("Downloading AOSC OS Base Tarball ...")

	//var wg sync.WaitGroup
	//wg.Add(1)
	go Download()
	/*
		log.Printf("Waiting for Download ...")
		<-dldone
		log.Printf("Downloading Finished ...")
	*/
}

// Install3: Extract tarball to /root
func Install3() {
	log.Printf("This should start install3() ...")
	//time.Sleep(5 * time.Second)
	go ExtractBaseTarball()
	/**/
}

// Install4: Finish Installation
func Install4() {
	UpdateInstallProgress(100)
	InfoMsg("AOSC OS successfully installed!", "Your AOSC OS 4.0 on WSL is ready to roll.\n"+
		"Now open cmd and type bash to have a try.")
}

// UpdateDownloadProgress: Update two progress bar for Download
func UpdateDownloadProgress(progress int) {
	progCurr.SetValue(progress)
	progTotal.SetValue(20 + progress*5/10)
}

// UpdateInstallProgress: Update two progress bar for Install
func UpdateInstallProgress(progress int) {
	progCurr.SetValue(progress)
	progTotal.SetValue(70 + progress*3/10)
}
