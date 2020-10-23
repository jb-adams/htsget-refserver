// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module urldao_test tests module urldao
package htsdao

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsticket"
	"github.com/stretchr/testify/assert"
)

// urlDaoGetContentLengthTC test cases for GetContentLength
var urlDaoGetContentLengthTC = []struct {
	id, url string
	exp     int64
}{
	{
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		"https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/A1-B000168-3_57_F-1-1_R2.mus.Aligned.out.sorted.bam",
		41158,
	},
	{
		"tabulamuris.10X_P4_0",
		"https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
		17093057664,
	},
}

// urlDaoGetByteRangeUrlsTC test cases for GetByteRangeUrls
var urlDaoGetByteRangeUrlsTC = []struct {
	id, url string
	exp     []*htsticket.URL
}{
	{
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		"https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/A1-B000168-3_57_F-1-1_R2.mus.Aligned.out.sorted.bam",
		[]*htsticket.URL{
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/A1-B000168-3_57_F-1-1_R2.mus.Aligned.out.sorted.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=0-41157",
				},
			},
		},
	},
	{
		"tabulamuris.10X_P4_0",
		"https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
		[]*htsticket.URL{
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=0-499999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=500000000-999999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=1000000000-1499999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=1500000000-1999999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=2000000000-2499999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=2500000000-2999999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=3000000000-3499999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=3500000000-3999999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=4000000000-4499999999",
				},
			},
			&htsticket.URL{
				URL: "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/10X_P4_0_possorted_genome.bam",
				Headers: &htsticket.Headers{
					Range: "bytes=4500000000-4999999999",
				},
			},
		},
	},
}

// urlDaoStringTC test cases for String
var urlDaoStringTC = []struct {
	id, url, exp string
}{
	{
		"object0001",
		"https://datasource.com/objects/object0001.bam",
		"URLDao id=object0001, url=https://datasource.com/objects/object0001.bam",
	},
	{
		"tabulamuris.A1",
		"https://tabulamuris.com/data/tabulamuris.A1.bam",
		"URLDao id=tabulamuris.A1, url=https://tabulamuris.com/data/tabulamuris.A1.bam",
	},
}

// minInt convenience function get the minimum of 2 numbers
func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestURLDaoGetContentLength tests GetContentLength function
func TestURLDaoGetContentLength(t *testing.T) {
	for _, tc := range urlDaoGetContentLengthTC {
		dao := NewURLDao(tc.id, tc.url)
		assert.Equal(t, tc.exp, dao.GetContentLength())
	}
}

// TestURLDaoGetByteRangeUrls tests GetByteRangeUrls function
func TestURLDaoGetByteRangeUrls(t *testing.T) {
	for _, tc := range urlDaoGetByteRangeUrlsTC {
		dao := NewURLDao(tc.id, tc.url)
		byteRangeURLs := dao.GetByteRangeUrls()
		// only compare the first ten URLs in returned list
		byteRangeURLsMax10 := byteRangeURLs[:minInt(len(byteRangeURLs), 10)]
		assert.Equal(t, tc.exp, byteRangeURLsMax10)
	}
}

// TestURLDaoString tests String function
func TestURLDaoString(t *testing.T) {
	for _, tc := range urlDaoStringTC {
		dao := NewURLDao(tc.id, tc.url)
		assert.Equal(t, tc.exp, dao.String())
	}
}
