package catalog

func techs() []Entry {
	return []Entry{
		e("研究全部科技", "All techs", "research_all_technologies", CatTech, "3.0.0", "不含循环科技"),
		e("研究全部含天灾", "All techs crisis", "research_all_technologies 1", CatTech, "3.0.0", ""),
		e("刷新科技选项", "Reroll techs", "techupdate", CatTech, "3.0.0", ""),
		e("完成当前研究", "Finish research", "finish_research", CatTech, "3.0.0", ""),
		e("解锁全部传统", "All traditions", "activate_all_traditions", CatTech, "3.0.0", ""),
		e("解锁泰坦科技", "Tech Titans", "research_technology tech_titans", CatTech, "3.0.0", ""),
		e("解锁巨像科技", "Tech Colossus", "research_technology tech_colossus", CatTech, "3.0.0", ""),
		e("解锁主宰科技", "Tech Juggernaut", "research_technology tech_juggernaut", CatTech, "3.0.0", ""),
		e("解锁巨型工程", "Mega-Engineering", "research_technology tech_mega_engineering", CatTech, "3.0.0", ""),
		e("解锁戴森球科技", "Tech Dyson", "research_technology tech_dyson_sphere", CatTech, "3.0.0", ""),
		e("解锁环形世界科技", "Tech Ring World", "research_technology tech_ring_world", CatTech, "3.0.0", ""),
		e("解锁星门科技", "Tech Gateway", "research_technology tech_gateway_construction", CatTech, "3.0.0", ""),
		e("解锁居住站科技", "Tech Habitat", "research_technology tech_habitat", CatTech, "3.0.0", ""),
		e("星河奇迹 飞升", "Galactic Wonders", "activate_ascension_perk ap_galactic_wonders", CatTech, "3.0.0", "ID 因 DLC 组合可能带后缀"),
		e("星河奇迹 Utopia+Megacorp", "Galactic Wonders UM", "activate_ascension_perk ap_galactic_wonders_utopia_and_megacorp", CatTech, "3.0.0", ""),
		e("化身天灾 飞升", "Become the Crisis", "activate_ascension_perk ap_become_the_crisis", CatTech, "3.0.0", ""),
		e("水之心 飞升", "Hydrocentric", "activate_ascension_perk ap_hydrocentric", CatTech, "3.2.0", ""),
		e("灵能飞升", "Mind over Matter", "activate_ascension_perk ap_mind_over_matter", CatTech, "3.0.0", ""),
		e("超凡入圣", "Transcendence", "activate_ascension_perk ap_transcendence", CatTech, "3.0.0", ""),
		e("合成进化", "Synthetic Evolution", "activate_ascension_perk ap_synthetic_evolution", CatTech, "3.0.0", ""),
		e("机械飞升", "The Flesh is Weak", "activate_ascension_perk ap_the_flesh_is_weak", CatTech, "3.0.0", ""),
		e("发动机建筑师", "Master Builders", "activate_ascension_perk ap_master_builders", CatTech, "3.0.0", ""),
		e("宇宙创生", "Cosmogenesis", "activate_ascension_perk ap_cosmogenesis", CatTech, "3.12.0", "机械纪元"),
		e("纳米飞升", "Nanotech", "activate_ascension_perk ap_nanotech", CatTech, "3.12.0", ""),
		e("模块化", "Modularity", "activate_ascension_perk ap_modularity", CatTech, "3.12.0", ""),
		e("巨兽之怒", "Behemoth Fury", "activate_ascension_perk ap_behemoth_fury", CatTech, "4.0.0", ""),
		e("星系热症", "Galactic Hyperthermia", "activate_ascension_perk ap_galactic_hyperthermia", CatTech, "4.1.0", ""),
		e("领袖技能 +10", "Skills +10", "skills 10", CatTech, "3.0.0", ""),
		e("解锁全部法令", "Unlock edicts", "unlock_edicts", CatTech, "3.0.0", "开关"),
		e("免费政策", "Free policies", "free_policies", CatTech, "3.0.0", "开关"),
		e("免费政府改革", "Free government", "free_government", CatTech, "3.0.0", "开关"),
	}
}

func relics() []Entry {
	return []Entry{
		e("全部遗珍", "All relics", "add_relic all", CatRelic, "3.0.0", ""),
		e("卡拉克水晶", "Khar-Toba", "add_relic r_unbidden_warlock", CatRelic, "3.0.0", "无形者"),
		e("赛博勒克斯战争锻造", "Cybrex War Forge", "add_relic r_war_forge", CatRelic, "3.0.0", ""),
		e("普雷索林蜂后", "Prethoryn Brood-Queen", "add_relic r_prethoryn_queen", CatRelic, "3.0.0", ""),
		e(" ethereal 灵能回响", "Psionic Archive", "add_relic r_zro_crystal", CatRelic, "3.0.0", ""),
		e("银河之头", "Head of Zarqlan", "add_relic r_head_of_zarqlan", CatRelic, "3.0.0", ""),
		e("加拉特隆", "The Galatron", "add_relic r_galatron", CatRelic, "3.0.0", ""),
		e("孤立肃正核心", "Isolated Contingency Core", "add_relic r_isolated_contingency_core", CatRelic, "3.0.0", ""),
		e("塞塔娜之心", "Cetana Heart", "add_relic r_cetanas_heart", CatRelic, "3.12.0", "合成女王"),
		e("启动遗珍 无形者", "Activate Unbidden Relic", "activate_relic r_unbidden_warlock", CatRelic, "3.0.0", ""),
		e("随机标本", "Random specimen", "add_specimen _random_", CatRelic, "3.12.0", "博物天枢"),
	}
}
