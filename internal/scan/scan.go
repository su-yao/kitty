package scan

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"stellaris-console-panel/internal/catalog"
)

var (
	reShipName = regexp.MustCompile(`(?m)^\s*name\s*=\s*"?(NAME_[A-Za-z0-9_]+)"?`)
	reShipSize = regexp.MustCompile(`(?m)^\s*ship_size\s*=\s*"?([A-Za-z0-9_]+)"?`)
	reMegaKey  = regexp.MustCompile(`(?m)^([A-Za-z0-9_]+)\s*=\s*\{`)
	reEventID  = regexp.MustCompile(`(?m)^\s*id\s*=\s*"?([A-Za-z0-9_.]+)"?`)
	reLocName  = regexp.MustCompile(`(?m)^\s*([A-Za-z0-9_.]+):[0-9]*\s+"([^"]+)"`)
	reVersion  = regexp.MustCompile(`(?i)([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)
	reEventOK  = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*\.[0-9]+$`)
)

var junkSub = []string{
	"background", "backdrop", "decorat", "ambient", "skybox", "sky_box",
	"portrait", "particle", "animat", "visual", "vfx", "sfx",
	"dummy", "placeholder", "unused", "shader", "icon_frame",
	"room_background", "nebula_gfx", "asteroid_art", "scene_texture",
	"music", "sound_effect", "click_sound", "gui_", "widget",
	"背景", "装饰", "肖像", "粒子", "天空盒", "房间背景",
}

var megaSkip = map[string]bool{
	"inline_script": true, "entity": true, "possible": true, "potential": true,
	"ai": true, "resources": true, "prerequisites": true, "show_in_outliner": true,
	"construction_blocks": true, "upgrade": true, "on_build_complete": true,
	"on_build_start": true, "triggered_planet_modifier": true, "planet_modifier": true,
	"country_modifier": true, "bypass": true, "graphical_culture": true,
}

func GameVersion(root string) string {
	if root == "" {
		return ""
	}
	cands := []string{
		filepath.Join(root, "launcher-settings.json"),
		filepath.Join(root, "README.txt"),
		filepath.Join(root, "checksum_manifest.txt"),
	}
	for _, p := range cands {
		if st, err := os.Stat(p); err != nil || st.IsDir() || st.Size() > 1<<20 {
			continue
		}
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		m := reVersion.FindStringSubmatch(string(b))
		if len(m) < 2 {
			continue
		}
		major, _, _ := strings.Cut(m[1], ".")
		if major != "3" && major != "4" {
			continue
		}
		v := m[1]
		if strings.Count(v, ".") == 1 {
			v += ".0"
		}
		return v
	}
	return ""
}

func Scan(gameRoot, workshopRoot, extraModRoot, version string) ([]catalog.Entry, error) {
	loc := map[string]string{}
	mergeLoc(loc, loadLocNames(gameRoot))
	var out []catalog.Entry
	if gameRoot != "" {
		out = append(out, scanShips(filepath.Join(gameRoot, "common", "global_ship_designs"), "vanilla", version, loc)...)
		out = append(out, scanMegas(filepath.Join(gameRoot, "common", "megastructures"), "vanilla", version, loc)...)
		out = append(out, scanEvents(filepath.Join(gameRoot, "events"), "vanilla", version, loc)...)
	}
	modRoots := collectMods(workshopRoot, extraModRoot)
	for _, root := range modRoots {
		src := "mod:" + filepath.Base(root)
		mergeLoc(loc, loadLocNames(root))
		out = append(out, scanShips(filepath.Join(root, "common", "global_ship_designs"), src, version, loc)...)
		out = append(out, scanMegas(filepath.Join(root, "common", "megastructures"), src, version, loc)...)
		out = append(out, scanEvents(filepath.Join(root, "events"), src, version, loc)...)
		out = append(out, scanShips(filepath.Join(root, "common", "ship_designs"), src, version, loc)...)
	}
	return out, nil
}

