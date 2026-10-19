package scan

import (
	"os"
	"path/filepath"
	"testing"
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
 giga.1234:0 "吉加事件"
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
	foundShip, foundMega, foundEvt := false, false, false
	for _, e := range items {
		if e.Command == "add_ship NAME_Synth_Queen_Titan" {
			foundShip = true
			if e.NameZH != "女王泰坦" {
				t.Fatalf("want chinese name, got %q", e.NameZH)
			}
		}
		if e.Command == "add_ship NAME_Background_Decor" {
			t.Fatal("decorative ship should be filtered")
		}
		if e.Command == "create_megastructure skybox_background" {
			t.Fatal("background mega should be filtered")
		}
		if e.Command == "event test.1" {
			t.Fatal("test event should be filtered")
		}
		if e.Command == "create_megastructure gigastructure_sun" {
			foundMega = true
		}
		if e.Command == "event giga.1234" {
			foundEvt = true
			if e.Since != "3.0.0" {
				t.Fatalf("scanned entry should be visible on every version, got since=%q", e.Since)
			}
		}
	}
	if !foundShip || !foundMega || !foundEvt {
		t.Fatalf("scan miss ship=%v mega=%v evt=%v n=%d", foundShip, foundMega, foundEvt, len(items))
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
	if r.EventScopes["test.1"] != "country" || r.EventScopes["test.2"] != "planet" || r.EventScopes["test.3"] != "ship" {
		t.Fatalf("wrong scopes: %+v", r.EventScopes)
	}
	if !r.DesignNames["NAME_Titan"] {
		t.Fatalf("missing design names: %+v", r.DesignNames)
	}
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
