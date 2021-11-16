package pdf

import (
	"fmt"
	"path"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/schollz/progressbar/v3"
)

func (cpdfm *CPDFManager) DemineFile(filename string, visuals bool) (bool, error) {
	urls, found, err := cpdfm.doCheckFile(filename, visuals)
	if found {
		cpdfm.Mutex.Lock()
		cpdfm.Honeys[filename] = urls
		cpdfm.Mutex.Unlock()
	}
	return found, err
}


func (cpdfm *CPDFManager) doCheckFile(filepath string, visuals bool) ([]string, bool, error) {

	var urls = make([]string,0)
	var bar *progressbar.ProgressBar


	file, err := pdfcpu.ReadFile(filepath, nil)
	if err != nil {
		return urls, false, err
	}

	if visuals {
		fmt.Printf("%s\n", filepath)
	}
	if visuals {
		bar = progressbar.Default(-1, path.Base(filepath))
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

