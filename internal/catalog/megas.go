package catalog

import (
	"fmt"
	"strings"
)

func mega(zh, en, id, since, note string) Entry {
	return e(zh, en, MegaSpawnCommand(id), CatMega, since, note)
}

const (
	megaStar   = "star"
	megaPlanet = "planet"
	megaRim    = "rim"
)

// MegaSpawnCommand builds a console effect that places the megastructure where
// the AI / normal construction would: star-centered types wrap the star,
// habitats and rings need a selected planet, gateways and similar sit on the
// system rim instead of stacking on the star.
func MegaSpawnCommand(id string) string {
	kind, dist, _ := megaPlacement(id)
	return MegaSpawnCommandKind(id, kind, dist)
}

func MegaSpawnCommandKind(id, kind string, dist int) string {
	switch kind {
	case megaPlanet:
		return fmt.Sprintf("effect spawn_megastructure = { type = %s planet = this }", id)
	case megaStar:
		return fmt.Sprintf("effect solar_system = { spawn_megastructure = { type = %s planet = star } }", id)
	default:
		if dist <= 0 {
			dist = 40
		}
		return fmt.Sprintf("effect solar_system = { spawn_megastructure = { type = %s planet = star orbit_distance = %d orbit_angle = random } }", id, dist)
	}
}

func MegaScope(id string) string {
	return MegaScopeKind(megaKindOf(id))
}

func MegaScopeKind(kind string) string {
	if kind == megaPlanet {
		return "planet"
	}
	return "system"
}

func megaTargetHint(command string) string {
	switch {
	case strings.Contains(command, "planet = this"):
		return "先在游戏中选中玩家行星。该巨构会生成在该行星上。"
	case strings.Contains(command, "orbit_distance"):
		return "先选中玩家星系中的恒星、星球或舰队。该巨构会放在星系外围，轨道角随机，避免叠在一起。"
	case strings.Contains(command, "planet = star"):
		return "先选中玩家星系中的恒星、星球或舰队。该巨构会包裹当前星系的恒星。"
	default:
		return "先选中玩家星系中的对象，再发送。"
	}
}

func MegaPlacementLore(kind string) string {
	switch kind {
	case megaPlanet:
		return "需先选中行星。"
	case megaStar:
		return "生成时包裹当前星系的恒星。"
	default:
		return "生成在星系外围，轨道角随机。"
	}
}

func megaKindOf(id string) string {
	k, _, _ := megaPlacement(id)
	return k
}

// ClassifyMega reads Stellaris megastructure script flags so scanned vanilla
// and MOD entries use the same placement as normal construction.
func ClassifyMega(id, script string) (kind string, dist int) {
	kind, dist, _ = megaPlacement(id)
	if script == "" {
		return
	}
	low := strings.ToLower(script)
	if d := parseOrbitDistance(low); d > 0 && kind == megaRim {
		dist = d
	}
	if hasFlag(low, "build_outside_gravity_well") {
		if d := parseOrbitDistance(low); d > 0 {
			dist = d
		}
		if dist <= 0 {
			dist = 40
		}
		return megaRim, dist
	}
	if hasFlag(low, "is_orbital_ring") || hasFlag(low, "place_entity_on_planet") ||
		hasFlag(low, "build_on_planet") || hasFlag(low, "is_habitat") ||
		strings.Contains(low, "habitat_level") {
		return megaPlanet, 0
	}
	if strings.Contains(low, "bypass_type") || hasFlag(low, "is_gateway") {
		if dist <= 0 {
			dist = 52
		}
		return megaRim, dist
	}
	starYes := strings.Contains(low, "is_star = yes") || strings.Contains(low, "is_star=yes")
	starNo := strings.Contains(low, "is_star = no") || strings.Contains(low, "is_star=no")
	negated := hasToken(low, "nor") || hasToken(low, "nand") || strings.Contains(low, "not =") || strings.Contains(low, "not=")
	if kind == megaRim && starYes && !starNo && !negated {
		return megaStar, 0
	}
	if kind == megaStar && starNo && !starYes {
		return megaPlanet, 0
	}
	return kind, dist
}

