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
	// Scope is the object that must be selected before the command works:
	// country, planet, fleet, ship, leader, pop, situation, first_contact, system or selected.
	Scope string `json:"scope,omitempty"`
	// Target means a numeric id from the picker is written into the command.
	Target bool `json:"target,omitempty"`
	// TargetMode is append, first, second or last.
	TargetMode string `json:"target_mode,omitempty"`
	// DLC names the required expansion, for example 天启 or 机械纪元.
	DLC string `json:"dlc,omitempty"`
	// Missing is set after a scan when the event id or ship design used by a
	// built-in entry does not exist in the current game files.
	Missing bool `json:"missing,omitempty"`
	// DefaultTarget is the numeric id used when the command needs a target.
	// Player country and the starting science ship are both 0.
	DefaultTarget string `json:"default_target,omitempty"`
	// Targets lists the player default plus other valid ids the user can pick.
	Targets []TargetOpt `json:"targets,omitempty"`
	// Lore is the Stellaris background shown when hovering a command button.
	Lore string `json:"lore,omitempty"`
	// TargetHint explains which in-game object the command acts on.
	TargetHint string `json:"target_hint,omitempty"`
	// PreCommands are typed and confirmed before Command, in order.
	PreCommands []string `json:"pre_commands,omitempty"`
}

