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
