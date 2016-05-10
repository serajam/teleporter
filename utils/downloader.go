package utils

import (
	"os"
	"fmt"
	"net/http"
	"io"
	"log"
	"encoding/json"
	"strings"
)

const (
	DownloadUrl = "https://stream16.mixcloud.com/c/m4a/64/%s.m4a"
	DownloadDir = "data/"
)


// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PassThru struct {
	io.Reader
	total int64 // Total # of bytes transferred
}

type Details struct {
	Url string `json:"waveform_url"`
}

func Download(m *Mix) {
	// get details with download url part
	aUrl := fmt.Sprintf(DownloadUrl, getAudioUrl(m.DetailsUrl()))
	fmt.Println("Download url: ", aUrl)

	out, err := os.Create(fmt.Sprintf("%s%s.mp3", DownloadDir, m.FileName()))
	defer out.Close()

	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", aUrl, nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36")

	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	// Wrap it with our custom io.Reader.
	src := &PassThru{Reader: resp.Body}

	count, err := io.Copy(out, src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("Transfer completed successfuly. Total transfered", count, "bytes")
}

func getAudioUrl(detailsUrl string) string {
	response, err := http.Get(detailsUrl)
	if err != nil {
		fmt.Printf("%s", err)
		return ""
	} else {
		defer response.Body.Close()
		x := new(Details)
		err = json.NewDecoder(response.Body).Decode(x)
		s := strings.Split(x.Url, "com")
		s = strings.Split(s[1], ".")
		return s[0];
	}
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.total += int64(n)

	if err == nil {
		fmt.Printf("\rDownloaded %d bytes", pt.total)
	}

	return n, err
}
