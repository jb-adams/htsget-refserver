// Package htsgetconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults.go contains default runtime properties when not overriden
// by environment properties
package htsgetconfig

// getDefaults gets all default properties
//
// Returns
//	(map[string]string): map of default properties
func getDefaults() map[string]string {
	defaults := map[string]string{
		"port": "3000",
		"host": "http://localhost:3000",
	}
	return defaults
}

// getDefaultReadsSourcesRegistry gets the default source registry for 'reads' endpoint
//
// Returns
//	(*DataSourceRegistry): default reads source registry. points to tabula muris and local test files
func getDefaultReadsSourcesRegistry() *DataSourceRegistry {
	sources := []map[string]string{
		{
			"pattern": "^tabulamuris\\.(?P<accession>10X.*)$",
			"path":    "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam",
		},
		{
			"pattern": "^tabulamuris\\.(?P<accession>.*)$",
			"path":    "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam",
		},
	}

	registry := newDataSourceRegistry()
	for i := 0; i < len(sources); i++ {
		registry.addDataSource(newDataSource(sources[i]["pattern"], sources[i]["path"]))
	}
	return registry
}
