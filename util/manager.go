package util

import (
	"fmt"
	"regexp"
	"sync"

	"coalmine/modules"
	"github.com/karrick/godirwalk"
)

func DemineFile(filename string, visuals bool,
	file2mod map[string]modules.Processor, file2reg map[string]*regexp.Regexp) string {

	var format string
	for k, v := range file2reg {
		if v.MatchString(filename) {
			format = k
			_, err := file2mod[k].DemineFile(filename, visuals)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
		}
	}
	return format
}


func DemineDir(dirname string, visuals bool,
	file2mod map[string]modules.Processor, file2reg map[string]*regexp.Regexp) {

	var wg sync.WaitGroup

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			for k, v := range file2reg {
				if v.MatchString(osPathname) {
					wg.Add(1)
					go func(format string) {
						flail(format,file2mod[k],osPathname,visuals)
						defer wg.Done()
					}(k)

					break
				}
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

func flail(format string, processor modules.Processor, filepath string, visuals bool){
	_, err := processor.DemineFile(filepath, visuals)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}