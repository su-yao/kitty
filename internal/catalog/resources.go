package catalog

func resources() []Entry {
	add := func(zh, en, cmd, note string) Entry {
		return e(zh, en, cmd, CatResource, "3.0.0", note)
	}
	items := []Entry{
		add("能量币 +5000", "Energy +5000", "cash 5000", "默认数量"),
		add("能量币 +50000", "Energy +50000", "cash 50000", ""),
		add("矿物 +5000", "Minerals +5000", "minerals 5000", ""),
		add("矿物 +50000", "Minerals +50000", "minerals 50000", ""),
		add("食物 +5000", "Food +5000", "food 5000", ""),
		add("食物 +50000", "Food +50000", "food 50000", ""),
		add("合金 +5000", "Alloys +5000", "alloys 5000", "2.2 后可用"),
		add("合金 +50000", "Alloys +50000", "alloys 50000", ""),
		add("消费品 +5000", "Consumer Goods +5000", "resource consumer_goods 5000", ""),
		add("消费品 +50000", "Consumer Goods +50000", "resource consumer_goods 50000", ""),
		add("影响力 +5000", "Influence +5000", "influence 5000", "未开瞬间建造时会卡上限"),
		add("凝聚力 +5000", "Unity +5000", "unity 5000", ""),
		add("凝聚力 +50000", "Unity +50000", "unity 50000", ""),
		add("物理研究 +5000", "Physics +5000", "physics 5000", ""),
		add("社会研究 +5000", "Society +5000", "society 5000", ""),
		add("工程研究 +5000", "Engineering +5000", "engineering 5000", ""),
		add("填满全部资源", "Max resources", "max_resources", "未解锁战略资源可能不加"),
		add("活金属 +5000", "Living Metal", "resource sr_living_metal 5000", ""),
		add("泽洛 +5000", "Zro", "resource sr_zro 5000", ""),
		add("暗物质 +5000", "Dark Matter", "resource sr_dark_matter 5000", ""),
		add("易爆微粒 +5000", "Volatile Motes", "resource sr_volatile_motes 5000", ""),
		add("异星天然气 +5000", "Exotic Gases", "resource sr_exotic_gases 5000", ""),
		add("稀有水晶 +5000", "Rare Crystals", "resource sr_rare_crystals 5000", ""),
		add("纳米机器人 +5000", "Nanites", "resource nanites 5000", "L星团/灰蛊"),
		add("次要文物 +500", "Minor Artifacts", "minor_artifacts 500", "远古遗物 DLC"),
		add("威慑度 +5000", "Menace", "menace 5000", "化身天灾 3.0+"),
		add("帝国权威 +50", "Imperial Authority", "imperial_authority 50", "星海帝国"),
	}
	items = append(items,
		e("星界丝线 +5000", "Astral Threads", "astral_threads 5000", CatResource, "3.10.0", "星界裂隙"),
		e("高级逻辑 +5000", "Advanced Logic", "advanced_logic 5000", CatResource, "3.12.0", "宇宙创生"),
	)
	return items
}
