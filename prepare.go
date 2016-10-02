package main

import (
	"log"
	"time"
)

var ReadyForInstall bool

func Prepare() {
	ReadyForInstall = true
	Prepare1()
}

func Prepare1() {
	progCurr.SetValue(0)
	progTotal.SetValue(0)

	// Program already under RUNAS
	{
		DetectUAC()
		log.Printf(MSG_UAC_ALREADY_ENABLED)
		UpdatePrepareProgress(10)
	}
	WinArch := GetGoArch()
	if WinArch == "amd64" {
		UpdatePrepareProgress(20)
		log.Printf("You are running 64bit Windows 10, continuing ...")
	}
	if DetectDevMode() == true {
		log.Printf(MSG_DEV_MODE_ALREADY_ENABLED)
		UpdatePrepareProgress(30)
	} else {
		if AskMsg("DevMode Required!", ASK_DEV_MODE) == true {
			EnableDevMode()
			log.Printf(MSG_DEV_MODE_JUST_ENABLED)
			UpdatePrepareProgress(30)
		} else {
			ReadyForInstall = false
			log.Fatal("Installation Canceled")
		}
	}
	Win10Ver := GetWindow10Version()
	if Win10Ver >= 14393 {
		log.Printf("Your Windows 10 Version is newer than or equal to WINTH-Build#14393, continuing ...")
		UpdatePrepareProgress(50)
	}
	Win10Ed := GetWindows10Edition()
	if Win10Ed == "Professional" {
		log.Printf("You are using a Professional Edition, continuing ...")
		UpdatePrepareProgress(70)
	} else {
		WarnMsg("Unsupported OS detected", "Do the world a favor, upgrade to Win10 Professional or just use Linux !")
	}

	LxssInstalled := DetectInstalledWSL()
	if LxssInstalled == false {
		log.Printf("You haven't installed WSL before, continuing ...")
		Prepare2()
	} else {
		if AskMsg("Previous install detected!", ASK_DEL_WSL) == false {
			ReadyForInstall = false
			WarnMsg("Warning", "The Installation is canceled.")
			log.Fatalf("Installation Canceled.")
		} else {
			log.Printf("Purging Previous WSL ...")
			go PurgeLx()
		}
		/*WarnMsg("Warning", "Previous install detected! Please remove it first. The Installation is canceled.")
		log.Fatalf("Installation Canceled.")*/
	}

	//progTotal.SetValue(20)
}

func Prepare2() {
	UpdatePrepareProgress(90)
	if ReadyForInstall == true {
		log.Printf("Preparing to install AOSC OS on WSL")
		time.Sleep(500 * time.Millisecond)
		UpdatePrepareProgress(100)
		InfoMsg("Ready to roll!", "Your computer meets all the requirements for AOSC OS on Windows!")
		btInstall.SetEnabled(true)
	}
}
func UpdatePrepareProgress(progress int) {
	progCurr.SetValue(progress)
	progTotal.SetValue(progress / 5)
}
