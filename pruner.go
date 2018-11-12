package pruner

import (
	"encoding/json"
	"io/ioutil"
	//"fmt"
)

type Indicator struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type Category struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Unit       string      `json:"unit"`
	Indicators []Indicator `json:"indicators"`
}

func (c Category) FilterIndicators(cb func(ind Indicator) bool) []Indicator {
	var new_indicators []Indicator
	for _, ind := range c.Indicators {
		if cb(ind) {
			new_indicators = append(new_indicators, ind)
		}
	}
	return new_indicators
}

type SubTheme struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

func (s SubTheme) FilterCategories(cb func(cat Category) bool) []Category {
	var new_categories []Category
	for _, cat := range s.Categories {
		if cb(cat) {
			new_categories = append(new_categories, cat)
		}
	}
	return new_categories

}

type Theme struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	SubThemes []SubTheme `json:"sub_themes"`
}

func (t Theme) FilterSubThemes(cb func(st SubTheme) bool) []SubTheme {
	var new_subthemes []SubTheme
	for _, st := range t.SubThemes {
		if cb(st) {
			new_subthemes = append(new_subthemes, st)
		}
	}
	return new_subthemes
}

type ThemeCollection struct {
	Themes []Theme
}

func (tc ThemeCollection) Filter(cb func(th Theme) bool) ThemeCollection {
	var new_themes []Theme
	for _, th := range tc.Themes {
		if cb(th) {
			new_themes = append(new_themes, th)
		}
	}
	tc.Themes = new_themes
	return tc
}
func fileGetContents(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}
func ParseTree(json_string []byte) (ThemeCollection, error) {
	var root []Theme
	var nodes ThemeCollection
	err := json.Unmarshal(json_string, &root)
	nodes.Themes = root
	return nodes, err
}
func inIntSlice(needle int, haystack []int) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}
