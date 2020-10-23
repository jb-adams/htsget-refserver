// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module dao defines the DataAccessObject interface
package htsdao

import "github.com/ga4gh/htsget-refserver/internal/htsticket"

// DataAccessObject interface defining common methods for accessing data
// from a reads/variants object. Implemented differently for different data
// sources
type DataAccessObject interface {
	GetContentLength() int64
	GetByteRangeUrls() []*htsticket.URL
	String() string
}
