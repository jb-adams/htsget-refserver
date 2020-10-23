// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module filepathdao defines access methods for data accessible by local file path
package htsdao

import (
	"math"
	"os"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

// FilePathDao data access to reads/variants from a local file path
type FilePathDao struct {
	id       string
	filePath string
}

// NewFilePathDao instantiates a new FilePathDao
func NewFilePathDao(id string, filePath string) *FilePathDao {
	dao := new(FilePathDao)
	dao.id = id
	dao.filePath = filePath
	return dao
}

// GetContentLength gets the size in bytes of a local reads/variants file
func (dao *FilePathDao) GetContentLength() int64 {
	fileInfo, _ := os.Stat(dao.filePath)
	return fileInfo.Size()
}

// constructByteRangeURL constructs a single byte range URL for the response ticket
func (dao *FilePathDao) constructByteRangeURL(start int64, end int64) *htsticket.URL {
	host := htsconfig.GetHost()
	path := host + htsconstants.FileByteRangeURLPath
	headers := htsticket.NewHeaders()
	headers.SetRangeHeader(start, end)
	headers.SetFilePathHeader(dao.filePath)
	url := htsticket.NewURL()
	url.SetURL(path)
	url.SetHeaders(headers)
	return url
}

// GetByteRangeUrls constructs URLs for an htsget ticket, which provide simple,
// non-overlapping byte range indices to download the target file in blocks
func (dao *FilePathDao) GetByteRangeUrls() []*htsticket.URL {
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
		url := dao.constructByteRangeURL(start, end)
		start = end + 1
		urls = append(urls, url)
	}
	return urls
}

// String get the FilePathDao object represented as a string
func (dao *FilePathDao) String() string {
	return "FilePathDao id=" + dao.id + ", filePath=" + dao.filePath
}
