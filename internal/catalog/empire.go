package catalog

func planets() []Entry {
	return []Entry{
		e("盖亚星球", "Gaia", "planet_class pc_gaia", CatPlanet, "3.0.0", "先选中星球"),
		e("大陆星球", "Continental", "planet_class pc_continental", CatPlanet, "3.0.0", ""),
		e("海洋星球", "Ocean", "planet_class pc_ocean", CatPlanet, "3.0.0", ""),
		e("热带星球", "Tropical", "planet_class pc_tropical", CatPlanet, "3.0.0", ""),
		e("沙漠星球", "Desert", "planet_class pc_desert", CatPlanet, "3.0.0", ""),
		e("干旱星球", "Arid", "planet_class pc_arid", CatPlanet, "3.0.0", ""),
		e("苔原星球", "Tundra", "planet_class pc_tundra", CatPlanet, "3.0.0", ""),
		e("极地星球", "Arctic", "planet_class pc_arctic", CatPlanet, "3.0.0", ""),
		e("高山星球", "Alpine", "planet_class pc_alpine", CatPlanet, "3.0.0", ""),
		e("墓世界", "Tomb", "planet_class pc_nuked", CatPlanet, "3.0.0", ""),
		e("生态都市", "Ecumenopolis", "planet_class pc_city", CatPlanet, "3.0.0", ""),
		e("蜂巢星球", "Hive", "planet_class pc_hive", CatPlanet, "3.0.0", ""),
		e("机械星球", "Machine", "planet_class pc_machine", CatPlanet, "3.0.0", ""),
		e("环世界区段", "Ring World", "planet_class pc_ringworld_habitable", CatPlanet, "3.0.0", ""),
		e("星球大小 25", "Planet size 25", "planet_size 25", CatPlanet, "3.0.0", "最大建议 25"),
		e("完成环境改造", "Finish terraform", "finish_terraform", CatPlanet, "3.0.0", ""),
		e("殖民选中星球", "Colonize", "colonize 0", CatPlanet, "3.0.0", "物种 ID 0"),
		e("添加机器人人口 10", "Build pops 10", "build_pops 10", CatPlanet, "3.0.0", "需机器人科技"),
		e("人口增长加速", "Grow pops 100", "grow_pops 100", CatPlanet, "3.0.0", ""),
		e("控制选中星球", "Control", "control", CatPlanet, "3.0.0", ""),
		e("拥有选中星球", "Own planet", "own", CatPlanet, "3.0.0", ""),
		e("歼星选定星球摧毁殖民地", "Destroy colony", "effect destroy_colony", CatPlanet, "3.0.0", ""),
		e("行星荒芜 +50", "Devastation +50", "effect add_planet_devastation = 50", CatPlanet, "3.0.0", ""),
		e("行星荒芜清零", "Devastation 0", "effect add_planet_devastation = -100", CatPlanet, "3.0.0", ""),
		e("添加城市区划", "City district", "effect add_district = district_city", CatPlanet, "3.0.0", ""),
		e("添加采矿区划", "Mining district", "effect add_district = district_mining", CatPlanet, "3.0.0", ""),
		e("添加发电区划", "Generator district", "effect add_district = district_generator", CatPlanet, "3.0.0", ""),
		e("添加农业区划", "Farming district", "effect add_district = district_farming", CatPlanet, "3.0.0", ""),
	}
}

func empires() []Entry {
	return []Entry{
		e("与全部帝国建交", "Communications", "communications", CatEmpire, "3.0.0", ""),
		e("接触全部帝国", "Contact", "contact", CatEmpire, "3.0.0", ""),
		e("全图视野 intel", "Intel", "intel", CatEmpire, "3.0.0", "开关"),
		e("调查全部星系", "Survey", "survey", CatEmpire, "3.0.0", ""),
		e("AI 全同意", "Yes-men", "debug_yesmen", CatEmpire, "3.0.0", "开关"),
		e("AI 全拒绝", "No-men", "debug_nomen", CatEmpire, "3.0.0", "开关"),
		e("观察者模式", "Observe", "observe", CatEmpire, "3.0.0", "用前存档"),
		e("切回玩家 0", "Play 0", "play 0", CatEmpire, "3.0.0", ""),
		e("切到帝国 1", "Play 1", "play 1", CatEmpire, "3.0.0", ""),
		e("AI 开关", "Toggle AI", "ai", CatEmpire, "3.0.0", ""),
		e("吞并选中/指定", "Annex", "annex", CatEmpire, "3.0.0", "可跟帝国 ID"),
		e("强制整合", "Force integrate", "force_integrate", CatEmpire, "3.0.0", "跟目标 ID"),
		e("结束议会休会", "Force senate vote", "force_senate_vote", CatEmpire, "3.0.0", ""),
		e("结束当前决议", "End senate session", "end_senate_session", CatEmpire, "3.0.0", ""),
		e("跳过共同体冷却", "Skip GC cooldowns", "skip_galactic_community_cooldowns", CatEmpire, "3.0.0", ""),
		e("开始选举", "Election", "election", CatEmpire, "3.0.0", ""),
		e("好感 +100 对 1", "Opinion +100", "add_opinion 0 1 100", CatEmpire, "3.0.0", "源 目标 数值"),
		e("情报 +100 对 1", "Intel +100", "add_intel 1 100", CatEmpire, "3.0.0", ""),
		e("忠诚 +100", "Loyalty +100", "add_loyalty 1 100", CatEmpire, "3.0.0", ""),
		e("间谍渗透 +1000", "Spy +1000", "add_spynetwork_value 1 1000", CatEmpire, "3.0.0", ""),
		e("激活全部星门", "Activate gateways", "activate_gateways", CatEmpire, "3.0.0", ""),
		e("切换到 2400 年银河", "Advanced galaxy", "advanced_galaxy", CatEmpire, "3.0.0", "破坏性"),
	}
}

func debugs() []Entry {
	return []Entry{
		e("调试提示开关", "Debug tooltip", "debugtooltip", CatDebug, "3.0.0", "悬停显示 ID"),
		e("帮助列表", "Help", "help", CatDebug, "3.0.0", ""),
		e("显示版本", "Version", "version", CatDebug, "3.0.0", ""),
		e("FPS 统计", "3D stats", "3dstats", CatDebug, "3.0.0", ""),
		e("成就状态", "Achievement status", "achievement_status", CatDebug, "3.0.0", ""),
		e("杀死选中国家", "Kill country", "kill_country", CatDebug, "3.0.0", "4.4 起可跟 ID"),
		e("时间快进 360 天", "Fast forward 360", "fast_forward 360", CatDebug, "3.0.0", ""),
		e("加 10 年", "Add 10 years", "add_time years 10", CatDebug, "3.0.0", ""),
		e("重载界面", "Reload interface", "reloadinterface", CatDebug, "3.0.0", ""),
		e("重载本地化", "Reload localization", "reload localization", CatDebug, "3.0.0", ""),
		e("清空控制台", "Clear", "clear", CatDebug, "3.0.0", ""),
		e("分支办公室", "Branch office", "branchoffice", CatDebug, "3.0.0", "选中星球"),
		e("所有舰队互相攻击", "Attack all fleets", "attackallfleets", CatDebug, "3.0.0", "测试用"),
	}
}
