package language

import (
	"github.com/astrolink/gutils/general"
	"strconv"
	"strings"
)

var translationKeys map[string]string
var contexts []string

func LoadLang(lang map[string]map[string]string, context string, idiom string) {
	contextIdiom := context + "_" + idiom

	testInArray, _ := general.InArray(contextIdiom, contexts)

	if testInArray == true {
		return
	}

	if len(translationKeys) == 0 {
		translationKeys = make(map[string]string)
	}

	if val, ok := lang[idiom]; ok {
		for key, value := range val {
			translationKeys[key+"_"+idiom] = value
		}
	}

	contexts = append(contexts, contextIdiom)
}

func Translate(line string, idiom string, replacements []string) string {
	value := ""

	if val, ok := translationKeys[line+"_"+idiom]; ok {
		value = val
	}

	if value == "" || replacements == nil {
		return value
	}

	for index, replace := range replacements {
		value = strings.Replace(value, "{"+strconv.Itoa(index)+"}", replace, 1)
	}

	return value
}
