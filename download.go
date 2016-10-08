package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/iovxw/downloader"
)

// Download: download AOSC Base Nokernel tarbal
func Download() {
	file, err := os.Create(path.Join(os.Getenv("localappdata"), "lxss/aosc.tar.xz"))
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	fileDl, err := downloader.NewFileDl(AOSC_AMD64_TARBAL, file, -1)
	if err != nil {
		log.Println(err)
	}
	var wg sync.WaitGroup
	var exit = make(chan bool)
	var resume = make(chan bool)
	var pause bool
	wg.Add(1)
	fileDl.OnStart(func() {
		fmt.Println("download started")
		format := "\033[2K\r%v/%v [%s] %v byte/s %v"
		for {
			status := fileDl.GetStatus()
			var i = float64(status.Downloaded) / float64(fileDl.Size) * 50
			h := strings.Repeat("=", int(i)) + strings.Repeat(" ", 50-int(i))

			select {
			case <-exit:
				fmt.Printf(format, status.Downloaded, fileDl.Size, h, 0, "[FINISH]")
				fmt.Println("\ndownload finished")
				wg.Done()
			default:
				if !pause {
					time.Sleep(100 * time.Millisecond)
					go UpdateDownloadProgress(int(status.Downloaded * 100 / fileDl.Size))
					fmt.Printf(format, status.Downloaded, fileDl.Size, h, status.Speeds, "[DOWNLOADING]")
					//os.Stdout.Sync()
				} else {
					fmt.Printf(format, status.Downloaded, fileDl.Size, h, 0, "[PAUSE]")
					os.Stdout.Sync()
					<-resume
					pause = false
				}

			}
		}
	})

	fileDl.OnFinish(func() {
		exit <- true
	})

	fileDl.OnError(func(errCode int, err error) {
		log.Println(errCode, err)
	})

	fmt.Printf("%+v\n", fileDl)

	fileDl.Start()
	wg.Wait()
	Install3()
}