func parseOrbitDistance(script string) int {
	const key = "orbit_distance"
	i := strings.Index(script, key)
	if i < 0 {
		return 0
	}
	rest := strings.TrimSpace(script[i+len(key):])
	rest = strings.TrimPrefix(rest, "=")
	rest = strings.TrimSpace(rest)
	n := 0
	for _, r := range rest {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
		if n > 400 {
			return 400
		}
	}
	return n
}

func hasFlag(blob, name string) bool {
	return strings.Contains(blob, name+" = yes") || strings.Contains(blob, name+"=yes")
}

func hasToken(blob, token string) bool {
	i := 0
	for {
		j := strings.Index(blob[i:], token)
		if j < 0 {
			return false
		}
		j += i
		beforeOK := j == 0 || !isIdentByte(blob[j-1])
		after := j + len(token)
		afterOK := after >= len(blob) || !isIdentByte(blob[after])
		if beforeOK && afterOK {
			return true
		}
		i = j + len(token)
	}
}

func isIdentByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_'
}

func megaType(command string) (string, bool) {
	if strings.HasPrefix(command, "create_megastructure ") {
		id := strings.TrimSpace(strings.TrimPrefix(command, "create_megastructure "))
		if id != "" && !strings.ContainsAny(id, " \t{") {
			return id, true
		}
	}
	const key = "type ="
	i := strings.Index(command, key)
	if i < 0 {
		i = strings.Index(command, "type=")
		if i < 0 {
			return "", false
		}
		keyLen := len("type=")
		rest := strings.TrimSpace(command[i+keyLen:])
		return cutMegaID(rest)
	}
	rest := strings.TrimSpace(command[i+len(key):])
	return cutMegaID(rest)
}

func cutMegaID(rest string) (string, bool) {
	rest = strings.Trim(rest, `"`)
	id := rest
	for j, r := range rest {
		if r == ' ' || r == '}' || r == '\t' {
			id = rest[:j]
			break
		}
	}
	id = strings.Trim(id, `"`)
	if id == "" {
		return "", false
	}
	return id, true
}

func megaPlacement(id string) (kind string, dist, angle int) {
	low := strings.ToLower(id)
	switch {
	case hasAny(low, "habitat", "orbital_ring", "arc_furnace", "cosmogenesis_world",
		"synaptic", "behemoth_egg", "behemoth_chrysalis"):
		return megaPlanet, 0, 0
	case hasAny(low, "dyson_sphere", "dyson_swarm", "matter_decompressor", "crisis_sphere",
		"quantum_catapult", "dyson_gun", "galactic_crucible", "ring_world",
		"star_lifter", "birch", "penrose", "nidavellir"):
		return megaStar, 0, 0
	case hasAny(low, "gateway", "lgate"):
		return megaRim, 52, 0
	case hasAny(low, "hyper_relay", "waystation"):
		return megaRim, 48, 0
	case hasAny(low, "deep_space_citadel", "arkship", "shroud_seal"):
		return megaRim, 56, 0
	case hasAny(low, "think_tank", "spy_orb", "mega_art", "interstellar_assembly",
		"mega_shipyard", "strategic_coordination", "grand_archive", "cosmogenesis_needle"):
		return megaRim, 42, 0
	default:
		return megaRim, 40, 0
	}
}

