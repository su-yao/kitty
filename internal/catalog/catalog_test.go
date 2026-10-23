package catalog

import (
	"strings"
	"testing"
)

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

func TestHintsMarkMissing(t *testing.T) {
	s := New()
	s.SetHints(Hints{
		EventScopes:     map[string]string{"crisis.105": "fleet"},
		EventIDs:        map[string]bool{"crisis.105": true, "crisis.199": true},
		DesignNames:     map[string]bool{"NAME_Titan": true},
		EventsVerified:  true,
		DesignsVerified: true,
	})
	page := s.List(CatEvent, "虫群女王万岁", "4.4.6", 1, 10)
	if page.Total != 1 || page.Items[0].Missing {
		t.Fatalf("known event should stay available: %+v", page.Items)
	}
	if page.Items[0].Scope != "fleet" {
		t.Fatalf("scope should come from game files, got %q", page.Items[0].Scope)
	}
	page = s.List(CatEvent, "异次元入侵者", "4.4.6", 1, 10)
	if page.Total != 1 || !page.Items[0].Missing {
		t.Fatal("unknown event id should be marked missing")
	}
	page = s.List(CatShip, "泰坦", "4.4.6", 1, 50)
	for _, e := range page.Items {
		if e.Command == "add_ship NAME_Titan" && e.Missing {
			t.Fatal("known design should stay available")
		}
		if e.Command == "add_ship NAME_Colossus" && !e.Missing {
			t.Fatal("unknown design should be marked missing")
		}
	}
}

func TestListFilteredHidesMissing(t *testing.T) {
	s := New()
	s.SetHints(Hints{EventIDs: map[string]bool{"crisis.199": true}, EventsVerified: true})
	all := s.List(CatEvent, "", "4.4.6", 1, 800)
	shown := s.ListFiltered(CatEvent, "", "4.4.6", 1, 800, true)
	if shown.Total >= all.Total {
		t.Fatalf("hide_missing should drop rows: %d vs %d", shown.Total, all.Total)
	}
	for _, e := range shown.Items {
		if e.Missing {
			t.Fatalf("missing entry leaked: %s", e.Command)
		}
	}
}

func TestDefaultPlayerTargetAndLore(t *testing.T) {
	s := New()
	page := s.List(CatEvent, "地平线信号", "4.4.6", 1, 10)
	if page.Total != 1 {
		t.Fatal("horizon signal missing")
	}
	e := page.Items[0]
	if e.DefaultTarget != "0" {
		t.Fatalf("ship event default target %q", e.DefaultTarget)
	}
	if !e.Target {
		t.Fatal("horizon signal needs a ship target id")
	}
	if len(e.Targets) == 0 || e.Targets[0].ID != "0" {
		t.Fatalf("want player science ship option, got %+v", e.Targets)
	}
	if e.Lore == "" || e.TargetHint == "" {
		t.Fatalf("missing lore or target hint: %+v", e)
	}
	cash := s.List(CatResource, "cash 5000", "4.4.6", 1, 5)
	if cash.Total == 0 || cash.Items[0].DefaultTarget != "0" || cash.Items[0].Lore == "" {
		t.Fatal("resource commands should default to the player and include lore")
	}
	play := s.List(CatEmpire, "Play 0", "4.4.6", 1, 10)
	if play.Total == 0 {
		t.Fatal("play 0 missing")
	}
	found := false
	for _, item := range play.Items {
		if item.Command == "play 0" {
			found = true
			if item.DefaultTarget != "0" {
				t.Fatalf("play default %q", item.DefaultTarget)
			}
			if len(item.Targets) < 2 {
				t.Fatalf("play should list extra empires: %+v", item.Targets)
			}
		}
	}
	if !found {
		t.Fatal("play 0 not found")
	}
}

func TestAllListedEntriesHavePlayerDefault(t *testing.T) {
	s := New()
	page := s.List("", "", "4.4.6", 1, 800)
	if page.Total == 0 {
		t.Fatal("empty catalog")
	}
	if page.Total > len(page.Items) {
		t.Fatalf("need all items, total=%d got=%d", page.Total, len(page.Items))
	}
	for _, e := range page.Items {
		if e.DefaultTarget != "0" {
			t.Fatalf("%s default target %q", e.Command, e.DefaultTarget)
		}
		if len(e.Targets) == 0 || e.Targets[0].ID != "0" {
			t.Fatalf("%s missing player target: %+v", e.Command, e.Targets)
		}
		if e.Lore == "" {
			t.Fatalf("%s missing lore", e.Command)
		}
		if e.TargetHint == "" {
			t.Fatalf("%s missing target hint", e.Command)
		}
	}
}

