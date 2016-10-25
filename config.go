package main

const (
	AOSC_AMD64_REPO = "https://mirror.anthonos.org/aosc-os/os-amd64/container/"
	LOG_PREFIX      = "[WSAOSC]@"
	ABOUT_WSAOSC    = "A small tool that helps you install " +
		"AOSC on your windows using Windows Subsystem for Linux. " +
		"This program is released under GNU General Public License. " +
		"See github.com/LER0ever/WSAOSC for details."
	MSG_UAC_ALREADY_ENABLED      = "You have proper administrative privilige, continuing ..."
	MSG_DEV_MODE_ALREADY_ENABLED = "You have already enabled Dev Mode, continuing ..."
	MSG_DEV_MODE_JUST_ENABLED    = ""
	ASK_DEV_MODE                 = "You need to enable Developer's Mode in order to " +
		"use the Windows Subsystem for Linux"
	ASK_DEL_WSL = "It seems that you have already installed the WSL before\n" +
		"Press OK to delete the whole previous install\n" +
		"And Cancel to backup your data."
)
