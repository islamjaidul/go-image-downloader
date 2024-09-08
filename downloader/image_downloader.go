package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
				downloader.Download(url, "test"+strconv.Itoa(i))
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

func (i *ImageDownloader) Download(url string, fileName string) {
	response, err := http.Get(url)
	if err != nil {
		log.Println("Error downloading:", err)
		return
	}
	defer response.Body.Close()

	// Ensure the tmp directory exists
	if _, err := os.Stat(imageDownloadDir); os.IsNotExist(err) {
		err := os.Mkdir(imageDownloadDir, os.ModePerm)
		if err != nil {
			log.Println("Error creating directory:", err)
			return
		}
	}

	fileName = fmt.Sprintf("%s.jpg", fileName)
	file, err := os.Create(imageDownloadDir + "/" + fileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("Error saving image:", err)
		return
	}

	fmt.Println("Downloaded:", fileName)
}