func TestMegaSpawnAvoidsStarCenter(t *testing.T) {
	s := New()
	page := s.List(CatMega, "星门 完成", "4.4.6", 1, 10)
	if page.Total == 0 {
		t.Fatal("gateway missing")
	}
	cmd := page.Items[0].Command
	if !strings.Contains(cmd, "spawn_megastructure") {
		t.Fatalf("gateway should spawn, got %s", cmd)
	}
	if !strings.Contains(cmd, "orbit_distance") {
		t.Fatalf("gateway should sit on the rim, got %s", cmd)
	}
	if !strings.Contains(cmd, "orbit_angle = random") {
		t.Fatalf("gateway should pick a random rim angle: %s", cmd)
	}
	if !strings.Contains(page.Items[0].TargetHint, "外围") {
		t.Fatalf("gateway hint should mention the rim: %q", page.Items[0].TargetHint)
	}
	dyson := s.List(CatMega, "戴森球 完成", "4.4.6", 1, 10)
	if dyson.Total == 0 || !strings.Contains(dyson.Items[0].Command, "planet = star") {
		t.Fatalf("dyson should wrap the star: %+v", dyson.Items)
	}
	if strings.Contains(dyson.Items[0].Command, "orbit_distance") {
		t.Fatalf("dyson should not offset from the star: %s", dyson.Items[0].Command)
	}
	hab := s.List(CatMega, "居住站核心", "4.4.6", 1, 10)
	if hab.Total == 0 || hab.Items[0].Scope != "planet" {
		t.Fatalf("habitat should require a planet: %+v", hab.Items)
	}
	if !strings.Contains(hab.Items[0].TargetHint, "行星") {
		t.Fatalf("habitat hint should mention planet: %q", hab.Items[0].TargetHint)
	}
	lathe := s.List(CatMega, "突触凝练机", "4.4.6", 1, 10)
	if lathe.Total == 0 || lathe.Items[0].Scope != "planet" {
		t.Fatalf("synaptic lathe should require a planet: %+v", lathe.Items)
	}
	egg := s.List(CatMega, "噬天巨兽卵", "4.4.6", 1, 10)
	if egg.Total == 0 || egg.Items[0].Scope != "planet" {
		t.Fatalf("behemoth egg should require a planet: %+v", egg.Items)
	}
	if strings.Contains(egg.Items[0].Command, "planet = star") {
		t.Fatalf("behemoth egg should not wrap the star: %s", egg.Items[0].Command)
	}
	ark := s.List(CatMega, "民用方舟", "4.4.6", 1, 10)
	if ark.Total == 0 || !strings.Contains(ark.Items[0].Command, "orbit_distance") {
		t.Fatalf("arkship should sit on the rim: %+v", ark.Items)
	}
	nexus := s.List(CatMega, "科学枢纽 完成", "4.4.6", 1, 10)
	if nexus.Total == 0 || !strings.Contains(nexus.Items[0].Command, "orbit_distance") {
		t.Fatalf("science nexus should sit on the rim: %+v", nexus.Items)
	}
	archive := s.List(CatMega, "博物天枢", "4.4.6", 1, 10)
	if archive.Total == 0 || !strings.Contains(archive.Items[0].Command, "orbit_distance") {
		t.Fatalf("grand archive should sit on the rim: %+v", archive.Items)
	}
	seal := s.List(CatMega, "虚境封印", "4.4.6", 1, 10)
	if seal.Total == 0 || !strings.Contains(seal.Items[0].Command, "orbit_distance") {
		t.Fatalf("shroud seal should sit on the rim: %+v", seal.Items)
	}
}

