package main

import (
	"fmt"
	"time"

	"github.com/islamjaidul/go-image-downloader/downloader"
)

func main() {
	// Start the timer
	startTime := time.Now()

	// Call the DownloadImages function from the downloader package
	downloader.DownloadImages("https://www.rightmove.co.uk/")

	// Calculate and print the total download time
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total download time: %s\n", elapsedTime)
}
