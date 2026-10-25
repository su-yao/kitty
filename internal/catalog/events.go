package catalog

func ev(zh, en, id, since, note string) Entry {
	return e(zh, en, "event "+id, CatEvent, since, note)
}

func events() []Entry {
	list := []Entry{
		ev("普雷索林虫群降临", "Prethoryn Scourge", "crisis.199", "3.0.0", "开启虫群天灾"),
		ev("哨兵崛起", "Rise of the Sentinels", "crisis.50", "3.0.0", "不依赖天灾，含音效"),
		ev("哨兵舰队捐赠", "Sentinel Fleet Donation", "crisis.71", "3.0.0", "不依赖天灾或哨兵"),
		ev("虫群女王万岁", "Long live the Queen!", "crisis.105", "3.0.0", "选中舰队，生成驯化蜂后"),
		ev("赛博勒克斯回归", "Cybrex Return", "crisis.2400", "3.0.0", "不依赖天灾"),
		ev("焚天神兵开火", "Star-Eater Firing", "crisis.4550", "3.0.0", "当前星系被摧毁"),
		ev("异次元入侵者", "The Unbidden", "crisis.1000", "3.0.0", "开启无形者天灾"),
		ev("变异者来袭", "The Aberrant", "crisis.1100", "3.0.0", "不依赖无形者"),
		ev("狂暴者来袭", "The Vehement", "crisis.1200", "3.0.0", "需先有无形者"),
		ev("肃正协议", "The Contingency", "crisis.2000", "3.0.0", "开启肃正天灾"),
		ev("宇宙创生结局", "Cosmogenesis Ending", "crisis.7516", "3.12.0", "帝国化为堕落并强制换国"),
		ev("合成女王起始", "The Synthetic Queen", "crisis.8005", "3.12.0", "机械纪元 塞塔娜"),
		ev("沉睡者苏醒", "Sleepers Awake", "fallen_empires_awakening.1", "3.0.0", "需先用 play 选中堕落帝国，否则自己觉醒"),
		ev("远古看护者苏醒", "Ancient Caretakers Awaken", "fallen_machine_empire.1", "3.0.0", "需先用 play 选中远古看护者"),
		ev("天堂之战", "War in Heaven", "war_in_heaven.1000", "3.0.0", "唤醒两个堕落帝国"),
		ev("堕落帝国馈赠", "FE Patronizing Gift", "fallen_empires_tasks.1", "3.0.0", "堕落帝国随机赠礼"),
		ev("战鼓敲响", "The Drums of War", "marauder.500", "3.0.0", "需先用 play 选中掠夺者，否则自己成可汗"),
		ev("雇佣兵可用", "Mercenaries Available", "marauder.85", "3.0.0", "每个掠夺者帝国需触发多次"),
		ev("银河共同体诞生", "Birth of Galactic Community", "galcom.16", "3.0.0", "立即创建或加入共同体"),
		ev("银河市场建立", "Galactic Market", "action.99", "3.0.0", "在首都星系建立银河市场"),
		ev("生成 L 星团", "Spawn L-Cluster", "distar.11000", "3.0.0", "生成封闭的 L 星团"),
		ev("L 星团(L 龙)", "L-Cluster (L-Drake)", "distar.13000", "3.0.0", "会打开所有 L 星门"),
		ev("灰蛊爆发", "Gray Tempest", "graygoo.1", "3.0.0", "需已有 L 星团"),
		ev("L 星团(德萨努)", "Dessanu Consonance", "graygoo.100", "3.0.0", "需已有 L 星团"),
		ev("静谧漫步", "A Quiet Stroll", "graygoo.400", "3.0.0", "选中舰队，遇见小灰"),
		ev("地平线信号", "The Horizon Signal", "akx.8888", "3.0.0", "选中黑洞星系中有领袖的舰船"),
		ev("旅行者一号", "Voyager 1", "anomaly.95", "3.0.0", "需星系中存在太阳系"),
		ev("林波", "Limbo", "anomaly.186", "3.0.0", "可复活为新殖民地或帝国"),
		ev("变形虫安抚", "Amoebas Pacified", "story.107", "3.0.0", ""),
		ev("水晶安抚", "Crystals Pacified", "story.207", "3.0.0", ""),
		ev("虚境导师", "Teachers of the Shroud", "origin.5605", "3.0.0", "解锁虚境信标恒星基地建筑"),
		ev("天选之人", "The Chosen One", "utopia.3190", "3.0.0", ""),
		ev("一个问题", "A Question", "syndaw.545", "3.0.0", "开启 AI 相关事件局势"),
		ev("雇佣兵驻地", "Mercenary Garrison", "enclave.220", "3.0.0", "可选购买雇佣兵驻地建筑"),
		ev("无限方程式", "The Infinity Equation", "leviathans.1040", "3.0.0", "击败无限神机奖励"),
		ev("恩尼格玛要塞停用", "Enigmatic Fortress Disabled", "leviathans.2101", "3.0.0", "击败恩尼格玛要塞奖励"),
		ev("无畏舰停用", "Dreadnought Disabled", "leviathans.3102", "3.0.0", "击败自动无畏舰奖励"),
		ev("无畏舰修复", "Dreadnought Repaired", "leviathans.3103", "3.0.0", "选中舰船"),
		ev("拾荒者的终结", "The Scavenger's End", "distar.137", "3.0.0", "击败拾荒者机器人奖励"),
		ev("母兽之死", "Death of the Matriarch", "distar.212", "3.0.0", "选中舰队"),
		ev("野生真核生物", "Wild Eukaryotes", "distar.260", "3.0.0", "选中星球，生成前哨物种"),
		ev("失乐园", "Paradise Lost", "distar.1001", "3.0.0", "生成盖亚世界与石器时代文明"),
		ev("异星实体", "Alien Entity", "distar.2050", "3.0.0", "生成神秘容器"),
		ev("尼夫拉克(敌对)", "The Nivlac (unfriendly)", "distar.3014", "3.0.0", "生成 -100 好感帝国"),
		ev("尼夫拉克(友好)", "The Nivlac (friendly)", "distar.3016", "3.0.0", "生成 +100 好感帝国"),
		ev("异星盒子开启", "Alien Box Opened", "distar.3055", "3.0.0", "获得一个随机免费特质"),
		ev("虚空孳孽破壳", "The Voidspawn", "distar.5006", "3.0.0", "选中星球，变为破碎世界"),
		ev("巨化进化", "Gargantuan Evolution", "distar.5012", "3.0.0", "击败虚空孳孽奖励"),
		ev("石中低语", "Whispers in the Stone", "ancrel.4000", "3.0.0", "选中星球，生成考古遗址"),
		ev("扎克兰的选民", "Zarqulan's Chosen", "ancrel.6130", "3.0.0", "可殖民圣地"),
		ev("伊尼提亚商人", "Inetian Traders", "cstorms.195", "3.13.0", "解锁超空间细化器建筑"),
		ev("原始呼唤姿态", "Primal Calling Stance", "grand_archive.1002", "3.14.0", "可解锁多个姿态建筑"),
		ev("沃尔特劳姆母星", "Vultraum Home System", "precursor.98", "3.0.0", "生成先驱者母星系"),
		ev("尤特母星", "Yuht Home System", "precursor.598", "3.0.0", ""),
		ev("第一联盟母星", "First League Home System", "precursor.1098", "3.0.0", ""),
		ev("伊拉西安母星", "Irassian Home System", "precursor.1598", "3.0.0", ""),
		ev("赛博勒克斯母星", "Cybrex Home System", "precursor.2098", "3.0.0", ""),
		ev("虚境契约 低语者", "Whisperers in the Void", "utopia.3304", "3.0.0", "可选签订契约"),
		ev("虚境契约 织缕", "Composer of Strands", "utopia.3305", "3.0.0", ""),
		ev("虚境契约 噬世者", "Eater of Worlds", "utopia.3306", "3.0.0", ""),
		ev("虚境契约 欲望", "Instrument of Desire", "utopia.3307", "3.0.0", ""),
		ev("轮回终末", "End of the Cycle", "utopia.3308", "3.0.0", "可选签订轮回终末契约"),
		ev("清算之日", "The Reckoning", "utopia.3320", "3.0.0", "需先接受上一条契约"),
		ev("轮回终末(暗影)", "End of the Cycle (Shroud)", "shroud.4135", "4.3.0", "暗影之幕版本契约"),
		ev("神经共生", "Neural Symbiosis", "distar.172", "3.0.0", "选中舰船，可获得脑蛞蝓特质"),
		ev("虚空之虫瘟疫", "Voidworm Plague", "grand_archive.2300", "3.14.0", "先用 play 选中虚空之虫，否则帝国可能崩溃"),
		ev("破碎片段苏醒", "Shattered Fragments Awaken", "fallen_empires_awakening.10", "4.3.0", "需先用 play 选中破碎片段"),
		ev("星系热症结局", "Galactic Hyperthermia", "infernals.1055", "4.1.0", "全恒星变红巨星，非嗜热帝国灭亡"),
		ev("游牧天命", "Vindication", "nomads.610", "4.4.0", "可选择成为游牧部落"),
		ev("获得哨兵军队", "The Sentinels", "ancrel.4058", "3.0.0", "选中星球，可获哨兵军队"),
		ev("尤特净化秘密", "Secrets of the Yuht", "ancrel.10050", "3.0.0", "解锁尤特净化决议"),
		e("完成全部特殊项目", "Finish special projects", "finish_special_projects", CatEvent, "3.0.0", ""),
		e("完成首次接触", "Complete first contact", "complete_first_contact", CatEvent, "3.0.0", "选中接触"),
		e("添加首次接触线索", "First contact clues", "add_first_contact_clues 10", CatEvent, "3.0.0", "选中进行中的首次接触"),
		e("完成全部帝国焦点", "Complete focus tasks", "complete_focus_tasks", CatEvent, "3.8.0", ""),
		e("局势进度 +100", "Situation +100", "effect add_situation_progress = 100", CatEvent, "3.0.0", "先选中局势"),
		e("结束选中局势", "Destroy situation", "effect destroy_situation = this", CatEvent, "3.0.0", ""),
		e("摧毁选中舰队", "Destroy fleet", "effect destroy_fleet = this", CatEvent, "3.0.0", "先选中舰队"),
	}
	// Selection requirements verified against the in-game event scopes. The
	// scanner overrides these when it can read the current game files.
	type req struct {
		scope  string
		target bool
	}
	reqs := map[string]req{
		"event crisis.105":      {"fleet", false},
		"event graygoo.400":     {"fleet", false},
		"event akx.8888":        {"ship", true},
		"event distar.212":      {"fleet", false},
		"event leviathans.3103": {"ship", true},
		"event distar.172":      {"ship", true},
		"event distar.260":      {"planet", false},
		"event distar.1001":     {"planet", false},
		"event distar.5006":     {"planet", false},
		"event ancrel.4000":     {"planet", false},
		"event ancrel.4058":     {"planet", false},
		"event enclave.220":     {"planet", false},
		"event anomaly.186":     {"planet", false},
		"complete_first_contact": {"first_contact", false},
		"add_first_contact_clues 10": {"first_contact", false},
		"effect add_situation_progress = 100": {"situation", false},
		"effect destroy_situation = this": {"situation", false},
		"effect destroy_fleet = this": {"fleet", false},
	}
	for i := range list {
		if r, ok := reqs[list[i].Command]; ok {
			list[i].Scope = r.scope
			list[i].Target = r.target
		}
	}
	return list
}
