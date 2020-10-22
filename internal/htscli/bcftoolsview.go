// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module bcftoolsview defines the job submission for the 'bcftools view'
// command, which streams VCF or BCF to stdout
package htscli

import (
	"fmt"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

// BcftoolsViewCommand represents a single 'bcftools view' command and associated
// arguments
type BcftoolsViewCommand struct {
	filePath   string
	headerOnly bool
	outputVCF  bool
	region     *htsrequest.Region
}

// BcftoolsView instantiates a new BcftoolsView Command
func BcftoolsView() *BcftoolsViewCommand {
	return new(BcftoolsViewCommand)
}

// SetFilePath sets the path to the requested variant file
func (bcftoolsViewCommand *BcftoolsViewCommand) SetFilePath(filePath string) {
	bcftoolsViewCommand.filePath = filePath
}

// SetHeaderOnly sets boolean parameter that, if true, will stream only the header.
// if false, the header is excluded entirely
func (bcftoolsViewCommand *BcftoolsViewCommand) SetHeaderOnly(headerOnly bool) {
	bcftoolsViewCommand.headerOnly = headerOnly
}

// SetOutputVCF sets boolean parameter that, if true, will stream variant data
// as VCF. if false, variant data is streamed as BCF
func (bcftoolsViewCommand *BcftoolsViewCommand) SetOutputVCF(outputVCF bool) {
	bcftoolsViewCommand.outputVCF = outputVCF
}

// SetOutputFormat sets output format based on string keywords "VCF" or "BCF"
func (bcftoolsViewCommand *BcftoolsViewCommand) SetOutputFormat(format string) {
	if format == "VCF" {
		bcftoolsViewCommand.SetOutputVCF(true)
	} else if format == "BCF" {
		bcftoolsViewCommand.SetOutputVCF(false)
	} else {
		bcftoolsViewCommand.SetOutputVCF(true)
	}
}

// SetRegion sets the requested genomic region for variant streaming
func (bcftoolsViewCommand *BcftoolsViewCommand) SetRegion(region *htsrequest.Region) {
	bcftoolsViewCommand.region = region
}

// GetCommand exports the BcftoolsViewCommand as a generic Command
func (bcftoolsViewCommand *BcftoolsViewCommand) GetCommand() *Command {
	// consistent base command and initial args
	command := NewCommand()
	command.SetBaseCommand("bcftools")
	command.AddArg("view")
	command.AddArg(bcftoolsViewCommand.filePath)
	command.AddArg("--no-version")

	// add header flag
	if bcftoolsViewCommand.headerOnly {
		command.AddArg("-h")
	} else {
		if bcftoolsViewCommand.outputVCF {
			command.AddArg("-H")
		}
	}

	// always output as uncompressed VCF
	// TODO make 'format' parameter effective
	command.AddArg("-O")
	if bcftoolsViewCommand.outputVCF {
		command.AddArg("v")
	} else {
		command.AddArg("u")
	}

	// add region interval flag
	if bcftoolsViewCommand.region != nil {
		command.AddArg("-r")
		command.AddArg(bcftoolsViewCommand.region.ExportBcftools())
	}

	fmt.Println(command.baseCommand)
	fmt.Println(command.args)

	return command
}
