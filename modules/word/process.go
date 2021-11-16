package word

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"path"
	"regexp"

	"github.com/schollz/progressbar/v3"
)

func (cdocm *CDOCManager) DemineFile(filename string, visuals bool) (bool, error) {
	urls, found, err := cdocm.doCheckFile(filename, visuals)
	if found {
		cdocm.Mutex.Lock()
		cdocm.Honeys[filename] = urls
		cdocm.Mutex.Unlock()
	}
	return found, err
}

func (cdocm *CDOCManager) doCheckFile(filepath string, visuals bool) ([]string, bool, error) {

	var urls = make([]string, 0)
	var bar *progressbar.ProgressBar

	zr, err := zip.OpenReader(filepath)
	if err != nil {
		return urls, false, err
	}
	defer zr.Close()

	if visuals {
		fmt.Printf("%s\n", filepath)
	}

	if visuals {
		//bar = progressbar.Default(int64(len(zr.File)), path.Base(filepath))
		bar = progressbar.Default(-1, path.Base(filepath))
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
