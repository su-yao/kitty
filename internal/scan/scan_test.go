package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"stellaris-console-panel/internal/catalog"
)

func TestScanShipsAndMegas(t *testing.T) {
	root := t.TempDir()
	shipDir := filepath.Join(root, "common", "global_ship_designs")
	megaDir := filepath.Join(root, "common", "megastructures")
	evtDir := filepath.Join(root, "events")
	locDir := filepath.Join(root, "localisation", "simp_chinese")
	_ = os.MkdirAll(shipDir, 0755)
	_ = os.MkdirAll(megaDir, 0755)
	_ = os.MkdirAll(evtDir, 0755)
	_ = os.MkdirAll(locDir, 0755)
	_ = os.WriteFile(filepath.Join(locDir, "ships_l_simp_chinese.yml"), []byte(`
l_simp_chinese:
 NAME_Synth_Queen_Titan:0 "女王泰坦"
 NAME_Synth_Queen_Titan_DESC:0 "塞塔娜的纳米泰坦，来自机械纪元。"
 giga.1234:0 "吉加事件"
 giga.1234.desc:0 "吉加结构展开，恒星周围出现巨环。"
 gigastructure_sun:0 "吉加恒星工程"
 gigastructure_sun_DESC:0 "把恒星改造成能量核心。"
 mystery_gate:0 "神秘星门"
 mystery_gate_DESC:0 "在星系外围生成，轨道角随机。"
`), 0644)
	_ = os.WriteFile(filepath.Join(shipDir, "demo.txt"), []byte(`
my_design = {
	name = "NAME_Synth_Queen_Titan"
	ship_size = titan
}
bg_ship = {
	name = "NAME_Background_Decor"
	ship_size = background_station
}
`), 0644)
	_ = os.WriteFile(filepath.Join(megaDir, "demo.txt"), []byte(`
gigastructure_sun = {
	entity = "foo"
	possible = { is_star = yes }
}
mystery_gate = {
	entity = "gate"
	build_outside_gravity_well = yes
}
skybox_background = {
	entity = "bar"
}
inline_script = {
}
`), 0644)
	_ = os.WriteFile(filepath.Join(evtDir, "demo.txt"), []byte(`
country_event = {
	id = giga.1234
}
country_event = {
	id = test.1
}
`), 0644)
	items, err := Scan(root, "", "", "4.4.6")
	if err != nil {
		t.Fatal(err)
	}
	foundShip, foundMega, foundEvt, foundRim := false, false, false, false
	for _, e := range items {
		if e.Command == "add_ship NAME_Synth_Queen_Titan" {
			foundShip = true
			if e.NameZH != "女王泰坦" {
				t.Fatalf("want chinese name, got %q", e.NameZH)
			}
			if !strings.Contains(e.Lore, "纳米泰坦") {
				t.Fatalf("ship lore missing localisation: %q", e.Lore)
			}
		}
		if e.Command == "add_ship NAME_Background_Decor" {
			t.Fatal("decorative ship should be filtered")
		}
		if strings.Contains(e.Command, "skybox_background") {
			t.Fatal("background mega should be filtered")
		}
		if e.Command == "event test.1" {
			t.Fatal("test event should be filtered")
		}
		if strings.Contains(e.Command, "type = gigastructure_sun") {
			foundMega = true
			if !strings.Contains(e.Lore, "能量核心") {
				t.Fatalf("mega lore missing localisation: %q", e.Lore)
			}
			if !strings.Contains(e.Command, "planet = star") {
				t.Fatalf("star mega should wrap the star: %s", e.Command)
			}
			if strings.Contains(e.Command, "orbit_distance") {
				t.Fatalf("star mega should not offset: %s", e.Command)
			}
			if !strings.Contains(e.Lore, "恒星") {
				t.Fatalf("star mega lore should mention the star: %q", e.Lore)
			}
		}
		if strings.Contains(e.Command, "type = mystery_gate") {
			foundRim = true
			if !strings.Contains(e.Command, "orbit_distance") || !strings.Contains(e.Command, "orbit_angle = random") {
				t.Fatalf("rim mega should randomize orbit: %s", e.Command)
			}
			if !strings.Contains(e.Lore, "外围") {
				t.Fatalf("rim mega lore should mention the rim: %q", e.Lore)
			}
		}
		if e.Command == "event giga.1234" {
			foundEvt = true
			if e.Since != "3.0.0" {
				t.Fatalf("scanned entry should be visible on every version, got since=%q", e.Since)
			}
			if !strings.Contains(e.Lore, "巨环") {
				t.Fatalf("event lore missing localisation: %q", e.Lore)
			}
		}
	}
	if !foundShip || !foundMega || !foundEvt || !foundRim {
		t.Fatalf("scan miss ship=%v mega=%v evt=%v rim=%v n=%d", foundShip, foundMega, foundEvt, foundRim, len(items))
	}
}

