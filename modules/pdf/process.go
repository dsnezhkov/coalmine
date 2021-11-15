package coalpdf

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/schollz/progressbar/v3"
)

func (cpdfm *CPDFManager) FumigateFile(filename string, visuals bool) {
	urls, found, _ := cpdfm.doCheckFile(filename, visuals)
	if found {
		//fmt.Printf("Appending %#v\n", urls)
		cpdfm.Honeys[filename] = urls
	}

}
func (cpdfm *CPDFManager) FumigateDir(dirname string, visuals bool) {

	var wg sync.WaitGroup

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, ".pdf") {

				wg.Add(1)
				// fmt.Printf("%s %s\n", de.ModeType(), osPathname)
				go func() {
					urls, found, _ := cpdfm.doCheckFile(osPathname, visuals)
					/*if err!=nil {
						fmt.Printf("%s\n", err)
					}*/
					if found {
						cpdfm.Mutex.Lock()
						cpdfm.Honeys[osPathname] = urls
						cpdfm.Mutex.Unlock()
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

func (cpdfm *CPDFManager) doCheckFile(path string, visuals bool) ([]string, bool, error) {

	var urls = make([]string,0)
	var bar *progressbar.ProgressBar

	file, err := pdfcpu.ReadFile(path, nil)
	if err != nil {
		return urls, false, err
	}

	if visuals {
		fmt.Printf("%s\n", path)
	}
	if visuals {
		bar = progressbar.Default(cpdfm.MaxObjects)
	}

	// We could not find a reliable way to get to the right object - brute force
	for i := 1; i <= int(cpdfm.MaxObjects); i++ {

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

