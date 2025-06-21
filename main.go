package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	extensionID := flag.String("id", "", "Unique extension ID")
	flag.Parse()
	baseURL := "https://clients2.google.com/service/update2/crx"

	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Failed to parse URL.")
		panic(err)
	}

	params := u.Query()
	params.Add("response", "redirect")
	params.Add("acceptformat", "crx3,puff")
	params.Add("prodversion", "137.0.7151")
	params.Add("x", fmt.Sprintf("id=%s&installsource=ondemand&uc", *extensionID))

	u.RawQuery = params.Encode()
	finalURL := u.String()

	fetchAndSave(finalURL)
}

func fetchAndSave(url string) {
	client := &http.Client{
		// Follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to make new request.")
		panic(err)
	}

	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to get response.")
		panic(err)
	}

	defer res.Body.Close()

	if res.Header.Get("Content-Type") != "application/x-chrome-extension" {
		fmt.Println("Response was not an extension file")
		panic("invalid response")
	}

	outFile, err := os.Create("temp.crx")
	if err != nil {
		fmt.Println("Failed to create download file.")
		res.Body.Close()
		panic(err)
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		fmt.Println("Failed to save to file.")
		outFile.Close()
		panic(err)
	}

	fmt.Println("Saved extension successfully.")
}
