package cli

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"coalmine/detector"
	"coalmine/modules"
	"coalmine/modules/pdf"
	"coalmine/util"
	"github.com/spf13/cobra"
)


var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Process pdf",
	Long: "Seek out canaries in PDF objects",
	Run: processPDF,
}

func init() {
	// PDF section
	pdfCmd.Flags().Int64VarP(
		&maxPasses, "pass", "p", 10,
		"Max passes for URL object extraction",
	)
	pdfCmd.Flags().StringVarP(
		&location, "location", "l", "",
		"File or Folder to fumigate",
	)
	pdfCmd.Flags().BoolVarP(
		&showCandidate, "candidate", "c", true,
		"Show unverified candidates",
	)
	_ = pdfCmd.MarkFlagRequired("location")
}

func processPDF (cmd *cobra.Command, args []string) {

	demineOptions := make (map[string]interface{}, 0)
	vSeq, err := cmd.Flags().GetBool("sequential")
	if err != nil {
		demineOptions["sequential"]	= false
	}
	demineOptions["sequential"]	= vSeq

	vJit, err := cmd.Flags().GetInt("jitter")
	if err != nil {
		demineOptions["jitter"]	= 0
	}
	demineOptions["jitter"]	= vJit


	file2reg = make(map[string]*regexp.Regexp, 0)

	// Compile once
	file2reg["pdf"] = regexp.MustCompile(`(?i)^.*\.pdf$`)

	file2mod = make(map[string]modules.Processor, 0)
	file2mod["pdf"] = pdf.CPDFManagerFactory()

	fileInfo, err := os.Stat(location)
	if err != nil {
		log.Fatalf("Unable to stat location: %s\n", err)
	}

	cod := detector.CanaryOrgDetectorFactory()

	if fileInfo.IsDir() {
		util.DemineDir(location, verbose, file2mod, file2reg, demineOptions)
		for k, v := range file2mod {
			fmt.Printf("=== Format: %s (verified canaries marked with `->`) === \n", k)
			cod.LocateHoneys(v.GetHoneys(), showCandidate)
		}
	}else{
		format := util.DemineFile(location, verbose, file2mod, file2reg)
		cod.LocateHoneys(file2mod[format].GetHoneys(), showCandidate)
	}
}