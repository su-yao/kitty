package catalog

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Entry struct {
	ID       string `json:"id"`
	NameZH   string `json:"name_zh"`
	NameEN   string `json:"name_en"`
	Command  string `json:"command"`
	Category string `json:"category"`
	Since    string `json:"since"`
	Until    string `json:"until,omitempty"`
	Source   string `json:"source"`
	Note     string `json:"note,omitempty"`
}

type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type Page struct {
	Items      []Entry `json:"items"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
	Category   string  `json:"category"`
}

const (
	CatResource  = "resource"
	CatShip      = "ship"
	CatCrisis    = "crisis_ship"
	CatLeviathan = "leviathan"
	CatMega      = "mega"
	CatEvent     = "event"
	CatTech      = "tech"
	CatRelic     = "relic"
	CatPlanet    = "planet"
	CatEmpire    = "empire"
	CatDebug     = "debug"
	CatMod       = "mod"
)

var Categories = []Category{
	{CatResource, "资源经济", 1},
	{CatShip, "建造舰船", 2},
	{CatCrisis, "天灾舰船", 3},
	{CatLeviathan, "利维坦", 4},
	{CatMega, "巨构", 5},
	{CatEvent, "事件危机", 6},
	{CatTech, "科技飞升", 7},
	{CatRelic, "遗珍标本", 8},
	{CatPlanet, "行星人口", 9},
	{CatEmpire, "帝国外交", 10},
	{CatDebug, "调试开关", 11},
	{CatMod, "MOD扫描", 12},
}

var catOrder = func() map[string]int {
	m := make(map[string]int, len(Categories))
	for _, c := range Categories {
		m[c.ID] = c.Order
	}
	return m
}()

var Versions = []string{
	"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0", "3.5.0", "3.6.0", "3.7.0",
	"3.8.0", "3.9.0", "3.10.0", "3.11.0", "3.12.0", "3.13.0", "3.14.0",
	"4.0.0", "4.1.0", "4.2.0", "4.3.0", "4.4.0", "4.4.6", "4.5.0",
}

type Store struct {
	mu      sync.RWMutex
	vanilla []Entry
	scanned []Entry
}

func New() *Store {
	s := &Store{vanilla: builtin()}
	s.reindex()
	return s
}

func (s *Store) reindex() {
	for i := range s.vanilla {
		if s.vanilla[i].ID == "" {
			s.vanilla[i].ID = makeID(s.vanilla[i])
		}
		if s.vanilla[i].Source == "" {
			s.vanilla[i].Source = "vanilla"
		}
		if s.vanilla[i].Since == "" {
			s.vanilla[i].Since = "3.0.0"
		}
	}
}

func makeID(e Entry) string {
	base := e.Category + "-" + e.Command
	base = strings.ReplaceAll(base, " ", "-")
	base = strings.ReplaceAll(base, "=", "-")
	return strings.ToLower(base)
}

func (s *Store) ReplaceScanned(items []Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Scanning the vanilla game re-discovers designs that are already in the
	// built-in catalog. Skipping them keeps one button per real command.
	builtinCmd := make(map[string]bool, len(s.vanilla))
	for _, e := range s.vanilla {
		builtinCmd[e.Command] = true
	}
	seen := map[string]bool{}
	out := make([]Entry, 0, len(items))
	for _, e := range items {
		if e.Command == "" || builtinCmd[e.Command] {
			continue
		}
		if e.Source == "" {
			e.Source = "mod"
		}
		if e.Category == "" {
			e.Category = CatMod
		}
		if e.Since == "" {
			e.Since = "3.0.0"
		}
		e.ID = makeID(e)
		key := e.Category + "|" + e.Command
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, e)
	}
	s.scanned = out
}

func (s *Store) ScannedCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.scanned)
}

func (s *Store) List(category, query, version string, page, pageSize int) Page {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 24
	}
	if pageSize > 800 {
		pageSize = 800
	}
	query = strings.ToLower(strings.TrimSpace(query))
	items := append([]Entry{}, s.vanilla...)
	items = append(items, s.scanned...)
	filtered := make([]Entry, 0, len(items))
	for _, e := range items {
		if !inVersion(e, version) {
			continue
		}
		if category != "" && category != "all" && e.Category != category {
			continue
		}
		if query != "" {
			blob := strings.ToLower(e.NameZH + " " + e.NameEN + " " + e.Command + " " + e.Note)
			if !strings.Contains(blob, query) {
				continue
			}
		}
		filtered = append(filtered, e)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Category != filtered[j].Category {
			return catOrderOf(filtered[i].Category) < catOrderOf(filtered[j].Category)
		}
		return filtered[i].NameZH < filtered[j].NameZH
	})
	total := len(filtered)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return Page{
		Items:      filtered[start:end],
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Category:   category,
	}
}

func catOrderOf(id string) int {
	if n, ok := catOrder[id]; ok {
		return n
	}
	return len(catOrder) + 1
}

func inVersion(e Entry, version string) bool {
	if version == "" {
		return true
	}
	if e.Since != "" && compareVer(version, e.Since) < 0 {
		return false
	}
	if e.Until != "" && compareVer(version, e.Until) > 0 {
		return false
	}
	return true
}

func compareVer(a, b string) int {
	as := parseVer(a)
	bs := parseVer(b)
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		av, bv := 0, 0
		if i < len(as) {
			av = as[i]
		}
		if i < len(bs) {
			bv = bs[i]
		}
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
	}
	return 0
}

func parseVer(s string) []int {
	s = strings.TrimPrefix(strings.ToLower(s), "v")
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
	out := make([]int, 0, 3)
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			break
		}
		out = append(out, n)
		if len(out) == 3 {
			break
		}
	}
	for len(out) < 3 {
		out = append(out, 0)
	}
	return out
}

func e(zh, en, command, cat, since, note string) Entry {
	return Entry{
		NameZH:   zh,
		NameEN:   en,
		Command:  command,
		Category: cat,
		Since:    since,
		Source:   "vanilla",
		Note:     note,
	}
}

func builtin() []Entry {
	out := make([]Entry, 0, 900)
	out = append(out, resources()...)
	out = append(out, ships()...)
	out = append(out, crisisShips()...)
	out = append(out, leviathans()...)
	out = append(out, megas()...)
	out = append(out, events()...)
	out = append(out, techs()...)
	out = append(out, relics()...)
	out = append(out, planets()...)
	out = append(out, empires()...)
	out = append(out, debugs()...)
	return out
}
