package locales

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"log"
)

//go:embed packs/*.yaml
var LocaleFS embed.FS

type I18n struct {
	dict map[string]map[string]string
}

type Locale struct {
	Name string
	Key  string
}

var AvailableLocales []Locale

func NewLocalizer() (*I18n, error) {
	i := &I18n{}
	err := i.load()
	if err != nil {
		return nil, err
	}
	for key, l := range i.dict {
		AvailableLocales = append(AvailableLocales, Locale{l["language_name"], key})
	}
	return i, nil
}

func (i *I18n) load() error {
	dir, err := LocaleFS.ReadDir("packs")
	if err != nil {
		return err
	}
	for _, f := range dir {
		if !f.IsDir() {
			data, err := fs.ReadFile(LocaleFS, "packs/"+f.Name())
			if err != nil {
				return fmt.Errorf("failed to read language pack %s: %w", f.Name(), err)
			}
			err = yaml.Unmarshal(data, &i.dict)
			if err != nil {
				return fmt.Errorf("failed to parse language pack %s: %w", f.Name(), err)
			}
			log.Println("loaded language pack:", f.Name())
		}
	}
	return nil
}

func (i *I18n) Get(lang, key string) string {
	if m := i.dict[lang]; m == nil {
		return "error: language not found"
	} else if v, ok := m[key]; ok {
		return v
	}
	log.Printf("localizer: key not found: lang=%s, key=%s\n", lang, key)
	return "error: " + key + " not found"
}

func (i *I18n) CheckLangAvailable(lang string) bool {
	return i.dict[lang] != nil
}