func collectMods(workshop, extra string) []string {
	var roots []string
	seen := map[string]bool{}
	add := func(dir string) {
		if dir == "" {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			if e.IsDir() {
				if seen[p] {
					continue
				}
				if looksLikeMod(p) {
					seen[p] = true
					roots = append(roots, p)
				}
				sub, _ := os.ReadDir(p)
				for _, s := range sub {
					if s.IsDir() {
						sp := filepath.Join(p, s.Name())
						if !seen[sp] && looksLikeMod(sp) {
							seen[sp] = true
							roots = append(roots, sp)
						}
					}
				}
			}
		}
	}
	add(workshop)
	add(extra)
	return roots
}

func looksLikeMod(p string) bool {
	if _, err := os.Stat(filepath.Join(p, "descriptor.mod")); err == nil {
		return true
	}
	if st, err := os.Stat(filepath.Join(p, "common")); err == nil && st.IsDir() {
		return true
	}
	if st, err := os.Stat(filepath.Join(p, "events")); err == nil && st.IsDir() {
		return true
	}
	return false
}

func walkText(root string, fn func(path string, data string)) {
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".txt" && ext != ".yml" && ext != ".mod" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		fn(path, string(b))
		return nil
	})
}

func scanShips(root, source, version string, loc map[string]string) []catalog.Entry {
	var out []catalog.Entry
	seen := map[string]bool{}
	walkText(root, func(_ string, data string) {
		for _, m := range reShipName.FindAllStringSubmatchIndex(data, -1) {
			name := data[m[2]:m[3]]
			if seen[name] {
				continue
			}
			seen[name] = true
			lo := m[0] - 240
			if lo < 0 {
				lo = 0
			}
			hi := m[1] + 360
			if hi > len(data) {
				hi = len(data)
			}
			size := ""
			if sm := reShipSize.FindStringSubmatch(data[lo:hi]); len(sm) > 1 {
				size = sm[1]
			}
			zh := locPick(loc, name)
			if !keepShip(name, zh, size) {
				continue
			}
			if zh == "" {
				zh = prettyName(name)
			}
			cat := catalog.CatShip
			if source != "vanilla" {
				cat = catalog.CatMod
			}
			out = append(out, catalog.Entry{
				NameZH:   zh,
				NameEN:   name,
				Command:  "add_ship " + name,
				Category: cat,
				Since:    orVer(version),
				Source:   source,
				Note:     "扫描自 " + source,
			})
		}
	})
	return out
}

func scanMegas(root, source, version string, loc map[string]string) []catalog.Entry {
	var out []catalog.Entry
	seen := map[string]bool{}
	walkText(root, func(_ string, data string) {
		for _, m := range reMegaKey.FindAllStringSubmatch(data, -1) {
			id := m[1]
			if megaSkip[id] || strings.HasPrefix(id, "#") {
				continue
			}
			if seen[id] {
				continue
			}
			seen[id] = true
			zh := locPick(loc, id, "MEGASTRUCTURE_"+id, id+"_MEGASTRUCTURE")
			if !keepMega(id, zh) {
				continue
			}
			if zh == "" {
				zh = prettyName(id)
			}
			cat := catalog.CatMega
			if source != "vanilla" {
				cat = catalog.CatMod
			}
			out = append(out, catalog.Entry{
				NameZH:   zh,
				NameEN:   id,
				Command:  "create_megastructure " + id,
				Category: cat,
				Since:    orVer(version),
				Source:   source,
				Note:     "扫描自 " + source,
			})
		}
	})
	return out
}

func scanEvents(root, source, version string, loc map[string]string) []catalog.Entry {
	var out []catalog.Entry
	seen := map[string]bool{}
	walkText(root, func(_ string, data string) {
		for _, m := range reEventID.FindAllStringSubmatch(data, -1) {
			id := m[1]
			if !strings.Contains(id, ".") {
				continue
			}
			if seen[id] {
				continue
			}
			seen[id] = true
			zh := locPick(loc, id, id+".t", id+".title", id+".desc")
			if !keepEvent(id, zh) {
				continue
			}
			if zh == "" || !hasHan(zh) {
				if t := locPick(loc, id+".t", id+".title"); t != "" {
					zh = t
				} else {
					zh = "事件 " + id
				}
			}
			cat := catalog.CatEvent
			if source != "vanilla" {
				cat = catalog.CatMod
			}
			out = append(out, catalog.Entry{
				NameZH:   zh,
				NameEN:   id,
				Command:  "event " + id,
				Category: cat,
				Since:    orVer(version),
				Source:   source,
				Note:     "扫描自 " + source,
			})
		}
	})
	return out
}

