package cli

import (
	"fmt"
	"os"
	"regexp"

	"coalmine/modules"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	jitter int
	sequential bool
	maxPasses int64
	location  string
	showCandidate bool
	file2mod map[string]modules.Processor
	file2reg map[string]*regexp.Regexp
)

var RootCmd = &cobra.Command{
	Use:   "coalmine",
	Short: "Canaries in the Coalmine",
	Long: `Seek out canaries and honeytokens in commonly used file formats`,
}


func init() {
	// Collection section
	RootCmd.AddCommand(pdfCmd)
	RootCmd.AddCommand(xlsCmd)
	RootCmd.AddCommand(docCmd)
	RootCmd.AddCommand(allCmd)
	RootCmd.AddCommand(CompletionCmd)

	// Root section
	RootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false,
		"Verbose output",
	)
	RootCmd.PersistentFlags().IntVarP(
		&jitter, "jitter", "j", 0,
		"Time variance (sec) on sequential file access",
	)
	RootCmd.PersistentFlags().BoolVarP(
		&sequential, "sequential", "s", false,
		"Sequential file access in directory scan (slower)",
	)

}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
