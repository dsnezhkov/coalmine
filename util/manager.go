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

var rtStats map[string]int64
var  rtStatsM sync.RWMutex

func init() {

	rtStats = make(map[string]int64,0)
	rtStats["DirSize"] = 0
	rtStats["TotalFilesProcessed"] = 0
	rtStats["FilesInScope"] = 0
	rtStats["TotalTimeSec"] = 0
}

func PrintStats() {
	fmt.Printf("\n\n=============== Stats ===============\n")
	for k, v :=range rtStats {
		fmt.Printf("%s : %d\n", k, v)
	}
}
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
	start := time.Now()

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			rtStatsM.Lock()
			rtStats["TotalFilesProcessed"] += 1
			rtStatsM.Unlock()

			for k, v := range file2reg {
				if v.MatchString(osPathname) {
					rtStatsM.Lock()
					rtStats["FilesInScope"] += 1
					rtStatsM.Unlock()

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

	rtStats["TotalTimeSec"] = int64(time.Since(start).Seconds())
}

func flail(format string, processor modules.Processor, filepath string, visuals bool) {
	_, err := processor.DemineFile(filepath, visuals)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}
