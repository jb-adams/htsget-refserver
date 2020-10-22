package htsserver

import (
	"net/http"
	"os/exec"

	"github.com/ga4gh/htsget-refserver/internal/htscli"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

func getVariantsData(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointVariantsData,
		addRegionFromQueryString,
		getVariantsDataHandler,
	).handleRequest(writer, request)
}

// getVariantsData serves the actual data from AWS back to client
func getVariantsDataHandler(handler *requestHandler) {
	fileURL, err := htsconfig.GetObjectPath(handler.HtsReq.GetEndpoint(), handler.HtsReq.GetID())
	if err != nil {
		return
	}

	commandChain := htscli.NewCommandChain()
	removedHeadBytes := 0
	removedTailBytes := 0

	if handler.HtsReq.IsHeaderBlock() {
		// only get the header for header blocks
		commandChain.AddCommand(bcftoolsViewHeaderOnlyVCF(handler.HtsReq, fileURL))
	} else {
		// body-based requests
		if handler.HtsReq.GetFormat() == "BCF" {
			removedHeadBytes, _ = getBCFHeaderByteSize(handler.HtsReq, fileURL)
		}
		commandChain.AddCommand(bcftoolsViewBodyVCF(handler.HtsReq, fileURL))
	}

	// execute command chain and stream output
	commandWriteStream(commandChain, removedHeadBytes, removedTailBytes, handler.Writer)
}

func bcftoolsViewHeaderOnlyVCF(htsgetReq *htsrequest.HtsgetRequest, fileURL string) *htscli.Command {
	cmd := htscli.BcftoolsView()
	cmd.SetFilePath(fileURL)
	cmd.SetHeaderOnly(true)
	cmd.SetOutputFormat(htsgetReq.GetFormat())
	return cmd.GetCommand()
}

func bcftoolsViewBodyVCF(htsgetReq *htsrequest.HtsgetRequest, fileURL string) *htscli.Command {
	cmd := htscli.BcftoolsView()
	cmd.SetFilePath(fileURL)
	cmd.SetHeaderOnly(false)
	cmd.SetOutputFormat(htsgetReq.GetFormat())
	if !htsgetReq.AllRegionsRequested() {
		cmd.SetRegion(htsgetReq.GetRegions()[0])
	}
	return cmd.GetCommand()
}

func getBCFHeaderByteSize(htsgetReq *htsrequest.HtsgetRequest, fileURL string) (int, error) {

	cmd := exec.Command("bcftools", "view", "--no-version", "-h", "-O", "u", fileURL)
	tmpHeader, err := htsconfig.CreateTempFile(htsgetReq.GetID() + "_header")
	if err != nil {
		return 0, err
	}

	cmd.Stdout = tmpHeader
	cmd.Run()

	fi, err := tmpHeader.Stat()
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	tmpHeader.Close()
	htsconfig.RemoveTempfile(tmpHeader)
	return int(size), nil
}