func loadLocNames(root string) map[string]string {
	out := map[string]string{}
	if root == "" {
		return out
	}
	walkText(filepath.Join(root, "localisation"), func(path string, data string) {
		base := strings.ToLower(filepath.Base(path))
		dir := strings.ToLower(filepath.Dir(path))
		zhFile := strings.Contains(base, "simp_chinese") || strings.Contains(dir, "simp_chinese")
		enFile := strings.Contains(base, "english") || strings.Contains(dir, "english")
		if !zhFile && !enFile {
			return
		}
		for _, m := range reLocName.FindAllStringSubmatch(data, -1) {
			key, val := m[1], m[2]
			if val == "" {
				continue
			}
			cur, ok := out[key]
			if !ok {
				out[key] = val
				continue
			}
			if zhFile && !hasHan(cur) {
				out[key] = val
			}
		}
	})
	return out
}

func mergeLoc(dst, src map[string]string) {
	for k, v := range src {
		cur, ok := dst[k]
		if !ok || (hasHan(v) && !hasHan(cur)) {
			dst[k] = v
		}
	}
}

func locPick(loc map[string]string, keys ...string) string {
	var en string
	for _, k := range keys {
		if v, ok := loc[k]; ok && v != "" {
			if hasHan(v) {
				return v
			}
			if en == "" {
				en = v
			}
		}
	}
	return en
}

func keepShip(name, zh, size string) bool {
	blob := strings.ToLower(name + " " + zh + " " + size)
	if isJunk(blob) {
		return false
	}
	if hasHan(zh) {
		return true
	}
	return strings.HasPrefix(name, "NAME_") && len(name) > 6
}

func keepMega(id, zh string) bool {
	blob := strings.ToLower(id + " " + zh)
	if isJunk(blob) {
		return false
	}
	if hasHan(zh) {
		return true
	}
	low := strings.ToLower(id)
	usable := []string{
		"dyson", "ring_world", "habitat", "gateway", "spy_orb", "think_tank",
		"matter_decompressor", "mega_shipyard", "strategic_coordination",
		"interstellar_assembly", "mega_art", "sentry", "lgate", "hyper_relay",
		"quantum_catapult", "orbital_ring", "stellar_cannon", "arc_furnace",
		"dyson_swarm", "megastructure", "gigastructure", "giga_", "aot_",
		"katzen", "compound", "ehof", "birch", "nidavellir", "star_lifter",
		"neutronium_giga", "penrose", "yggdrasil", "behemoth",
	}
	for _, p := range usable {
		if strings.Contains(low, p) {
			return true
		}
	}
	return false
}

func keepEvent(id, zh string) bool {
	low := strings.ToLower(id + " " + zh)
	if isJunk(low) {
		return false
	}
	ns := strings.ToLower(strings.Split(id, ".")[0])
	if ns == "test" || ns == "debug" || ns == "hidden" || ns == "tutorial" || ns == "dummy" {
		return false
	}
	if strings.Contains(low, "test_event") || strings.Contains(low, "debug_event") {
		return false
	}
	if hasHan(zh) {
		return true
	}
	return reEventOK.MatchString(id)
}

func isJunk(s string) bool {
	for _, p := range junkSub {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

func hasHan(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func prettyName(id string) string {
	id = strings.TrimPrefix(id, "NAME_")
	id = strings.ReplaceAll(id, "_", " ")
	return id
}

func orVer(v string) string {
	if v == "" {
		return "3.0.0"
	}
	return v
}
