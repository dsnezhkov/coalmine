package detector

import (
	"fmt"
	"log"
	"regexp"
)

type CanaryOrgDetector struct {
	Tokens []*regexp.Regexp
}


var (
	// http://j2d9n4auf7b5aeaph3jhlbtp3.canarytokens.net/FADYIMGCFISKUHKSKTQGUOFQTFYHOQITFK
	canaryOrgRules = []string{
		".*canarytokens.net",
		"http?.*\\/\\/[a-z0-9A-Z]{25}.",
		"http?.*\\/.*\\/[A-Z]{34}",
	}
	cod *CanaryOrgDetector
)
func CanaryOrgDetectorFactory() *CanaryOrgDetector {

	if cod == nil {
		cod = new(CanaryOrgDetector)
		cod.rules2Regex()
	}
	return cod
}

func (cod *CanaryOrgDetector) rules2Regex() {
	for _,r :=range canaryOrgRules {
		e, err := regexp.Compile(r)
		if err!=nil {
			log.Printf("Unable to compile %s\n", err)
			continue
		}
		cod.Tokens = append(cod.Tokens, e)
	}
}


func (cod *CanaryOrgDetector) LocateHoneys(honeys map[string][]string, showCandidates bool) {
	var located bool
	for k, v := range honeys {
		located = false
		fmt.Printf("%s:\n", k)
		for i,u :=range v {
			for _,r :=range cod.Tokens {
				if r.MatchString(u) {
					located = true
					break
				}
			}
			if located{
				fmt.Printf("-->\t%d: %s\n", i+1, u)
				break
			}
			if showCandidates {
				fmt.Printf("\t%d: %s\n", i+1, u)
			}
		}
	}
}
