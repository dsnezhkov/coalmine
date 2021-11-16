package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"sync"
	"time"

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
	file2mod map[string]modules.Processor, file2reg map[string]*regexp.Regexp,
	options map[string]interface{}) {

	var wg sync.WaitGroup

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			for k, v := range file2reg {
				if v.MatchString(osPathname) {

					if options["sequential"].(bool) {

						// TODO: refactor
						jitter  := options["jitter"].(int)
						if jitter != 0 {
							randPause := rand.Intn(jitter)
							time.Sleep(time.Duration(randPause) * time.Second)
						}

						flail(k, file2mod[k], osPathname, visuals)
					} else {
						wg.Add(1)

						// TODO: refactor
						jitter  := options["jitter"].(int)
						if jitter != 0 {
							randPause := rand.Intn(jitter)
							time.Sleep(time.Duration(randPause) * time.Second)
						}
						go func(format string) {
							flail(format, file2mod[k], osPathname, visuals)
							defer wg.Done()
						}(k)
					}

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

func flail(format string, processor modules.Processor, filepath string, visuals bool) {
	_, err := processor.DemineFile(filepath, visuals)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}
