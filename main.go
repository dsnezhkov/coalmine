package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/schollz/progressbar/v3"
)

const MaxObjects = 50
var honeys map[string][]string

func main() {

	honeys = make(map[string][]string)

	dirname := flag.String("d", "", "Directory to fumigate")
	filename := flag.String("f", "", "File to fumigate")
	visuals := flag.Bool("v", true, "Visual progress")

	flag.Parse()
	if *dirname == "" && *filename == "" {
		usage()
		return
	}

	if *dirname != "" {
		fumigateDir(*dirname, *visuals)
	}
	if *filename != "" {
		fumigateFile(*filename, *visuals)
	}

	printResults()
}

func fumigateFile(filename string, visuals bool) {
	urls, found, _ := doCheckFile(filename, visuals)
	if found {
		//fmt.Printf("Appending %#v\n", urls)
		honeys[filename] = urls
	}

}
func fumigateDir(dirname string, visuals bool) {

	var wg sync.WaitGroup

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, ".pdf") {

				wg.Add(1)
				// fmt.Printf("%s %s\n", de.ModeType(), osPathname)
				go func() {
					urls, found, _ := doCheckFile(osPathname, visuals)
					/*if err!=nil {
						fmt.Printf("%s\n", err)
					}*/
					if found {
						honeys[osPathname] = urls
					}
					defer wg.Done()
				}()
				return nil
			}
			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		fmt.Printf("Walk err: %s\n", err)
	}

	wg.Wait()
}

func doCheckFile(path string, visuals bool) ([]string, bool, error) {

	var urls = make([]string,0)
	var bar *progressbar.ProgressBar

	if visuals {
		fmt.Printf("%s\n", path)
	}
	file, err := pdfcpu.ReadFile(path, nil)
	if err != nil {
		return urls, false, err
	}

	// We could not find a reliable way to get to the right object - brute force
	if visuals {
		bar = progressbar.Default(MaxObjects)
	}
	for i := 1; i <= MaxObjects; i++ {

		if visuals {
			_ = bar.Add(1)
		}
		time.Sleep(10 * time.Millisecond)

		o, found := file.Find(i)
		if !found || o == nil {
			continue
		}
		o1, ok := o.Object.(pdfcpu.Dict)
		if !ok {
			continue
		}
		if uri, ok := o1["URI"]; ok {
			//fmt.Printf("Appending %s..", uri.String())
			urls = append(urls,uri.String())
		} else {
			continue
		}
	}

	if len(urls) > 0 {
		return urls, true, nil
	}
	return urls, false, nil
}

func printResults(){
	fmt.Println("\n\n=== Seek Canaries in URLs ===")
	for k, v := range honeys {
		fmt.Printf("%s:\n", k)
		for i,u :=range v {
			fmt.Printf("\t%d: %s\n", i+1, u)
		}
	}
}
func usage() {
	fmt.Printf("Usage: coalmine -d=<dir>\n")
}