func hasAny(s string, parts ...string) bool {
	for _, p := range parts {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

func megas() []Entry {
	return []Entry{
		mega("戴森球 完成", "Dyson Sphere", "dyson_sphere_5", "3.0.0", "当前星系"),
		mega("戴森球 基址", "Dyson Site", "dyson_sphere_0", "3.0.0", ""),
		mega("戴森球 框架", "Dyson Frame", "dyson_sphere_1", "3.0.0", ""),
		mega("戴森球 25%", "Dyson 25%", "dyson_sphere_2", "3.0.0", ""),
		mega("戴森球 50%", "Dyson 50%", "dyson_sphere_3", "3.0.0", ""),
		mega("戴森球 75%", "Dyson 75%", "dyson_sphere_4", "3.0.0", ""),
		mega("损坏戴森球", "Ruined Dyson", "dyson_sphere_ruined", "3.0.0", "可修复"),
		mega("环形世界 基址", "Ring World Site", "ring_world_1", "3.0.0", ""),
		mega("环形世界 框架", "Ring World Frame", "ring_world_2_intermediate", "3.0.0", ""),
		mega("环形世界 完成区段", "Ring World Section", "ring_world_3_intermediate", "3.0.0", "重复执行生成下一区段"),
		mega("损坏环形世界", "Ruined Ring", "ring_world_ruined", "3.0.0", ""),
		mega("科学枢纽 完成", "Science Nexus", "think_tank_3", "3.0.0", ""),
		mega("科学枢纽 基址", "Science Nexus Site", "think_tank_0", "3.0.0", ""),
		mega("科学枢纽 核心", "Science Nexus Hub", "think_tank_1", "3.0.0", ""),
		mega("科学枢纽 侧翼", "Science Nexus Wings", "think_tank_4", "3.0.0", "Groik 神经库版本"),
		mega("损坏科学枢纽", "Ruined Nexus", "think_tank_ruined", "3.0.0", ""),
		mega("哨兵阵列 完成", "Sentry Array", "spy_orb_4", "3.0.0", ""),
		mega("哨兵阵列 基址", "Sentry Site", "spy_orb_0", "3.0.0", ""),
		mega("损坏哨兵阵列", "Ruined Sentry", "spy_orb_ruined", "3.0.0", ""),
		mega("星门 完成", "Gateway", "gateway_final", "3.0.0", ""),
		mega("星门 基址", "Gateway Site", "gateway_0", "3.0.0", ""),
		mega("损坏星门", "Ruined Gateway", "gateway_ruined", "3.0.0", ""),
		mega("L星门", "L-Gate", "lgate_base", "3.0.0", "遥远繁星"),
		mega("轨道居住站", "Habitat", "habitat_0", "3.0.0", "3.9 前旧居住站"),
		mega("居住站核心", "Habitat Central", "habitat_central_complex", "3.9.0", "居住站重做"),
		mega("大型轨道", "Major Orbital", "habitat_major_orbital", "3.9.0", ""),
		mega("小型轨道", "Minor Orbital", "habitat_minor_orbital", "3.9.0", ""),
		mega("物质解压器 完成", "Matter Decompressor", "matter_decompressor_4", "3.0.0", "黑洞星系"),
		mega("物质解压器 基址", "Matter Decompressor Site", "matter_decompressor_0", "3.0.0", ""),
		mega("损坏物质解压器", "Ruined Decompressor", "matter_decompressor_ruined", "3.0.0", ""),
		mega("战略指挥中心 完成", "Strategic Coordination", "strategic_coordination_center_3", "3.0.0", "寰宇企业"),
		mega("战略指挥中心 基址", "SCC Site", "strategic_coordination_center_0", "3.0.0", ""),
		mega("巨型艺术设施 完成", "Mega Art", "mega_art_installation_4", "3.0.0", ""),
		mega("巨型艺术设施 基址", "Mega Art Site", "mega_art_installation_0", "3.0.0", ""),
		mega("星际集会 完成", "Interstellar Assembly", "interstellar_assembly_4", "3.0.0", ""),
		mega("星际集会 基址", "Assembly Site", "interstellar_assembly_0", "3.0.0", ""),
		mega("巨型船坞 完成", "Mega Shipyard", "mega_shipyard_3", "3.0.0", "纵横捭阖"),
		mega("巨型船坞 基址", "Mega Shipyard Site", "mega_shipyard_0", "3.0.0", ""),
		mega("超空间中继器", "Hyper Relay", "hyper_relay", "3.4.0", "纵横捭阖"),
		mega("损坏超空间中继器", "Ruined Relay", "hyper_relay_ruined", "3.4.0", ""),
		mega("量子弹弓 完成", "Quantum Catapult", "quantum_catapult_3", "3.4.0", "中子星/脉冲星"),
		mega("量子弹弓 基址", "Catapult Site", "quantum_catapult_0", "3.4.0", ""),
		mega("损坏量子弹弓", "Ruined Catapult", "quantum_catapult_ruined", "3.4.0", ""),
		mega("星环基地", "Orbital Ring", "orbital_ring", "3.4.0", "需选中行星"),
		mega("损坏星环基地", "Ruined Orbital Ring", "orbital_ring_ruined", "3.4.0", "需选中行星，可能无法修复"),
		mega("以太相引擎 完成", "Aetherophasic Engine", "crisis_sphere_4", "3.0.0", "化身天灾终局"),
		mega("以太相引擎 框架", "Engine Frame", "crisis_sphere_0", "3.0.0", ""),
		mega("以太相引擎 阶段3", "Engine Stage 3", "crisis_sphere_3", "3.0.0", ""),
		mega("损坏以太相引擎", "Ruined Engine", "crisis_sphere_ruined", "3.0.0", ""),
		mega("戴森云 完成", "Dyson Swarm", "dyson_swarm_3", "3.12.0", "机械纪元"),
		mega("戴森云 阵列", "Dyson Swarm Array", "dyson_swarm_1", "3.12.0", ""),
		mega("电弧熔炉 完成", "Arc Furnace", "orbital_arc_furnace_4", "3.12.0", "熔岩行星"),
		mega("电弧熔炉 基址", "Arc Furnace Band", "orbital_arc_furnace_1", "3.12.0", ""),
		mega("突触凝练机", "Synaptic Lathe", "cosmogenesis_world_0", "3.12.0", "宇宙创生"),
		mega("视界刺针巨构", "Horizon Needle Mega", "cosmogenesis_needle", "3.12.0", ""),
		mega("虚境封印", "Shroud Seal", "shroud_seal", "3.12.0", ""),
		mega("博物天枢", "Grand Archive", "grand_archive_0", "3.14.0", "博物天枢 DLC"),
		mega("深空擎天堡 基址", "Deep Space Citadel", "deep_space_citadel_0", "3.14.0", ""),
		mega("噬天巨兽卵", "Behemoth Egg", "behemoth_egg", "4.0.0", "万物竞发"),
		mega("巨兽蛹", "Behemoth Chrysalis", "behemoth_chrysalis", "4.0.0", ""),
		mega("星系熔炉 完成", "Galactic Crucible", "galactic_crucible_4", "4.1.0", "炼狱物种包"),
		mega("星系熔炉 框架", "Crucible Frame", "galactic_crucible_0", "4.1.0", ""),
		mega("星际大炮 基址", "Stellar Cannon", "dyson_gun_0", "4.4.0", "游牧 DLC 中子星"),
		mega("星际大炮 完成", "Stellar Cannon 3", "dyson_gun_3", "4.4.0", ""),
		mega("修复星际大炮", "Restored Cannon", "dyson_gun_0_restored", "4.4.0", ""),
		mega("民用方舟", "Civilian Arkship", "civilian_arkship_megastructure", "4.4.0", "游牧"),
		mega("军用方舟", "Military Arkship", "military_arkship_megastructure", "4.4.0", ""),
		mega("科研方舟", "Science Arkship", "science_arkship_megastructure", "4.4.0", ""),
		mega("驿站", "Waystation", "waystation_megastructure", "4.4.0", ""),
		e("移除选中巨构", "Remove megastructure", "effect remove_megastructure = this", CatMega, "3.0.0", "先选中巨构"),
	}
}
