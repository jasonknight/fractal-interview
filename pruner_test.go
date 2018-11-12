package pruner
import "testing"
func TestPruner(t *testing.T) {
	txt,err := fileGetContents("./test_data/input.json");
	if err != nil {
		t.Errorf("failed to load input %v",err)
		return
	}
	root,err := ParseTree(txt);
	if err != nil {
		t.Errorf("failed to parse json %v",err)
		return
	}
	cat := root.Themes[0].SubThemes[0].Categories[0]
	ind_filter := func (ind Indicator) bool {
		if ind.Id == 299 {
			return true
		}
		return false
	}
	new_indicators := cat.FilterIndicators(ind_filter)
	if len(new_indicators) != 1 {
		t.Errorf("expected only one indicator")
	}
	sub_theme := root.Themes[0].SubThemes[0]
	sub_filter := func (cat Category) bool {
		ninds := cat.FilterIndicators(ind_filter)
		if len(ninds) > 0 {
			return true
		}
		return false
	}
	new_cats := sub_theme.FilterCategories(sub_filter)
	if len(new_cats) != 1 {
		t.Errorf("expected 1 category")
	}
	theme := root.Themes[0]
	t_filter := func (st SubTheme)bool {
		new_cats := st.FilterCategories(sub_filter)
		if len(new_cats) > 0 {
			return true
		}
		return false
	}
	new_subs := theme.FilterSubThemes(t_filter)
	if len(new_subs) != 1 {
		t.Errorf("expected 1 subtheme")
	}
	new_root := root.Filter(func (th Theme)bool {
		new_subs := th.FilterSubThemes(t_filter)
		if len(new_subs) > 0 {
			return true
		}
		return false
	})
	if len(new_root.Themes) != 1 {
		t.Errorf("expected filtering collection would return 1 theme")
	}
}