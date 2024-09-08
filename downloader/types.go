package downloader

type IDownloader interface {
	Download(url string, fileName string)
}