type TargetOpt struct {
	ID    string `json:"id"`
	Label string `json:"label"`
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

type Hints struct {
	EventScopes     map[string]string
	EventIDs        map[string]bool
	DesignNames     map[string]bool
	EventsVerified  bool
	DesignsVerified bool
}

type Store struct {
	mu      sync.RWMutex
	vanilla []Entry
	scanned []Entry
	hints   Hints
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

// SetHints records what the current game files contain, so built-in entries
// can be marked as unavailable instead of silently doing nothing.
func (s *Store) SetHints(h Hints) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hints = h
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
	return s.ListFiltered(category, query, version, page, pageSize, false)
}

func (s *Store) ListFiltered(category, query, version string, page, pageSize int, hideMissing bool) Page {
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
		e = s.decorate(e)
		if hideMissing && e.Missing {
			continue
		}
		if !inVersion(e, version) {
			continue
		}
		if category != "" && category != "all" && e.Category != category {
			continue
		}
		if query != "" {
			blob := strings.ToLower(e.NameZH + " " + e.NameEN + " " + e.Command + " " + e.Note + " " + e.Lore)
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

func (s *Store) decorate(e Entry) Entry {
	if e.Source == "vanilla" {
		if e.Category == CatEvent {
			if id, ok := eventID(e.Command); ok {
				if sc, found := s.hints.EventScopes[id]; found {
					e.Scope = sc
				}
				if s.hints.EventsVerified && !s.hints.EventIDs[id] {
					e.Missing = true
				}
			}
		}
		if name, ok := shipDesign(e.Command); ok {
			// add_ship places the new ship next to the selected object, so
			// something must be selected first.
			if e.Scope == "" {
				e.Scope = "selected"
			}
			// Player-created designs are not in the game files, so only preset
			// NAME_ designs can be verified.
			if strings.HasPrefix(name, "NAME_") && s.hints.DesignsVerified &&
				!s.hints.DesignNames[name] && !s.hints.DesignNames[strings.TrimPrefix(name, "NAME_")] {
				e.Missing = true
			}
		}
	} else if _, ok := shipDesign(e.Command); ok && e.Scope == "" {
		e.Scope = "selected"
	}
	if id, ok := megaType(e.Command); ok {
		if strings.HasPrefix(strings.TrimSpace(e.Command), "create_megastructure ") {
			e.Command = MegaSpawnCommand(id)
		}
		if e.Scope == "" {
			e.Scope = MegaScope(id)
		}
	}
	if e.Scope == "" {
		e.Scope = inferScope(e.Command)
	}
	if e.TargetMode == "" {
		e.TargetMode = inferTargetMode(e.Command, e.Scope, e.Target)
	}
	if e.TargetMode != "" {
		e.Target = true
	} else if e.Scope == "ship" {
		e.Target = true
		e.TargetMode = "append"
	} else if e.Scope != "" && e.Scope != "selected" {
		e.Target = false
	}
	e = fillTargets(e)
	e = applyDLC(e)
	if _, ok := megaType(e.Command); ok {
		e.TargetHint = megaTargetHint(e.Command)
	}
	switch e.Command {
	case "own":
		e.TargetHint = "先在游戏中左键选中舰队、星球或巨构。命令会把所有权转给玩家。"
	case "effect remove_megastructure = this":
		e.TargetHint = "先在游戏中选中要拆除的巨构。"
	case "finish_rift_stage":
		e.TargetHint = "先在游戏中选中星界裂隙。"
	case "spawn_astral_rift":
		e.TargetHint = "先在游戏中选中目标星系。裂隙会生成在该星系。"
	case "effect enrage_behemoth":
		e.TargetHint = "先选中玩家舰队、恒星基地或星球。"
	}
	if e.Lore == "" {
		e.Lore = loreFor(e)
	}
	if len(e.PreCommands) == 0 {
		e.PreCommands = PreCommandsFor(e.Command)
	}
	return e
}

func fillTargets(e Entry) Entry {
	if e.TargetHint != "" && len(e.Targets) > 0 {
		return e
	}
	player := TargetOpt{ID: "0", Label: "玩家"}
	switch {
	case e.Scope == "pop":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{{ID: "0", Label: "人口 0"}}
		e.TargetHint = "参数是人口 ID。控制台输入 debugtooltip 后，把鼠标移到人口上可查看 ID。右上角可改。"
	case e.Scope == "leader":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{{ID: "0", Label: "领袖 0"}}
		e.TargetHint = "先在游戏中选中玩家领袖。带数字的命令里，右上角 ID 会写入命令。"
	case e.Target && (e.Scope == "ship" || strings.HasPrefix(e.Command, "event ")):
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{{ID: "0", Label: "玩家开局科研船"}}
		e.TargetHint = "默认玩家开局科研船（ID 0）。控制台输入 debugtooltip 后，把鼠标移到舰船上可查看其它舰船 ID。"
	case e.Scope == "fleet":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先在游戏中选中玩家舰队，再发送。命令本身不附加数字 ID。"
	case e.Scope == "planet":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先在游戏中选中玩家星球，再发送。命令本身不附加数字 ID。"
	case e.Scope == "ship":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{{ID: "0", Label: "玩家开局科研船"}}
		e.TargetHint = "先在游戏中选中玩家舰船。需要目标 ID 的事件会默认使用开局科研船 0。"
	case e.Scope == "starbase":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先在游戏中选中玩家恒星基地。"
	case e.Scope == "first_contact":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先在游戏中选中进行中的首次接触。"
	case e.Scope == "system":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先在游戏中选中目标星系。星门等会放在星系外围，戴森球等会套在恒星上。"
	case e.Scope == "situation":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先在局势界面选中对应局势。"
	case e.Scope == "selected":
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "先选中玩家舰队、星球或恒星，新对象会生成在该单位旁边。"
	case e.Target || needsCountryID(e.Command):
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{
			player,
			{ID: "1", Label: "帝国 1"},
			{ID: "2", Label: "帝国 2"},
		}
		e.TargetHint = "默认玩家帝国（ID 0）。右上角可改选其它帝国 ID，数字会写入命令。"
	default:
		e.DefaultTarget = "0"
		e.Targets = []TargetOpt{player}
		e.TargetHint = "默认作用于玩家帝国（ID 0），无需额外目标 ID。"
	}
	return e
}

func PreCommandsFor(cmd string) []string {
	cmd = strings.TrimSpace(cmd)
	if v, ok := preByCommand[cmd]; ok {
		out := make([]string, len(v))
		copy(out, v)
		return out
	}
	if strings.HasPrefix(cmd, "activate_relic ") {
		id := strings.TrimSpace(strings.TrimPrefix(cmd, "activate_relic "))
		if id != "" && id != "all" {
			return []string{"add_relic " + id}
		}
	}
	return nil
}

var preByCommand = map[string][]string{
	"event crisis.1200":                 {"event crisis.1000"},
	"event utopia.3320":                 {"event utopia.3308"},
	"event graygoo.1":                   {"event distar.11000"},
	"event graygoo.100":                 {"event distar.11000"},
	"event distar.13000":                {"event distar.11000"},
	"event fallen_empires_awakening.1":  {"play 1"},
	"event fallen_machine_empire.1":     {"play 1"},
	"event marauder.500":                {"play 1"},
	"event fallen_empires_awakening.10": {"play 1"},
	"event grand_archive.2300":          {"play 1"},
}

func inferScope(cmd string) string {
	switch {
	case cmd == "own" || cmd == "effect remove_megastructure = this" ||
		cmd == "finish_rift_stage" || cmd == "effect enrage_behemoth":
		return "selected"
	case cmd == "invincible" || strings.HasPrefix(cmd, "damage ") ||
		cmd == "effect destroy_fleet = this" || cmd == "break_fleet_contract":
		return "fleet"
	case cmd == "complete_first_contact" || strings.HasPrefix(cmd, "add_first_contact_clues "):
		return "first_contact"
	case strings.Contains(cmd, "add_situation_progress") || strings.Contains(cmd, "destroy_situation"):
		return "situation"
	case cmd == "spawn_astral_rift":
		return "system"
	case strings.HasPrefix(cmd, "kill_pop "):
		return "pop"
	case cmd == "kill_ruler" || strings.HasPrefix(cmd, "kill_leader ") ||
		strings.HasPrefix(cmd, "skills ") || strings.HasPrefix(cmd, "add_trait_leader "):
		return "leader"
	case cmd == "branchoffice" || cmd == "finish_arc_stage" || planetScoped(cmd):
		return "planet"
	}
	return ""
}

func planetScoped(cmd string) bool {
	switch cmd {
	case "finish_terraform", "grow_pops", "control":
		return true
	}
	for _, p := range []string{
		"planet_class ", "planet_size ", "planet_happiness ", "planet_ascension_tier ",
		"colonize ", "add_pops ", "add_anomaly ",
		"effect destroy_colony", "effect add_planet_devastation",
		"effect add_district", "effect add_building",
		"effect add_deposit", "effect remove_deposit",
		"effect create_archaeological_site", "effect set_planet_flag",
	} {
		if strings.HasPrefix(cmd, p) {
			return true
		}
	}
	return false
}

func needsCountryID(cmd string) bool {
	for _, p := range []string{
		"play ", "annex", "force_integrate", "add_opinion ", "add_trust ",
		"add_intel ", "add_loyalty ", "add_spynetwork_value ", "add_subject_xp ",
		"war ", "casusbelli ", "surrender ", "kill_country",
	} {
		if strings.HasPrefix(cmd, p) || cmd == strings.TrimSpace(p) {
			return true
		}
	}
	return false
}

func inferTargetMode(cmd, scope string, target bool) string {
	switch {
	case strings.HasPrefix(cmd, "kill_pop "):
		return "last"
	case strings.HasPrefix(cmd, "kill_leader "), strings.HasPrefix(cmd, "add_trait_leader "):
		return "first"
	case strings.HasPrefix(cmd, "play "):
		return "last"
	case strings.HasPrefix(cmd, "add_intel "), strings.HasPrefix(cmd, "add_loyalty "),
		strings.HasPrefix(cmd, "add_spynetwork_value "), strings.HasPrefix(cmd, "add_subject_xp "),
		strings.HasPrefix(cmd, "casusbelli "), strings.HasPrefix(cmd, "surrender "):
		return "first"
	case strings.HasPrefix(cmd, "add_opinion "), strings.HasPrefix(cmd, "add_trust "),
		strings.HasPrefix(cmd, "war "):
		return "second"
	case target && strings.HasPrefix(cmd, "event "):
		return "append"
	case target && scope == "ship":
		return "append"
	}
	return ""
}

func ApplyTarget(command, mode, id string) string {
	command = strings.TrimSpace(command)
	if command == "" || mode == "" || !isDigits(id) {
		return command
	}
	f := strings.Fields(command)
	if len(f) == 0 {
		return command
	}
	switch mode {
	case "append":
		if len(f) >= 3 && isDigits(f[len(f)-1]) {
			f[len(f)-1] = id
			return strings.Join(f, " ")
		}
		return command + " " + id
	case "first":
		for i := 1; i < len(f); i++ {
			if isDigits(f[i]) {
				f[i] = id
				return strings.Join(f, " ")
			}
		}
	case "second":
		n := 0
		for i := 1; i < len(f); i++ {
			if isDigits(f[i]) {
				n++
				if n == 2 {
					f[i] = id
					return strings.Join(f, " ")
				}
			}
		}
	case "last":
		for i := len(f) - 1; i >= 1; i-- {
			if isDigits(f[i]) {
				f[i] = id
				return strings.Join(f, " ")
			}
		}
	}
	return command
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func applyDLC(e Entry) Entry {
	if e.DLC != "" {
		return e
	}
	if d, ok := dlcByCommand[e.Command]; ok {
		e.DLC = d
		return e
	}
	if id, ok := eventID(e.Command); ok {
		if d, ok := dlcByEvent[id]; ok {
			e.DLC = d
		}
	}
	if name, ok := shipDesign(e.Command); ok {
		if d, ok := dlcByShip[name]; ok {
			e.DLC = d
		}
	}
	if id, ok := megaType(e.Command); ok {
		if d, ok := dlcByMega[id]; ok {
			e.DLC = d
		}
	}
	return e
}

var dlcByCommand = map[string]string{
	"research_technology tech_titans":     "天启",
	"research_technology tech_colossus":   "天启",
	"research_technology tech_juggernaut": "联邦",
	"activate_ascension_perk ap_become_the_crisis": "天灾",
	"activate_ascension_perk ap_cosmogenesis":      "机械纪元",
	"activate_ascension_perk ap_nanotech":          "机械纪元",
	"activate_ascension_perk ap_modularity":        "机械纪元",
	"activate_ascension_perk ap_behemoth_fury":     "万物竞发",
	"activate_ascension_perk ap_galactic_hyperthermia": "炼狱",
	"effect unlock_council_slots = 1": "星系典范",
	"add_specimen random":             "寰宇博物",
	"add_relic r_head_of_zarqlan":     "远古遗物",
	"add_relic r_cetanas_heart":       "机械纪元",
	"spawn_astral_rift":               "星界位面",
	"finish_rift_stage":               "星界位面",
	"set_completed_rifts 2":           "星界位面",
	"branchoffice":                    "寰宇企业",
}

var dlcByEvent = map[string]string{
	"crisis.8005":                    "机械纪元",
	"crisis.7516":                    "机械纪元",
	"graygoo.1":                      "遥远繁星",
	"graygoo.100":                    "遥远繁星",
	"graygoo.400":                    "遥远繁星",
	"distar.11000":                   "遥远繁星",
	"distar.13000":                   "遥远繁星",
	"akx.8888":                       "地平线信号",
	"ancrel.4000":                    "远古遗物",
	"ancrel.4058":                    "远古遗物",
	"ancrel.6130":                    "远古遗物",
	"ancrel.10050":                   "远古遗物",
	"cstorms.195":                    "星体风暴",
	"grand_archive.1002":             "寰宇博物",
	"grand_archive.2300":             "寰宇博物",
	"shroud.4135":                    "暗影之幕",
	"infernals.1055":                 "炼狱",
	"nomads.610":                     "游牧",
	"fallen_empires_awakening.10":    "暗影之幕",
	"cyber.1505":                     "机械纪元",
	"synth.1505":                     "机械纪元",
	"bio.170":                        "生物创生",
	"bio.175":                        "生物创生",
	"bio.180":                        "生物创生",
	"shroud.20":                      "暗影之幕",
}

var dlcByShip = map[string]string{
	"NAME_Synth_Queen_Titan":       "机械纪元",
	"NAME_Synth_Queen_Royal_Guard": "机械纪元",
	"NAME_Synth_Queen_Harrier":     "机械纪元",
	"NAME_Paradox_Titan":           "机械纪元",
	"NAME_Behemoth_01":             "万物竞发",
	"NAME_Behemoth_02":             "万物竞发",
	"NAME_Behemoth_03":             "万物竞发",
	"NAME_Behemoth_04":             "万物竞发",
	"NAME_Venomous_Deity":          "骑士与领主",
	"NAME_Crisis_Star_Eater":       "天灾",
}

var dlcByMega = map[string]string{
	"lgate_base":                       "遥远繁星",
	"strategic_coordination_center_3":  "寰宇企业",
	"strategic_coordination_center_0":  "寰宇企业",
	"mega_shipyard_3":                  "纵横捭阖",
	"mega_shipyard_0":                  "纵横捭阖",
	"hyper_relay":                      "纵横捭阖",
	"hyper_relay_ruined":               "纵横捭阖",
	"quantum_catapult_3":               "纵横捭阖",
	"quantum_catapult_0":               "纵横捭阖",
	"quantum_catapult_ruined":          "纵横捭阖",
	"orbital_ring":                     "纵横捭阖",
	"orbital_ring_ruined":              "纵横捭阖",
	"dyson_swarm_3":                    "机械纪元",
	"dyson_swarm_1":                    "机械纪元",
	"orbital_arc_furnace_4":            "机械纪元",
	"orbital_arc_furnace_1":            "机械纪元",
	"cosmogenesis_world_0":             "机械纪元",
	"cosmogenesis_needle":              "机械纪元",
	"shroud_seal":                      "机械纪元",
	"grand_archive_0":                  "寰宇博物",
	"deep_space_citadel_0":             "寰宇博物",
	"behemoth_egg":                     "万物竞发",
	"behemoth_chrysalis":               "万物竞发",
	"galactic_crucible_4":              "炼狱",
	"galactic_crucible_0":              "炼狱",
	"dyson_gun_0":                      "游牧",
	"dyson_gun_3":                      "游牧",
	"dyson_gun_0_restored":             "游牧",
	"civilian_arkship_megastructure":   "游牧",
	"military_arkship_megastructure":   "游牧",
	"science_arkship_megastructure":    "游牧",
	"waystation_megastructure":         "游牧",
}

func eventID(command string) (string, bool) {
	if !strings.HasPrefix(command, "event ") {
		return "", false
	}
	f := strings.Fields(command)
	if len(f) < 2 {
		return "", false
	}
	return f[1], true
}

func shipDesign(command string) (string, bool) {
	if !strings.HasPrefix(command, "add_ship ") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(command, "add_ship "))
	if name == "" || strings.ContainsAny(name, " \t") {
		return "", false
	}
	return name, true
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
