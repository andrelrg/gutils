package language

import (
	"github.com/astrolink/gutils/general"
	"strconv"
	"strings"
)

var translationKeys map[string]string
var contexts []string

func LoadLang(lang map[string]map[string]string, context string, idiom string) {
	testInArray, _ := general.InArray(context, contexts)

	if testInArray == true {
		return
	}

	if len(translationKeys) == 0 {
		translationKeys = make(map[string]string)
	}

	if val, ok := lang[idiom]; ok {
		for key, value := range val {
			translationKeys[key] = value
		}
	}

	contexts = append(contexts, context)
}

func Translate(line string, replacements []string) string {
	value := ""

	if val, ok := translationKeys[line]; ok {
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
