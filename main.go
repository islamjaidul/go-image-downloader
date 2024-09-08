package main

import (
	"fmt"
	"time"

	"github.com/islamjaidul/go-image-downloader/downloader"
)

func main() {
	// Start the timer
	startTime := time.Now()

	var d downloader.IDownloader = &downloader.ImageDownloader{}

	// Call the DownloadImages function from the downloader package
	downloader.DownloadImages("hello", d)

	// Calculate and print the total download time
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total download time: %s\n", elapsedTime)
}
