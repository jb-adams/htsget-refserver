// Package preprocessing involves preprocessing of state-related objects of
// an incoming request
//
// Module regionpreprocessor contains tools for sorting, merging an incoming
// set of requested regions to match the order in the file
package preprocessing

import (
	"io/ioutil"
	"math"
	"sort"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htscli"
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

func PreprocessRegions(request *htsrequest.HtsgetRequest) [][]*htsrequest.Region {
	sortRegions(request)
	mergeRegions(request)
	return bridgeRegions(request)
}

func sortRegions(request *htsrequest.HtsgetRequest) {

	orderedSequenceNamesFunctions := map[htsconstants.APIEndpoint]func(*htsrequest.HtsgetRequest) ([]string, error){
		htsconstants.APIEndpointReadsTicket:    htsrequest.GetReferenceNamesInReadsObject,
		htsconstants.APIEndpointVariantsTicket: htsrequest.GetReferenceNamesInVariantsObject,
	}
	orderedSequenceNameFunction := orderedSequenceNamesFunctions[request.GetEndpoint()]
	orderedSequenceNames, err := orderedSequenceNameFunction(request)
	if err != nil {
		return
	}
	rankedSequenceNames := make(map[string]int)
	for i, seqName := range orderedSequenceNames {
		rankedSequenceNames[seqName] = i
	}

	regions := request.GetRegions()
	sort.Slice(regions, func(i int, j int) bool {
		iRefName := regions[i].GetReferenceName()
		iRefNameRank := rankedSequenceNames[iRefName]

		jRefName := regions[j].GetReferenceName()
		jRefNameRank := rankedSequenceNames[jRefName]

		// first, compare ranked reference sequence names
		if iRefNameRank < jRefNameRank {
			return true
		} else if iRefNameRank > jRefNameRank {
			return false
		} else if iRefNameRank == jRefNameRank {
			// if reference sequence names are equal, compare start coordinate
			iStart := regions[i].GetStart()
			jStart := regions[j].GetStart()

			if iStart < jStart {
				return true
			} else if iStart > jStart {
				return false
			} else if iStart == jStart {
				// if start positions are equal, compare end coordinate
				iEnd := regions[i].GetEnd()
				jEnd := regions[j].GetEnd()

				if iEnd < jEnd {
					return true
				}
				return false
			}
		}
		return true
	})
}

func mergeRegions(request *htsrequest.HtsgetRequest) {
	regions := request.GetRegions()
	i := 0

	for i < len(regions)-1 {
		regionI := regions[i]
		regionJ := regions[i+1]

		// cannot merge regions if they are on different contigs
		if regionI.GetReferenceName() != regionJ.GetReferenceName() {
			i++
		} else {
			// as regions are sorted, if the start of j is on or after the start
			// of i, AND on or before the end of i, they are contained one
			// within the other
			if regionJ.GetStart() >= regionI.GetStart() && regionJ.GetStart() <= regionI.GetEnd() {
				// successfully found two regions to merge
				newStart := htsutils.Min(regionI.GetStart(), regionJ.GetStart())
				newEnd := htsutils.Max(regionI.GetEnd(), regionJ.GetEnd())
				newRegion := htsrequest.NewRegion()
				newRegion.SetReferenceName(regionI.GetReferenceName())
				newRegion.SetStart(newStart)
				newRegion.SetEnd(newEnd)

				// delete the two former regions
				regions = append(regions[:i+1], regions[i+2:]...) // delete j
				regions = append(regions[:i], regions[i+1:]...)   // delete i

				// replace with the new merged region
				a := append([]*htsrequest.Region{newRegion}, regions[i:]...)
				regions = append(regions[:i], a...)
			} else {
				i++
			}
		}
	}
	request.SetRegions(regions)
}

// bridgeLengthDeterminationReads determines an estimate of read length within
// the file by inspecting the first few reads. This determines a reasonable
// estimate for bridging separated regions into a single command
func bridgeLengthDeterminationReads(request *htsrequest.HtsgetRequest) int {

	// get the file URL
	fileURL, err := htsconfig.GetObjectPath(request.GetEndpoint(), request.GetID())
	if err != nil {
		return 0
	}

	// pipe a basic samtools view command into a `head` command to get the first
	// 100 reads
	samtoolsView := htscli.SamtoolsView().AddFilePath(fileURL)
	head := htscli.NewCommand()
	head.SetBaseCommand("head")
	head.AddArg("-100")

	// run the command chain
	commandChain := htscli.NewCommandChain()
	commandChain.AddCommand(samtoolsView.GetCommand())
	commandChain.AddCommand(head)
	commandChain.SetupCommandChain()
	pipe := commandChain.ExecuteCommandChain()
	bytes, err := ioutil.ReadAll(pipe)
	if err != nil {
		return 0
	}
	stdout := string(bytes)
	reads := strings.Split(stdout, "\n")
	minReadLength := math.MaxInt32

	for _, read := range reads {
		seq := strings.Split(read, "\t")[9]
		seqlen := len(seq)
		if seqlen < minReadLength {
			minReadLength = seqlen
		}
	}
	return minReadLength
}

func bridgeLengthDeterminationVariants(request *htsrequest.HtsgetRequest) int {
	return 0
}

func bridgeRegions(request *htsrequest.HtsgetRequest) [][]*htsrequest.Region {

	bridgeLengthFunctions := map[htsconstants.APIEndpoint]func(*htsrequest.HtsgetRequest) int{
		htsconstants.APIEndpointReadsTicket:    bridgeLengthDeterminationReads,
		htsconstants.APIEndpointVariantsTicket: bridgeLengthDeterminationVariants,
	}

	i := 0
	regions := request.GetRegions()
	bridgedRegions := [][]*htsrequest.Region{}

	bridgeLengthFunction := bridgeLengthFunctions[request.GetEndpoint()]
	bridgeLength := bridgeLengthFunction(request)

	bridgedRegion := []*htsrequest.Region{regions[0]}
	for i < len(regions)-1 {
		regionI := regions[i]
		regionJ := regions[i+1]

		newBridgeRegion := true

		// if reference names are different, they are not in the same bridge
		if regionI.GetReferenceName() == regionJ.GetReferenceName() {
			distance := regionJ.GetStart() - regionI.GetEnd()
			if distance < bridgeLength {
				newBridgeRegion = false
			}
		}

		if newBridgeRegion {
			bridgedRegions = append(bridgedRegions, bridgedRegion)
			bridgedRegion = []*htsrequest.Region{regionJ}
		} else {
			bridgedRegion = append(bridgedRegion, regionJ)
		}
		i++
	}
	bridgedRegions = append(bridgedRegions, bridgedRegion)

	return bridgedRegions
}
