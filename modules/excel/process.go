package excel

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
)

func (cxlsm *CXLSManager) FumigateFile(filename string, visuals bool) {
	urls, found, err := cxlsm.doCheckFile(filename, visuals)
	if err != nil {
		fmt.Printf("error fumigating file %s : %#v\n", filename, urls)
	}
	if found {
		cxlsm.Honeys[filename] = urls
	}

}
func (cxlsm *CXLSManager) FumigateDir(dirname string, visuals bool) {

	var wg sync.WaitGroup

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, ".xlsx") {

				wg.Add(1)
				// fmt.Printf("%s %s\n", de.ModeType(), osPathname)
				go func() {
					urls, found, _ := cxlsm.doCheckFile(osPathname, visuals)
					/*if err!=nil {
						fmt.Printf("%s\n", err)
					}*/
					if found {
						cxlsm.Mutex.Lock()
						cxlsm.Honeys[osPathname] = urls
						cxlsm.Mutex.Unlock()
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

func (cxlsm *CXLSManager) doCheckFile(path string, visuals bool) ([]string, bool, error) {

	var urls = make([]string, 0)
	var bar *progressbar.ProgressBar

	zr, err := zip.OpenReader(path)
	if err != nil {
		return urls, false, err
	}
	defer zr.Close()

	if visuals {
		fmt.Printf("%s\n", path)
	}

	if visuals {
		bar = progressbar.Default(int64(len(zr.File)))
	}

	for _, f := range zr.File {
		fr, err := f.Open()
		if err != nil {
			return urls, false, err
		}
		defer fr.Close()

		if visuals {
			_ = bar.Add(1)
		}

		// Should not be any dirs but skip if present
		if !f.FileInfo().IsDir() {
			filesize := f.FileInfo().Size()

			buffer := make([]byte, filesize)
			_, err := fr.Read(buffer)
			if err != nil {
				// Done reading
				if err != io.EOF {
					log.Println("fr.Read:", err)
					return urls, false, err
				}
			}
			rUrl := regexp.MustCompile(`(?i)Target="(?P<link>http.+?)" `)
			if ok, urlst := cxlsm.processZFile(buffer, rUrl); ok {

				for _, u := range urlst {
					// TODO: cleanup
					// fmt.Printf("\tUrl: %s\n", u)
					urls = append(urls, u)
				}
			}
		}
	}

	if len(urls) > 0 {
		return urls, true, nil
	}
	return urls, false, nil
}

func (cxlsm *CXLSManager) processZFile(content []byte, pattern *regexp.Regexp) (bool, []string) {

	var result []string
	if pattern.Match(content) {
		matches := pattern.FindAllSubmatch(content, -1)
		for _, v := range matches {
			result = append(result, string(v[1]))
		}
		return true, result
	} else {
		return false, result
	}
}
