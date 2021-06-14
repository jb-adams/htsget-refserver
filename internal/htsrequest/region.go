// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module region contains genomic intervals
package htsrequest

import (
	"errors"
	"strconv"
	"strings"
)

// Region defines a simple genomic interval: contig name, start, and end position
type Region struct {
	ReferenceName string `json:"referenceName"`
	Start         *int   `json:"start"`
	End           *int   `json:"end"`
}

// NewRegion instantiates a Region instance
func NewRegion() *Region {
	return new(Region)
}

func RegionFromString(regionString string) (*Region, error) {
	region := new(Region)
	colonSplit := strings.Split(regionString, ":")
	region.SetReferenceName(colonSplit[0])
	if len(colonSplit) > 1 {
		dashSplit := strings.Split(colonSplit[1], "-")
		startI, errStart := strconv.Atoi(dashSplit[0])
		if errStart != nil {
			return nil, errors.New("Error parsing region from string")
		}
		region.SetStart(startI)

		if len(dashSplit) > 1 {
			endI, errEnd := strconv.Atoi(dashSplit[1])
			if errEnd != nil {
				return nil, errors.New("Error parsing region from string")
			}
			region.SetEnd(endI)
		}
	}

	return region, nil
}

/* SETTERS AND GETTERS */

// SetReferenceName sets a region's reference name
func (region *Region) SetReferenceName(referenceName string) {
	region.ReferenceName = referenceName
}

// GetReferenceName retrieves a region's reference name
func (region *Region) GetReferenceName() string {
	return region.ReferenceName
}

// SetStart sets a region's start position
func (region *Region) SetStart(start int) {
	region.Start = &start
}

// GetStart retrieves a region's start position
func (region *Region) GetStart() int {
	return *region.Start
}

// SetEnd sets a region's end position
func (region *Region) SetEnd(end int) {
	region.End = &end
}

// GetEnd retrieves a region's end position
func (region *Region) GetEnd() int {
	return *region.End
}

// StartString retrieves the start position as a string
func (region *Region) StartString() string {
	return strconv.Itoa(region.GetStart())
}

// EndString retrieves the end position as a string
func (region *Region) EndString() string {
	return strconv.Itoa(region.GetEnd())
}

/* API METHODS */

// ReferenceNameRequested validates whether a real reference name has been requested
func (region *Region) ReferenceNameRequested() bool {
	return !(region.GetReferenceName() == "")
}

// StartRequested validates whether a real start position has been requested
func (region *Region) StartRequested() bool {
	if region.Start == nil {
		return false
	}
	return !(region.GetStart() == -1)
}

// EndRequested validates whether a real end position has been requested
func (region *Region) EndRequested() bool {
	if region.End == nil {
		return false
	}
	return !(region.GetEnd() == -1)
}

// String gets a representation of a genomic region
func (region *Region) String() string {
	if !region.StartRequested() && !region.EndRequested() {
		return region.ReferenceName
	}
	if region.StartRequested() && !region.EndRequested() {
		return region.ReferenceName + ":" + region.StartString()
	}
	if !region.StartRequested() && region.EndRequested() {
		return region.ReferenceName + ":" + "0-" + region.EndString()
	}
	return region.ReferenceName + ":" + region.StartString() + "-" + region.EndString()
}

// ExportSamtools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (region *Region) ExportSamtools() string {
	return region.String()
}

// ExportBcftools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (region *Region) ExportBcftools() string {
	if !region.StartRequested() && !region.EndRequested() {
		return region.ReferenceName
	}
	if region.StartRequested() && !region.EndRequested() {
		return region.ReferenceName + ":" + region.StartString() + "-"
	}
	if !region.StartRequested() && region.EndRequested() {
		return region.ReferenceName + ":" + "0-" + region.EndString()
	}
	return region.ReferenceName + ":" + region.StartString() + "-" + region.EndString()
}
