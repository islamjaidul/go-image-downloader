package downloader

type IDownloader interface {
	download(url string, fileName string)
}
