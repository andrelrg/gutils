package language

const Context string = "example"

func GetLang() map[string]map[string]string {
	lang := map[string]map[string]string{
		"es": {
			"example_lang_1": "Hola {0}!",
			"example_lang_2": "¿Todo bien?",
		},
		"pt_br": {
			"example_lang_1": "Olá {0}!",
			"example_lang_2": "tudo bem?",
		},
		"en_us": {
			"example_lang_1": "Hello {0}!",
			"example_lang_2": "How are you?",
		},
	}

	return lang
}
