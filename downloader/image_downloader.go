package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

const imageDownloadDir = "./tmp"

type ImageDownloader struct{}

func DownloadImages(url string, downloader IDownloader) {
	cleanTmpDirectory(imageDownloadDir)

	var urlList []string
	c := colly.NewCollector()

	// Find and visit all image links
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		url := e.Attr("src")
		urlList = append(urlList, url)
	})

	c.OnScraped(func(r *colly.Response) {
		downloadableUrl := filterUrl(urlList)
		var wg sync.WaitGroup
		for i, url := range downloadableUrl {
			wg.Add(1)
			go func(url string, i int) {
				defer wg.Done()
				downloader.download(url, "test"+strconv.Itoa(i))
			}(url, i)
		}
		wg.Wait()
		fmt.Println("All images downloaded!")
	})

	c.Visit(url)
}

func cleanTmpDirectory(directory string) {
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		// If the directory exists, remove it
		err := os.RemoveAll(directory)
		if err != nil {
			log.Println("Error removing tmp directory:", err)
			return
		}
	}
}

func filterUrl(urlList []string) []string {
	var filteredUrlList []string
	for _, url := range urlList {
		if strings.Contains(url, "jpg") || strings.Contains(url, "jpeg") {
			filteredUrlList = append(filteredUrlList, strings.TrimSpace(url))
		}
	}
	return filteredUrlList
}

func (i *ImageDownloader) download(imageUrl string, fileName string) {
	// Validate and parse the image URL
	parsedURL, err := url.ParseRequestURI(imageUrl)
	if err != nil {
		log.Println("Invalid URL:", err)
		return
	}

	// Check that the URL scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		log.Println("Unsupported URL scheme:", parsedURL.Scheme)
		return
	}

	// Perform a HEAD request to verify the URL and its content type
	resp, err := http.Head(parsedURL.String())
	if err != nil {
		log.Println("Error checking URL:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the URL is reachable and the content type is an image
	if resp.StatusCode != http.StatusOK {
		log.Println("URL is not reachable, status code:", resp.StatusCode)
		return
	}
	if !isImageContentType(resp.Header.Get("Content-Type")) {
		log.Println("URL does not point to an image, Content-Type:", resp.Header.Get("Content-Type"))
		return
	}

	// Ensure the tmp directory exists
	if _, err := os.Stat(imageDownloadDir); os.IsNotExist(err) {
		err := os.Mkdir(imageDownloadDir, os.ModePerm)
		if err != nil {
			log.Println("Error creating directory:", err)
			return
		}
	}

	// Download the image
	response, err := http.Get(parsedURL.String())
	if err != nil {
		log.Println("Error downloading image:", err)
		return
	}
	defer response.Body.Close()

	// Create the file
	fileName = fmt.Sprintf("%s.jpg", fileName)
	file, err := os.Create(filepath.Join(imageDownloadDir, fileName))
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Save the image
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("Error saving image:", err)
		return
	}

	fmt.Println("Downloaded:", fileName)
}

func isImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}
