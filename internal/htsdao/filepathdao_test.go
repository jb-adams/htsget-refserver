// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module filepathdao_test tests module filepathdao
package htsdao

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsticket"
	"github.com/stretchr/testify/assert"
)

// filePathDaoGetContentLengthTC test cases for GetContentLength
var filePathDaoGetContentLengthTC = []struct {
	id, filePath string
	exp          int64
}{
	{
		"gatk.NA12878",
		"../../data/gcp/gatk-test-data/wgs_bam/NA12878.bam",
		15236350,
	},
	{
		"gatk.NA12878_20k_b37",
		"../../data/gcp/gatk-test-data/wgs_bam/NA12878_20k_b37.bam",
		8820570,
	},
}

// filePathDaoGetByteRangeUrlsTC test cases for GetByteRangeUrls
var filePathDaoGetByteRangeUrlsTC = []struct {
	id, filePath string
	exp          []*htsticket.URL
}{
	{
		"gatk.NA12878",
		"../../data/gcp/gatk-test-data/wgs_bam/NA12878.bam",
		[]*htsticket.URL{
			&htsticket.URL{
				URL: "http://localhost:3000/file-bytes",
				Headers: &htsticket.Headers{
					FilePath: "../../data/gcp/gatk-test-data/wgs_bam/NA12878.bam",
					Range:    "bytes=0-15236349",
				},
			},
		},
	},
	{
		"gatk.NA12878_20k_b37",
		"../../data/gcp/gatk-test-data/wgs_bam/NA12878_20k_b37.bam",
		[]*htsticket.URL{
			&htsticket.URL{
				URL: "http://localhost:3000/file-bytes",
				Headers: &htsticket.Headers{
					FilePath: "../../data/gcp/gatk-test-data/wgs_bam/NA12878_20k_b37.bam",
					Range:    "bytes=0-8820569",
				},
			},
		},
	},
}

// filePathDaoStringTC test cases for String
var filePathDaoStringTC = []struct {
	id, filePath, exp string
}{
	{
		"gatk.NA12878",
		"../../data/gcp/gatk-test-data/wgs_bam/NA12878.bam",
		"FilePathDao id=gatk.NA12878, filePath=../../data/gcp/gatk-test-data/wgs_bam/NA12878.bam",
	},
}

// TestFilePathDaoGetContentLength tests GetContentLength function
func TestFilePathDaoGetContentLength(t *testing.T) {
	for _, tc := range filePathDaoGetContentLengthTC {
		dao := NewFilePathDao(tc.id, tc.filePath)
		assert.Equal(t, tc.exp, dao.GetContentLength())
	}
}

// TestFilePathDaoGetByteRangeUrls tests GetByteRangeUrls function
func TestFilePathDaoGetByteRangeUrls(t *testing.T) {
	for _, tc := range filePathDaoGetByteRangeUrlsTC {
		dao := NewFilePathDao(tc.id, tc.filePath)
		assert.Equal(t, tc.exp, dao.GetByteRangeUrls())
	}
}

// TestFilePathDaoString tests String function
func TestFilePathDaoString(t *testing.T) {
	for _, tc := range filePathDaoStringTC {
		dao := NewFilePathDao(tc.id, tc.filePath)
		assert.Equal(t, tc.exp, dao.String())
	}
}
