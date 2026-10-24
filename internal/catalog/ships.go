package catalog

func ship(zh, en, name, since, note string) Entry {
	return e(zh, en, "add_ship "+name, CatShip, since, note)
}

func ships() []Entry {
	list := []Entry{
		ship("建造舰", "Constructor", "NAME_Constructor", "3.0.0", "工程船"),
		ship("科研船", "Science", "NAME_Science", "3.0.0", ""),
		ship("殖民船", "Colony Ship", "NAME_Colony_Ship", "3.0.0", ""),
		ship("运输船", "Transport", "NAME_Transport", "3.0.0", ""),
		ship("护卫舰", "Corvette", "NAME_Corvette", "3.0.0", "玩家最新设计名因帝国而异，此为全局设计"),
		ship("驱逐舰", "Destroyer", "NAME_Destroyer", "3.0.0", ""),
		ship("巡洋舰", "Cruiser", "NAME_Cruiser", "3.0.0", ""),
		ship("战列舰", "Battleship", "NAME_Battleship", "3.0.0", ""),
		ship("泰坦", "Titan", "NAME_Titan", "3.0.0", ""),
		ship("巨像", "Colossus", "NAME_Colossus", "3.0.0", "天启 DLC"),
		ship("主宰", "Juggernaut", "NAME_Juggernaut", "3.0.0", "刷出会无法造船坞功能"),
		ship("护卫舰(玩家设计)", "Player Corvette", "Corvette", "3.0.0", "使用当前帝国设计名，失败则改用 NAME_ 前缀"),
		ship("驱逐舰(玩家设计)", "Player Destroyer", "Destroyer", "3.0.0", ""),
		ship("巡洋舰(玩家设计)", "Player Cruiser", "Cruiser", "3.0.0", ""),
		ship("战列舰(玩家设计)", "Player Battleship", "Battleship", "3.0.0", ""),
		ship("泰坦(玩家设计)", "Player Titan", "Titan", "3.0.0", ""),
		e("满编制海军", "Create Navy 100%", "create_navy 1", CatShip, "3.0.0", "按最新设计填满海军上限"),
		e("半编制海军", "Create Navy 50%", "create_navy 0.5", CatShip, "3.0.0", ""),
		ship("失落帝国护卫", "FE Escort", "NAME_Scholar", "3.0.0", "知识管理者护卫"),
		ship("失落帝国战列巡洋舰", "FE Battlecruiser", "NAME_Custodian", "3.0.0", ""),
		ship("失落帝国泰坦", "FE Titan", "NAME_Reclaimer", "3.0.0", "拾荒者机器人同名需注意"),
		ship("先祖荣耀级", "NAME_Glory", "NAME_Glory", "3.0.0", "传送门事件舰"),
		ship("原型科研船", "Prototype", "NAME_Prototype", "3.0.0", ""),
		ship("外星科研船", "Alien Ship", "NAME_Alien_Ship", "3.0.0", ""),
		ship("幽灵船", "Hillos", "NAME_Hillos_ship", "3.0.0", ""),
		ship("神秘容器", "Enigmatic Cache", "NAME_enigmatic_cache", "3.0.0", ""),
		ship("纳米机器人工厂", "Nanite Factory", "NAME_Nanite_Factory", "3.0.0", "灰蛊"),
		ship("纳米机器人母舰", "Nanite Mothership", "NAME_Nanite_Mothership", "3.0.0", ""),
		ship("灰蛊战列舰", "Gray Warship", "NAME_Gray_Warship", "3.0.0", ""),
		ship("灰蛊航空母舰", "Dessanu Mothership", "NAME_Dessanu_Mothership", "3.0.0", ""),
		ship("灰蛊护卫(能量)", "Nanite Interdictor", "NAME_Nanite_Interdictor", "3.0.0", ""),
		ship("灰蛊护卫(动能)", "Dessanu Interdictor", "NAME_Dessanu_Interdictor", "3.0.0", ""),
		ship("灵能化身", "Shroud Avatar", "NAME_Shroud_Avatar", "3.0.0", ""),
		ship("堕恶化身", "Corrupted Avatar", "NAME_Corrupted_Avatar", "3.0.0", ""),
		ship("扭曲意识", "Warped Consciousness", "NAME_Warped_Consciousness", "3.0.0", ""),
		ship("虚境显圣", "Shroud Manifestation", "NAME_Shroud_Manifestation", "3.0.0", ""),
		ship("无限神机", "Infinity Machine", "NAME_Infinity_Machine", "3.0.0", "直接刷出可能是彩箱，需切国家 own"),
		ship("先驱者级泰坦", "Grand Herald", "NAME_Herald", "3.0.0", "伟大先驱考古"),
		ship("女王泰坦", "Synth Queen Titan", "NAME_Synth_Queen_Titan", "3.12.0", "机械纪元 塞塔娜"),
		ship("女王哨兵", "Synth Queen Harrier", "NAME_Synth_Queen_Harrier", "3.12.0", ""),
		ship("女王捍卫者", "Synth Queen Royal Guard", "NAME_Synth_Queen_Royal_Guard", "3.12.0", ""),
		ship("焚天神兵", "Star Eater", "NAME_Crisis_Star_Eater", "3.0.0", "化身天灾"),
		ship("视界刺针", "Horizon Needle", "NAME_Horizon_Needle", "3.12.0", "宇宙创生"),
		ship("圣骑士舰", "Paladin Ship", "NAME_Paladin", "4.4.0", "游牧/银河守护者"),
		e("瞬间建造开关", "Instant Build", "instant_build", CatShip, "3.0.0", "对 AI 同样生效，用完关闭"),
		e("瞬间移动开关", "Instant Move", "tweakergui instant_move", CatShip, "3.0.0", "舰船右键即瞬移"),
		e("无敌开关", "Invincible", "invincible", CatShip, "3.0.0", "仅舰队"),
		e("舰队受伤 1000", "Damage 1000", "damage 1000", CatShip, "3.0.0", "先选中舰队"),
		e("摧毁选中舰队伤害", "Damage huge", "damage 100000000", CatShip, "3.0.0", ""),
		e("获得选中单位所有权", "Own selected", "own", CatShip, "3.0.0", "先左键选中舰队/星球"),
	}
	// DLC-locked designs still exist in the files but cannot be built or
	// spawned without the matching expansion.
	dlc := map[string]string{
		"add_ship NAME_Colossus":                "天启",
		"add_ship NAME_Titan":                   "天启",
		"add_ship NAME_Juggernaut":              "联邦",
		"add_ship NAME_Crisis_Star_Eater":       "复仇女神",
		"add_ship NAME_Synth_Queen_Titan":       "机械纪元",
		"add_ship NAME_Synth_Queen_Harrier":     "机械纪元",
		"add_ship NAME_Synth_Queen_Royal_Guard": "机械纪元",
		"add_ship NAME_Horizon_Needle":          "机械纪元",
		"add_ship NAME_Paladin":                 "游牧",
	}
	for i := range list {
		if d, ok := dlc[list[i].Command]; ok {
			list[i].DLC = d
		}
	}
	return list
}
