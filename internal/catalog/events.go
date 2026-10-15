package catalog

func ev(zh, en, id, since, note string) Entry {
	return e(zh, en, "event "+id, CatEvent, since, note)
}

func events() []Entry {
	return []Entry{
		ev("普雷索林虫群降临", "Prethoryn Crisis", "crisis.10", "3.0.0", "天灾起始"),
		ev("虫群女王事件", "Prethoryn Queen", "crisis.100", "3.0.0", ""),
		ev("获得虫群舰队", "Prethoryn Fleet Gift", "crisis.250", "3.0.0", "选中己方舰队后更稳"),
		ev("异次元入侵者", "Unbidden", "crisis.1000", "3.0.0", ""),
		ev("肃正协议", "Contingency", "crisis.2000", "3.0.0", "部分资料写 crisis.3000，以游戏内 TAB 为准"),
		ev("肃正核心相关", "Contingency Alt", "crisis.3000", "3.0.0", "备用天灾 ID"),
		ev("化身天灾起始", "Become the Crisis", "crisis.5010", "3.0.0", "需已点飞升 perk"),
		ev("合成女王起始", "Synth Queen Start", "crisis.8005", "3.12.0", "机械纪元 塞塔娜"),
		ev("合成女王可攻击", "Queen Attackable", "crisis.9010", "3.12.0", "先选中塞塔娜局势"),
		ev("获得女王泰坦剧情", "Cetana Titan Event", "crisis.24250", "3.12.0", "粘土艺术家改写"),
		ev("合成女王胜利", "Queen Wins", "crisis.8505", "3.12.0", ""),
		ev("获得赛博勒克斯舰队", "Cybrex Fleet", "crisis.2473", "3.0.0", "可能无法解散"),
		ev("灰蛊爆发", "Gray Tempest", "graygoo.1", "3.0.0", "L星门相关"),
		ev("小灰事件", "Gray", "graygoo.401", "3.0.0", ""),
		ev("掠夺者可汗", "Great Khan", "marauder.500", "3.0.0", ""),
		ev("堕落帝国觉醒", "FE Awaken", "fallen_empire.1", "3.0.0", "ID 因版本可能变化"),
		ev("银河之星", "Galactic Star", "cara.320", "3.0.0", ""),
		ev("失落帝国舰队奖励", "FE Fleet Gift", "galactic_features.301", "3.0.0", "选中己方舰队"),
		ev("古代舰队异常", "Ancient Fleet", "anomaly.3116", "3.0.0", "选中科研船"),
		ev("F79 巡防舰", "F79 Frigate", "anomaly.2552", "3.0.0", ""),
		ev("殖民地树人链", "Tree of Life", "colony.1", "3.0.0", "选中殖民地"),
		ev("地下种族链", "Subterranean", "colony.50", "3.0.0", ""),
		ev("星球改造链", "Terraform Event", "colony.100", "3.0.0", ""),
		e("完成全部特殊项目", "Finish special projects", "finish_special_projects", CatEvent, "3.0.0", ""),
		e("完成首次接触", "Complete first contact", "complete_first_contact", CatEvent, "3.0.0", "选中接触"),
		e("完成全部帝国焦点", "Complete focus tasks", "complete_focus_tasks", CatEvent, "3.8.0", ""),
		e("局势进度 +100", "Situation +100", "effect add_situation_progress = 100", CatEvent, "3.0.0", "先选中局势"),
		e("结束选中局势", "Destroy situation", "effect destroy_situation = this", CatEvent, "3.0.0", ""),
	}
}
