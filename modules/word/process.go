package word

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

func (cdocm *CDOCManager) FumigateFile(filename string, visuals bool) {
	urls, found, err := cdocm.doCheckFile(filename, visuals)
	if err != nil {
		fmt.Printf("error fumigating file %s : %#v\n", filename, urls)
	}
	if found {
		cdocm.Honeys[filename] = urls
	}

}
func (cdocm *CDOCManager) FumigateDir(dirname string, visuals bool) {

	var wg sync.WaitGroup

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.Contains(osPathname, ".docx") ||
				strings.Contains(osPathname, ".doc") ||
				strings.Contains(osPathname, ".docm") {

				wg.Add(1)
				// fmt.Printf("%s %s\n", de.ModeType(), osPathname)
				go func() {
					urls, found, _ := cdocm.doCheckFile(osPathname, visuals)
					/*if err!=nil {
						fmt.Printf("%s\n", err)
					}*/
					if found {
						cdocm.Mutex.Lock()
						cdocm.Honeys[osPathname] = urls
						cdocm.Mutex.Unlock()
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

func (cdocm *CDOCManager) doCheckFile(path string, visuals bool) ([]string, bool, error) {

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
			// seeking Target= constructs with external attribute
			rUrl := regexp.MustCompile(`(?i)Target="(?P<link>http.+?)" TargetMode="External"`)
			if ok, urlst := cdocm.processZFile(buffer, rUrl); ok {

				for _, u := range urlst {
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

func (cdocm *CDOCManager) processZFile(content []byte, pattern *regexp.Regexp) (bool, []string) {

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
