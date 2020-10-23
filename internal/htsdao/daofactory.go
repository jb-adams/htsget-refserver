// Package htsdao defines Data Access Objects, which
// opens handles to reads/variants data based on different data sources
// (local files, http URLs, etc.)
//
// Module daofactory constructs the correct Data Access Object based on the
// type of object path passed to it
package htsdao

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

// getMatchingDao constructs either a FilePathDao or URLDao based on whether
// the requested id maps to a local file path or a URL
func getMatchingDao(id string, registry *htsconfig.DataSourceRegistry) (DataAccessObject, error) {
	path, err := registry.GetMatchingPath(id)
	if err != nil {
		return nil, err
	}
	if htsutils.IsValidURL(path) {
		return NewURLDao(id, path), nil
	}
	return NewFilePathDao(id, path), nil
}

// GetDao gets the correct Data Access Object based on the ID of an htsget request
func GetDao(req *htsrequest.HtsgetRequest) (DataAccessObject, error) {
	req.SetMatchingDataSource()
	registry := req.GetDataSourceRegistry()
	return getMatchingDao(req.GetID(), registry)
}
