package main

import (
	"fmt"
	"log"
	//	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mw *walk.MainWindow
var inTE, statusTE *walk.TextEdit
var btDetect, btInstall, btAbout *walk.PushButton
var progCurr, progTotal *walk.ProgressBar
var lv *LogView
var debug bool

func ui() {
	debug = false
	if err := (MainWindow{
		Title:    "WSAOSC - AOSC OS on WSL | Installer",
		AssignTo: &mw,
		MinSize:  Size{600, 400},
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				MaxSize: Size{9999, 50},
				Children: []Widget{
					PushButton{
						Text:     "Detect Your System",
						AssignTo: &btDetect,
						OnClicked: func() {
							Prepare()
						},
					},
					PushButton{
						Text:     "Install !",
						Enabled:  false,
						AssignTo: &btInstall,
						OnClicked: func() {
							btDetect.SetEnabled(false)
							btInstall.SetEnabled(false)
							log.Printf("This should start installing\n")
							Install()
						},
					},
					PushButton{
						Text:     "About",
						AssignTo: &btAbout,
						OnClicked: func() {
							log.Printf("This should start about\n")
							walk.MsgBox(mw,
								"About WSAOSC",
								ABOUT_WSAOSC,
								walk.MsgBoxIconInformation)
						},
					},
				},
			},
			Composite{
				MaxSize: Size{9999, 10},
				Layout:  Grid{Columns: 10},
				//StretchFactor: 10,
				Children: []Widget{
					Label{
						Text: "Current Progress",
						//Size: Size{20, 10},
					},
					ProgressBar{
						AssignTo: &progCurr,
					},
				},
			},
			Composite{
				MaxSize: Size{9999, 10},
				Layout:  Grid{Columns: 10},
				//StretchFactor: 10,
				Children: []Widget{
					Label{
						Text: "Total Progress",
						//Size: Size{20, 10},
					},
					ProgressBar{
						AssignTo: &progTotal,
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lv, err := NewLogView(mw)
	if err != nil {
		log.Fatal(err)
	}

	lv.PostAppendText("===== Welcome to WSAOSC Installer! =====\n\n")
	log.SetPrefix(LOG_PREFIX)
	log.SetFlags(log.Ltime)
	if debug != true {
		log.SetOutput(lv)
	}

	mw.Run()

}

func AskMsg(title string, text string) bool {
	response := walk.MsgBox(mw, title, text, walk.MsgBoxOKCancel)
	log.Printf("AskMsg User Choice: %d", response)
	return response == 1
}

func WarnMsg(title string, text string) {
	walk.MsgBox(mw, title, text, walk.MsgBoxIconWarning)
}

func InfoMsg(title string, text string) {
	walk.MsgBox(mw, title, text, walk.MsgBoxIconInformation)
}

func ErrMsg(text string, v ...interface{}) {
	CompMsg := fmt.Sprintf(text+":\n%s", v)
	walk.MsgBox(mw, text, CompMsg, walk.MsgBoxIconError)
	log.Fatalln(CompMsg)
}
