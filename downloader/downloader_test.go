package downloader

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDownloader struct {
	mock.Mock
}

func (m *MockDownloader) download(url string, fileName string) {
	m.Called(url, fileName)
}

func Test_Downloader(t *testing.T) {

	t.Run("cleanTmpDirectory: It should remove the given directory", func(t *testing.T) {
		err := os.Mkdir("../tmp", os.ModePerm)
		if err != nil {
			log.Println("Error creating directory:", err)
			return
		}

		cleanTmpDirectory("../tmp")
		_, err = os.Stat("../tmp")
		assert.Equal(t, os.IsNotExist(err), true)
	})

	t.Run("filterUrl: It should filter only jpg and jpeg from given slice of string", func(t *testing.T) {
		urlList := []string{
			"https://s3.aws.com/media/test1.jpg",
			"https://s3.aws.com/media/test2.jpg",
			"https://s3.aws.com/media/test3.jpeg",
			"https://s3.aws.com/media/test4.bmp",
			"https://s3.aws.com/media/test5.webp",
		}
		filteredUrlList := filterUrl(urlList)
		assert.Equal(t, len(filteredUrlList), 3)
		for i, url := range filteredUrlList {
			assert.Equal(t, urlList[i], url)
		}
	})

	t.Run("download: mock this ", func(t *testing.T) {
		mockDownloader := new(MockDownloader)
		mockDownloader.On("Download", mock.Anything, mock.Anything).Times(4)
		DownloadImages("https://www.rightmove.co.uk/", mockDownloader)

		mockDownloader.AssertExpectations(t)
	})
}
