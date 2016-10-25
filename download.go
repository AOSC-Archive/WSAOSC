package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iovxw/downloader"
	"golang.org/x/net/html"
)

// Download : download AOSC Base Nokernel tarball
func Download() {
	file, err := os.Create(path.Join(os.Getenv("localappdata"), "lxss/aosc.tar.xz"))
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	var fileDl *downloader.FileDl
	if CustomTarball == false {
		fileDl, err = downloader.NewFileDl(GetLatestTarballURL(), file, -1)
	} else {
		fileDl, err = downloader.NewFileDl(cbSelectTarball.Text(), file, -1)
	}

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

// GetTarballURLs returns a slice of strings of tarball URLs
func GetTarballURLs() []string {
	var RetURLs []string
	resp, _ := http.Get(AOSC_AMD64_REPO)
	//bytes, _ := ioutil.ReadAll(resp.Body)
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return RetURLs
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						fmt.Println("Found href:", a.Val)
						a.Val = AOSC_AMD64_REPO + a.Val
						if !(a.Val[len(a.Val)-9:] == "sha256sum" || a.Val[len(a.Val)-6:] == "md5sum" || a.Val[len(a.Val)-3:] == "../") {
							RetURLs = append(RetURLs, a.Val)
						}
						break
					}
				}
			}
		}
	}
	resp.Body.Close()
	fmt.Printf("RetURL size: %d", len(RetURLs))
	for _, url := range RetURLs {
		fmt.Println(url)
	}
	return RetURLs
}

// GetLatestTarballURL : Latest one
func GetLatestTarballURL() string {
	URLs := GetTarballURLs()
	sort.Strings(URLs)
	return URLs[0]
}

// FillComboTarball : fill combobox with tarball URLs
func FillComboTarball() {
	TarballURLs := GetTarballURLs()
	cbSelectTarball.SetModel(TarballURLs)
}
