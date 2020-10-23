// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module daofactory_test tests module daofactory
package htsdao

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// daoFactoryGetDaoTC test cases for GetDao
var daoFactoryGetDaoTC = []struct {
	endpoint htsconstants.APIEndpoint
	id       string
	expError bool
	expClass string
}{
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		false,
		"*htsdao.URLDao",
	},
	{
		htsconstants.APIEndpointVariantsTicket,
		"HG002_GIAB",
		false,
		"*htsdao.FilePathDao",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"HG002_GIAB",
		true,
		"",
	},
}

// TestDaoFactoryGetDao tests GetDao function
func TestDaoFactoryGetDao(t *testing.T) {
	for _, tc := range daoFactoryGetDaoTC {
		req := htsrequest.NewHtsgetRequest()
		req.SetEndpoint(tc.endpoint)
		req.SetID(tc.id)
		dao, err := GetDao(req)
		daoType := fmt.Sprintf("%T", dao)

		if tc.expError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expClass, daoType)
		}
	}
}
