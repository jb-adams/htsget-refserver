// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module urldao defines access methods for data accessible by URL
package htsdao

import (
	"math"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

// URLDao data access to reads/variants over a network (from an HTTP URL)
type URLDao struct {
	id  string
	url string
}

// NewURLDao instantiates a new URLDao
func NewURLDao(id string, url string) *URLDao {
	dao := new(URLDao)
	dao.id = id
	dao.url = url
	return dao
}

// GetContentLength gets the size in bytes of a file over the internet
func (dao *URLDao) GetContentLength() int64 {
	res, _ := http.Head(dao.url)
	return res.ContentLength
}

// GetByteRangeUrls constructs URLs for an htsget ticket, which provide simple,
// non-overlapping byte range indices to download the target file in blocks
func (dao *URLDao) GetByteRangeUrls() []*htsticket.URL {

	numBytes := dao.GetContentLength()
	blockSize := htsconstants.SingleBlockByteSize
	var start, end int64 = 0, 0
	numBlocks := int(math.Ceil(float64(numBytes) / float64(blockSize)))
	urls := []*htsticket.URL{}
	for i := 1; i <= numBlocks; i++ {
		end = start + blockSize - 1
		if end >= numBytes {
			end = numBytes - 1
		}
		headers := htsticket.NewHeaders()
		headers.SetRangeHeader(start, end)
		url := htsticket.NewURL()
		url.SetURL(dao.url)
		url.SetHeaders(headers)
		start = end + 1
		urls = append(urls, url)
	}
	return urls
}

// String get the URLDao object represented as a string
func (dao *URLDao) String() string {
	return "URLDao id=" + dao.id + ", url=" + dao.url
}
