package main

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
// So, the main idea is that we will 'filter' the tree from
// the bottom (or leaf) up. So each node will have a filter
// function, and we can pass any arbitrary filter, and it
// will return it's children filtered.
// 
// Ideally this could be generic, but that would take more
// time than I'm willing to invest for the moment.
// 
// Also, I will put the functions directly after the struct
// that they are defined for... it's an arbitrary organization
// but for such a small lib, it makes more sense this way.
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
	// What, again? Actually yes, in this sense, over
	// generalizing can lead to "one function to rule them all"
	// which has drawbacks. 
	// 
	// Over time, these functions are going to get crudded
	// up with hotfixes and hacks due to the
	// unforeseen.
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
// So what to name this guy? Eh, we'll call him a ThemeCollection,
// but this is mainly there to encapsulate the naughty bits
// of iterating over the collection returned by the server.
// If we wanted to wrap this up into a library, this would
// make a good interface with their code as all they need
// to do is pass in a callback to do the filtering
// We can also add onto the ThemeCollection some utility
// functions to further conceal the internals
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
func (tc ThemeCollection) FilterByIndicators(cb func(ind Indicator)bool) ThemeCollection {
	cat_filter := func(cat Category) bool {
		ninds := cat.FilterIndicators(cb)
		if len(ninds) > 0 {
			return true
		}
		return false
	}
	subtheme_filter := func(st SubTheme) bool {
		new_cats := st.FilterCategories(cat_filter)
		if len(new_cats) > 0 {
			return true
		}
		return false
	}
	theme_filter := func(th Theme) bool {
		new_subs := th.FilterSubThemes(subtheme_filter)
		if len(new_subs) > 0 {
			return true
		}
		return false
	}
	return tc.Filter(theme_filter)
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
