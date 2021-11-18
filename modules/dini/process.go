package dini

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/ini.v1"
)

func (cdinim *CDINIManager) DemineFile(filename string, visuals bool) (bool, error) {
	urls, found, err := cdinim.doCheckFile(filename, visuals)
	if found {
		cdinim.Mutex.Lock()
		cdinim.Honeys[filename] = urls
		cdinim.Mutex.Unlock()
	}
	return found, err
}

func (cdinim *CDINIManager) doCheckFile(filepath string, visuals bool) ([]string, bool, error) {

	var urls = make([]string, 0)
	var bar *progressbar.ProgressBar


	iniOpts := ini.LoadOptions{
		Loose:                       true,
		Insensitive:                 true,
		InsensitiveSections:         true,
		InsensitiveKeys:             true,
		IgnoreContinuation:          false,
		IgnoreInlineComment:         true,
		SkipUnrecognizableLines:     true,
		ShortCircuit:                false,
		AllowBooleanKeys:            true,
		AllowNestedValues:           true,
		AllowPythonMultilineValues:  false,
		SpaceBeforeInlineComment:    true,
		UnparseableSections:         nil,
		KeyValueDelimiters:          "=",
		ChildSectionDelimiter:       "",
	}
	cfg, err := ini.LoadSources(iniOpts, filepath)
	if err != nil {
		return urls, false, err
	}

	if visuals {
		fmt.Printf("%s\n", filepath)
	}
	if visuals {
		bar = progressbar.Default(-1, path.Base(filepath))
	}

	sections := cfg.Sections()
	if len(sections) == 0 {
		return urls, false, err
	}

	// seeking constructs with urls
	// We assume any UNC and external pointers are valuable to alert on
	rUrl := regexp.MustCompile(`.+=\s?(\\\\|\\|\/|http.?:\/\/|file:\/\/).+`)
	if ok, urls := cdinim.processFile(sections, rUrl, bar); ok {
		return urls, true, nil
	}

	return urls, false, nil
}

func (cdinim *CDINIManager) processFile(cfgSections []*ini.Section, pattern *regexp.Regexp, bar *progressbar.ProgressBar) (bool, []string) {

	var result []string

	for _, k := range cfgSections{
		_ = bar.Add(1)
		for _, v :=range k.Keys(){

			// TODO: I really don't know how to deal with BOMs and multibytes in Regexes
			// Clean them up
			vNameCanonical := strings.ReplaceAll(v.Name(), "\x00", "")
			vValCanonical := strings.ReplaceAll(v.Value(), "\x00", "")

			// Construct row out of keys and values and match
			row := strings.Join([]string{vNameCanonical, vValCanonical},"=")
			if pattern.Match([]byte(row)) {
				result = append(result, vValCanonical)
			}
		}
	}
	if len(result) > 0 {
		return true, result
	}
	return false, result
}