func TestClassifyMegaFromScript(t *testing.T) {
	kind, dist := ClassifyMega("mystery_gate", "build_outside_gravity_well = yes")
	if kind != megaRim || dist <= 0 {
		t.Fatalf("gravity-well flag should sit on the rim: %s %d", kind, dist)
	}
	cmd := MegaSpawnCommandKind("mystery_gate", kind, dist)
	if !strings.Contains(cmd, "orbit_angle = random") {
		t.Fatalf("rim spawn should randomize angle: %s", cmd)
	}
	kind, _ = ClassifyMega("mystery_ring", "is_orbital_ring = yes")
	if kind != megaPlanet {
		t.Fatalf("orbital ring should need a planet, got %s", kind)
	}
	kind, _ = ClassifyMega("mystery_star", "possible = { is_star = yes }")
	if kind != megaStar {
		t.Fatalf("star trigger should wrap the star, got %s", kind)
	}
	kind, dist = ClassifyMega("mystery_gate2", "possible = { NOR = { is_star = yes } }")
	if kind != megaRim || dist <= 0 {
		t.Fatalf("negated star check should stay on the rim: %s %d", kind, dist)
	}
	kind, _ = ClassifyMega("mystery_habitat", "place_entity_on_planet = yes")
	if kind != megaPlanet {
		t.Fatalf("place_entity_on_planet should need a planet, got %s", kind)
	}
	kind, _ = ClassifyMega("mystery_hab2", "is_habitat = yes")
	if kind != megaPlanet {
		t.Fatalf("is_habitat should need a planet, got %s", kind)
	}
	kind, dist = ClassifyMega("civilian_arkship_megastructure", "")
	if kind != megaRim || dist <= 0 {
		t.Fatalf("arkship should sit on the rim: %s %d", kind, dist)
	}
	kind, dist = ClassifyMega("mystery_bypass", "bypass_type = gateway")
	if kind != megaRim || dist <= 0 {
		t.Fatalf("bypass_type should sit on the rim: %s %d", kind, dist)
	}
	kind, dist = ClassifyMega("mystery_gate3", "build_outside_gravity_well = yes\norbit_distance = 72")
	if kind != megaRim || dist != 72 {
		t.Fatalf("script orbit_distance should be used: %s %d", kind, dist)
	}
}

func TestAllMegaSpawnCommandsHavePlacement(t *testing.T) {
	s := New()
	page := s.List(CatMega, "", "4.5.0", 1, 800)
	if page.Total == 0 {
		t.Fatal("empty mega catalog")
	}
	for _, e := range page.Items {
		if !strings.Contains(e.Command, "spawn_megastructure") {
			continue
		}
		star := strings.Contains(e.Command, "planet = star")
		planet := strings.Contains(e.Command, "planet = this")
		rim := strings.Contains(e.Command, "orbit_distance")
		if planet {
			if e.Scope != "planet" {
				t.Fatalf("%s planet spawn needs planet scope: %s", e.NameZH, e.Command)
			}
			continue
		}
		if star && !rim {
			continue
		}
		if rim && strings.Contains(e.Command, "orbit_angle = random") {
			continue
		}
		t.Fatalf("%s missing placement: %s", e.NameZH, e.Command)
	}
}

func TestScannedModLore(t *testing.T) {
	s := New()
	s.ReplaceScanned([]Entry{
		{NameZH: "吉加恒星工程", Command: MegaSpawnCommand("gigastructure_sun"), Category: CatMod, Since: "3.0.0",
			Source: "mod:Gigastructures", Lore: "来自 MOD Gigastructures。把恒星改造成能量核心。"},
	})
	page := s.List(CatMod, "吉加恒星工程", "4.4.6", 1, 10)
	if page.Total != 1 {
		t.Fatalf("mod entry missing: %+v", page)
	}
	e := page.Items[0]
	if !strings.Contains(e.Lore, "Gigastructures") || !strings.Contains(e.Lore, "能量核心") {
		t.Fatalf("mod lore incomplete: %q", e.Lore)
	}
	if !strings.Contains(e.Command, "orbit_distance") && !strings.Contains(e.Command, "planet = star") {
		t.Fatalf("mod mega should use spawn placement: %s", e.Command)
	}
}

func TestEventScopesAreValid(t *testing.T) {
	valid := map[string]bool{"": true, "country": true, "planet": true, "fleet": true, "ship": true,
		"leader": true, "starbase": true, "first_contact": true, "system": true, "situation": true,
		"selected": true}
	for _, e := range events() {
		if !valid[e.Scope] {
			t.Fatalf("invalid scope %q for %s", e.Scope, e.Command)
		}
		if e.Target && e.Scope != "ship" {
			t.Fatalf("target id is only used by ship events: %s", e.Command)
		}
	}
}