func TestKeepFilters(t *testing.T) {
	if keepShip("NAME_Background_Decor", "背景装饰", "station") {
		t.Fatal("junk ship")
	}
	if !keepShip("NAME_Synth_Queen_Titan", "女王泰坦", "titan") {
		t.Fatal("queen titan")
	}
	if keepMega("skybox_background", "") {
		t.Fatal("skybox")
	}
	if !keepMega("dyson_sphere_5", "戴森球") {
		t.Fatal("dyson")
	}
	if !keepMega("cosmogenesis_needle", "") {
		t.Fatal("cosmogenesis")
	}
	if !keepMega("civilian_arkship_megastructure", "") {
		t.Fatal("arkship")
	}
	if !keepMega("grand_archive_0", "") {
		t.Fatal("grand archive")
	}
	if !keepMega("orbital_arc_furnace_4", "") {
		t.Fatal("arc furnace")
	}
	if !keepMega("spy_orb_4", "") {
		t.Fatal("sentry array")
	}
	if !keepMega("lgate_base", "") {
		t.Fatal("l-gate")
	}
	if !keepMega("think_tank_3", "") {
		t.Fatal("science nexus")
	}
	if !keepMega("hyper_relay", "") {
		t.Fatal("hyper relay")
	}
	if !keepMega("habitat_central_complex", "") {
		t.Fatal("habitat")
	}
	if !keepMega("gateway_final", "") {
		t.Fatal("gateway")
	}
	if !keepMega("behemoth_egg", "") {
		t.Fatal("behemoth")
	}
	if !keepMega("galactic_crucible_4", "") {
		t.Fatal("crucible")
	}
	if !keepMega("dyson_gun_3", "") {
		t.Fatal("stellar cannon")
	}
	if !keepMega("shroud_seal", "") {
		t.Fatal("shroud seal")
	}
	if !keepMega("deep_space_citadel_0", "") {
		t.Fatal("citadel")
	}
	if !keepMega("waystation_megastructure", "") {
		t.Fatal("waystation")
	}
	if keepEvent("test.1", "") {
		t.Fatal("test event")
	}
	if keepEvent("giga.1234", "") {
		t.Fatal("event without a localised title should be dropped")
	}
	if !keepEvent("giga.1234", "吉加事件") {
		t.Fatal("localised event should be kept")
	}
	if keepShip("NAME_Portrait_Dummy", "肖像", "corvette") {
		t.Fatal("portrait ship")
	}
}

func TestScanReportScopesAndDesigns(t *testing.T) {
	root := t.TempDir()
	evtDir := filepath.Join(root, "events")
	locDir := filepath.Join(root, "localisation", "simp_chinese")
	shipDir := filepath.Join(root, "common", "global_ship_designs")
	for _, d := range []string{evtDir, locDir, shipDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	_ = os.WriteFile(filepath.Join(locDir, "e.yml"), []byte(`
l_simp_chinese:
 test.1:0 "国家事件"
 test.2:0 "星球事件"
 test.3:0 "舰船事件"
 test.4:0 "局势事件"
`), 0644)
	_ = os.WriteFile(filepath.Join(evtDir, "demo.txt"), []byte(`
country_event = {
	id = test.1
}
planet_event = {
	id = test.2
}
ship_event = {
	id = test.3
}
situation_event = {
	id = test.4
}
`), 0644)
	_ = os.WriteFile(filepath.Join(shipDir, "demo.txt"), []byte(`
d1 = {
	name = "NAME_Titan"
	ship_size = titan
}
`), 0644)
	r, err := ScanReport(root, "", "", "4.4.6")
	if err != nil {
		t.Fatal(err)
	}
	if !r.EventsVerified || !r.DesignsVerified {
		t.Fatalf("expected verification, got %+v", r)
	}
	if r.EventScopes["test.1"] != "country" || r.EventScopes["test.2"] != "planet" ||
		r.EventScopes["test.3"] != "ship" || r.EventScopes["test.4"] != "situation" {
		t.Fatalf("wrong scopes: %+v", r.EventScopes)
	}
	if !r.DesignNames["NAME_Titan"] {
		t.Fatalf("missing design names: %+v", r.DesignNames)
	}
}

func TestKeepMegaCoversBuiltinIDs(t *testing.T) {
	s := catalog.New()
	page := s.List(catalog.CatMega, "", "4.5.0", 1, 800)
	n := 0
	for _, e := range page.Items {
		id := megaIDFromCommand(e.Command)
		if id == "" {
			continue
		}
		n++
		if !keepMega(id, "") {
			t.Fatalf("keepMega dropped builtin id %s (%s)", id, e.NameZH)
		}
	}
	if n < 50 {
		t.Fatalf("too few mega ids: %d", n)
	}
}

func megaIDFromCommand(cmd string) string {
	const key = "type = "
	i := strings.Index(cmd, key)
	if i < 0 {
		return ""
	}
	rest := strings.TrimSpace(cmd[i+len(key):])
	id := rest
	for j, r := range rest {
		if r == ' ' || r == '}' {
			id = rest[:j]
			break
		}
	}
	return id
}

func TestScanReportWithoutGameRoot(t *testing.T) {
	mod := t.TempDir()
	_ = os.MkdirAll(filepath.Join(mod, "events"), 0755)
	r, err := ScanReport("", mod, "", "4.4.6")
	if err != nil {
		t.Fatal(err)
	}
	if r.EventsVerified || r.DesignsVerified {
		t.Fatal("mod-only scan must not verify vanilla entries")
	}
}
