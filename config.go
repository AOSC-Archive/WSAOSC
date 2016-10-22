package main

const (
	AOSC_AMD64_TARBAL = "https://mirror.anthonos.org/aosc-os/os-amd64/container/aosc-os_container_20161021.tar.xz"
	LOG_PREFIX        = "[WSAOSC]@"
	ABOUT_WSAOSC      = "A small tool that helps you install " +
		"AOSC on your windows using Windows Subsystem for Linux." +
		"This program is released under Affero GNU Public License." +
		"I don't know what to say here."
	MSG_UAC_ALREADY_ENABLED      = "You have proper administrative privilige, continuing ..."
	MSG_DEV_MODE_ALREADY_ENABLED = "You have already enabled Dev Mode, continuing ..."
	MSG_DEV_MODE_JUST_ENABLED    = ""
	ASK_DEV_MODE                 = "You need to enable Developer's Mode in order to " +
		"use the Windows Subsystem for Linux"
	ASK_DEL_WSL = "It seems that you have already installed the WSL before\n" +
		"Press yes to delete the whole previous install\n" +
		"And no to backup your data."
)
