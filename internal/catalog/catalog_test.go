package catalog

import "testing"

func TestCompareVer(t *testing.T) {
	if compareVer("3.0.0", "3.0.0") != 0 {
		t.Fatal("eq")
	}
	if compareVer("3.12.0", "3.0.0") <= 0 {
		t.Fatal("gt")
	}
	if compareVer("4.4.6", "4.5.0") >= 0 {
		t.Fatal("lt")
	}
}

func TestVersionFilter(t *testing.T) {
	s := New()
	page := s.List(CatShip, "女王泰坦", "3.0.0", 1, 50)
	if page.Total != 0 {
		t.Fatalf("3.0 should hide synth queen, got %d", page.Total)
	}
	page = s.List(CatShip, "女王泰坦", "3.12.0", 1, 50)
	if page.Total == 0 {
		t.Fatal("3.12 should show synth queen")
	}
	page = s.List(CatMega, "星际大炮", "4.0.0", 1, 50)
	if page.Total != 0 {
		t.Fatalf("4.0 should hide stellar cannon, got %d", page.Total)
	}
	page = s.List(CatMega, "星际大炮", "4.4.6", 1, 50)
	if page.Total == 0 {
		t.Fatal("4.4.6 should show stellar cannon")
	}
}

func TestSearchCash(t *testing.T) {
	s := New()
	page := s.List("", "cash", "4.4.6", 1, 24)
	if page.Total == 0 {
		t.Fatal("search cash")
	}
	if page.Items[0].Command == "" {
		t.Fatal("empty command")
	}
}

func TestBuiltinIntegrity(t *testing.T) {
	seen := map[string]bool{}
	for _, e := range builtin() {
		if e.NameZH == "" || e.Command == "" || e.Category == "" || e.Since == "" {
			t.Fatalf("incomplete entry: %+v", e)
		}
		key := e.Category + "|" + e.Command
		if seen[key] {
			t.Fatalf("duplicate entry in category %s: %s", e.Category, e.Command)
		}
		seen[key] = true
	}
}

func TestPagination(t *testing.T) {
	s := New()
	p1 := s.List("", "", "4.4.6", 1, 24)
	if p1.TotalPages < 2 {
		t.Fatalf("expect multiple pages, total=%d", p1.Total)
	}
	p2 := s.List("", "", "4.4.6", 2, 24)
	if len(p2.Items) == 0 {
		t.Fatal("page 2 empty")
	}
	if p1.Items[0].Command == p2.Items[0].Command && p1.Items[0].NameZH == p2.Items[0].NameZH {
		t.Fatal("page 1 and 2 identical")
	}
}

func TestScannedSkipsBuiltinCommands(t *testing.T) {
	s := New()
	before := s.List("", "", "4.4.6", 1, 1).Total
	s.ReplaceScanned([]Entry{
		{NameZH: "重复战列舰", Command: "add_ship NAME_Battleship", Category: CatShip, Since: "3.0.0"},
		{NameZH: "重复战列舰副本", Command: "add_ship NAME_Battleship", Category: CatMod, Since: "3.0.0"},
		{NameZH: "新设计", Command: "add_ship NAME_New_Design", Category: CatMod, Since: "3.0.0"},
	})
	after := s.List("", "", "4.4.6", 1, 1).Total
	if after != before+1 {
		t.Fatalf("want exactly one new entry, before=%d after=%d", before, after)
	}
	if s.ScannedCount() != 1 {
		t.Fatalf("want 1 scanned entry, got %d", s.ScannedCount())
	}
}

func TestListSortsByCategoryOrder(t *testing.T) {
	s := New()
	page := s.List("", "", "4.4.6", 1, 500)
	if len(page.Items) < 2 {
		t.Fatal("not enough entries")
	}
	for i := 1; i < len(page.Items); i++ {
		prev, cur := page.Items[i-1], page.Items[i]
		if prev.Category == cur.Category {
			continue
		}
		if catOrderOf(prev.Category) > catOrderOf(cur.Category) {
			t.Fatalf("category order broken: %s before %s", prev.Category, cur.Category)
		}
	}
}

func TestBuiltinVersionsAndCategories(t *testing.T) {
	knownVer := map[string]bool{}
	for _, v := range Versions {
		knownVer[v] = true
	}
	knownCat := map[string]bool{}
	for _, c := range Categories {
		knownCat[c.ID] = true
	}
	for _, e := range builtin() {
		if !knownVer[e.Since] {
			t.Fatalf("unknown since %q for %s", e.Since, e.Command)
		}
		if e.Until != "" && !knownVer[e.Until] {
			t.Fatalf("unknown until %q for %s", e.Until, e.Command)
		}
		if !knownCat[e.Category] {
			t.Fatalf("unknown category %q for %s", e.Category, e.Command)
		}
	}
}
