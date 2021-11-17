package main

import (
	"coalmine/cmd/coalmine/cli"
)




func main() {
	cli.Execute()




	/*dirname := flag.String("d", "", "Directory to fumigate")
	filename := flag.String("f", "", "File to fumigate")
	visuals := flag.Bool("v", true, "Visual progress")
	 */

	// flag.Parse()
	// if *dirname == "" && *filename == "" {
	// 	usage()
	// 	return
	// }
	//
	// cpdfm := cpdf.CPDFManagerFactory()
	// if *dirname != "" {
	// 	cpdfm.FumigateDir(*dirname, *visuals)
	// }
	// if *filename != "" {
	// 	cpdfm.FumigateFile(*filename, *visuals)
	// }
	//
	// cpdfm.PrintResults()
}

// func usage() {
// 	fmt.Printf("Usage: coalmine -d=<dir>\n")
// }
