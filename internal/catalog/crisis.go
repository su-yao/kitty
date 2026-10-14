package catalog

func cship(zh, en, name, since, note string) Entry {
	return e(zh, en, "add_ship "+name, CatCrisis, since, note)
}

func crisisShips() []Entry {
	return []Entry{
		cship("虫群幼体", "Prethoryn Swarmling", "NAME_Swarmling", "3.0.0", "普雷索林"),
		cship("虫群战士", "Prethoryn Warrior", "NAME_Warrior", "3.0.0", ""),
		cship("虫群育母", "Prethoryn Brood Mother", "NAME_Brood_Mother", "3.0.0", ""),
		cship("虫群女王", "Prethoryn Queen", "NAME_Queen", "3.0.0", ""),
		cship("无形者护卫", "Unbidden Escort", "NAME_Revenant", "3.0.0", "异次元入侵者"),
		cship("无形者巡洋", "Unbidden Cruiser", "NAME_Phantom", "3.0.0", ""),
		cship("无形者战列", "Unbidden Battleship", "NAME_Wraith", "3.0.0", ""),
		cship("无形者锚点", "Unbidden Anchor", "NAME_Unbidden_Anchor", "3.0.0", ""),
		cship("肃正寻求者", "Contingency Seeker", "NAME_Sterilizer", "3.0.0", "痛苦终结者"),
		cship("肃正歼灭者", "Contingency Euthanizer", "NAME_Euthanizer", "3.0.0", ""),
		cship("肃正恒星基地", "AI Starbase", "NAME_AI_STARBASE", "3.0.0", ""),
		cship("灰蛊拦截者", "Graygoo Interdictor", "NAME_Nanite_Interdictor", "3.0.0", "L星团"),
		cship("灰蛊母舰", "Graygoo Mothership", "NAME_Nanite_Mothership", "3.0.0", ""),
		cship("焚天神兵", "Star Eater", "NAME_Crisis_Star_Eater", "3.0.0", "化身天灾"),
		cship("女王泰坦 塞塔娜", "Cetana Titan", "NAME_Synth_Queen_Titan", "3.12.0", "合成女王"),
		cship("女王护卫", "Queen Royal Guard", "NAME_Synth_Queen_Royal_Guard", "3.12.0", ""),
		cship("女王猎手", "Queen Harrier", "NAME_Synth_Queen_Harrier", "3.12.0", ""),
		cship("悖论泰坦", "Paradox Titan", "NAME_Paradox_Titan", "3.12.0", "宇宙创生危机舰"),
		cship("巨兽 I 级", "Behemoth I", "NAME_Behemoth_01", "4.0.0", "巨兽之怒"),
		cship("巨兽 II 级", "Behemoth II", "NAME_Behemoth_02", "4.0.0", ""),
		cship("巨兽 III 级", "Behemoth III", "NAME_Behemoth_03", "4.0.0", ""),
		cship("巨兽 IV 级", "Behemoth IV", "NAME_Behemoth_04", "4.0.0", ""),
	}
}

func lev(zh, en, name, since, note string) Entry {
	return e(zh, en, "add_ship "+name, CatLeviathan, since, note)
}

func leviathans() []Entry {
	return []Entry{
		lev("以太巨龙", "Ether Drake", "NAME_Grand_Dragon", "3.0.0", ""),
		lev("天龙", "Sky Dragon", "NAME_Sky_Dragon", "3.0.0", "起源龙"),
		lev("以太幼龙", "Dragonspawn", "NAME_Dragonspawn", "3.0.0", ""),
		lev("幼天龙", "Sky Dragon Baby", "NAME_Sky_Dragon_Baby", "3.0.0", ""),
		lev("L星龙", "Nanite Dragon", "NAME_Nanite_Dragon", "3.0.0", ""),
		lev("裂爪龙", "Shard Dragon", "NAME_Shard_Dragon", "3.0.0", ""),
		lev("噬星者", "Stellarite", "NAME_Stellarite", "3.0.0", ""),
		lev("位面之魇", "Dimensional Horror", "NAME_Dimensional_Horror", "3.0.0", "无法移动"),
		lev("虚空孳孽", "Voidspawn", "NAME_Voidspawn", "3.0.0", ""),
		lev("拾荒者机器人", "Scavenger Bot", "NAME_Reclaimer", "3.0.0", ""),
		lev("泡泡(百年变形虫)", "Amoeba Centenarian", "NAME_Adopted_Amoeba_Centenarian", "3.0.0", ""),
		lev("光谱幽魂 蓝", "Wraith Blue", "NAME_Spectral_Wraith_650THz", "3.0.0", ""),
		lev("光谱幽魂 红", "Wraith Red", "NAME_Spectral_Wraith_450THz", "3.0.0", ""),
		lev("光谱幽魂 黄", "Wraith Yellow", "NAME_Spectral_Wraith_520THz", "3.0.0", ""),
		lev("毒圣", "Venomous Deity", "NAME_Venomous_Deity", "3.7.0", "骑士与领主"),
		lev("古代采矿无人机", "Ancient Mining Drone", "NAME_Ancient_Mining_Drone", "3.0.0", ""),
		lev("古代战斗无人机", "Ancient Combat Drone", "NAME_Ancient_Combat_Drone", "3.0.0", ""),
		lev("古代驱逐舰", "Ancient Destroyer", "NAME_Ancient_Destroyer", "3.0.0", ""),
		lev("水晶巢穴", "Crystal Nidus", "NAME_Crystal_Nidus", "3.0.0", ""),
		lev("大蓝水晶", "Large Crystal Blue", "NAME_Large_Crystal_Entity_Blue", "3.0.0", ""),
		lev("中蓝水晶", "Medium Crystal Blue", "NAME_Medium_Crystal_Entity_Blue", "3.0.0", ""),
		lev("小蓝水晶", "Small Crystal Blue", "NAME_Small_Crystal_Entity_Blue", "3.0.0", ""),
		lev("泰扬基母兽", "Tiyanki Matriarch", "NAME_Tiyanki_Matriarch", "3.0.0", ""),
		lev("无畏舰", "Automated Dreadnought", "NAME_Dreadnought", "3.0.0", ""),
		lev("恩尼格玛守卫", "Enigmatic Fortress", "NAME_Enigmatic_Fortress", "3.0.0", "可能无法移动"),
	}
}
