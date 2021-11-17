package cli

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"coalmine/detector"
	"coalmine/modules"
	"coalmine/modules/word"
	"coalmine/util"
	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "word",
	Short: "Process word",
	Long: "Seek out canaries in Word objects",
	Run: processDOC,
}

func init() {
	// DOC section
	docCmd.Flags().StringVarP(
		&location, "location", "l", "",
		"File or Folder to fumigate",
	)
	docCmd.Flags().BoolVarP(
		&showCandidate, "candidate", "c", true,
		"Show unverified candidates",
	)
	_ = docCmd.MarkFlagRequired("location")
}

func processDOC (cmd *cobra.Command, args []string) {

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
	file2reg["doc"] = regexp.MustCompile(`(?i)^.*\.(doc|docx|docm)$`)

	file2mod = make(map[string]modules.Processor, 0)
	file2mod["doc"] = word.CDOCManagerFactory()

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
