package catalog

import "strings"

func loreFor(e Entry) string {
	if s := loreExact[e.Command]; s != "" {
		return s
	}
	if id, ok := eventID(e.Command); ok {
		if s := loreEvent[id]; s != "" {
			return s
		}
	}
	if name, ok := shipDesign(e.Command); ok {
		if s := loreShip[name]; s != "" {
			return s
		}
	}
	if strings.HasPrefix(e.Command, "create_megastructure ") {
		id := strings.TrimPrefix(e.Command, "create_megastructure ")
		if s := loreMega[id]; s != "" {
			return s
		}
	}
	return loreByCategory(e)
}

func loreByCategory(e Entry) string {
	switch e.Category {
	case CatResource:
		return "直接增加玩家帝国库存。战略资源需先解锁对应科技，未解锁时可能不加。"
	case CatShip:
		return "在当前选中的舰队、星球或恒星旁生成一艘舰船。官方 Wiki：add_ship 使用全局设计名，玩家自制设计名因帝国而异。"
	case CatCrisis:
		return "生成天灾阵营舰船。默认归属玩家，可用 own 把所有权交给对应天灾国家。"
	case CatLeviathan:
		return "生成银河守卫者或利维坦级单位。部分无法移动，击败后常解锁奖励事件或遗珍。"
	case CatMega:
		return "在当前选中星系生成巨构。部分阶段可重复执行以推进下一阶段；损坏版可修复。"
	case CatEvent:
		return "触发对应事件。国家事件一般无需选中；星球/舰队/舰船事件需先选中正确对象。舰船事件可附加目标 ID。"
	case CatTech:
		return "立即研究科技、激活传统或飞升。部分飞升 ID 会随已安装 DLC 变化。"
	case CatRelic:
		return "把遗珍加入玩家帝国收藏。首次激活部分遗珍需重进存档才显示按钮。"
	case CatPlanet:
		return "作用于当前选中星球：改类型、大小、区划、人口或殖民。过大体积会把星球推出星系。"
	case CatEmpire:
		return "作用于玩家帝国外交、联邦、共同体或战争。带数字参数时，0 通常是玩家。"
	case CatDebug:
		return "调试命令。debugtooltip 可在悬停时显示星球、舰船、帝国和人口 ID。"
	case CatMod:
		return "扫描游戏或 MOD 文件得到的条目。发送前请确认当前版本与 DLC 已安装对应内容。"
	}
	return e.Note
}

var loreExact = map[string]string{
	"cash 5000":                       "增加玩家能量币。能量币用于维护、建造和市场交易。",
	"cash 50000":                      "大量增加玩家能量币。",
	"minerals 5000":                   "增加矿物。矿物是建筑、区划和舰船船体的基础资源。",
	"minerals 50000":                  "大量增加矿物。",
	"food 5000":                       "增加食物。生物人口增长依赖食物；机械帝国通常不需要。",
	"food 50000":                      "大量增加食物。",
	"alloys 5000":                     "增加合金。合金用于舰船、恒星基地和多数巨构。2.2 利维坦经济重做后加入。",
	"alloys 50000":                    "大量增加合金。",
	"resource consumer_goods 5000":    "增加消费品。常规帝国用于人口维护和部分建筑。",
	"resource consumer_goods 50000":   "大量增加消费品。",
	"influence 5000":                  "增加影响力。用于宣称、法令、外交和恒星基地扩建。瞬间建造开启时可能卡在上限。",
	"unity 5000":                      "增加凝聚力。用于传统、飞升perk与议程。",
	"unity 50000":                     "大量增加凝聚力。",
	"physics 5000":                    "增加物理研究点数，立即推进物理科研。",
	"society 5000":                    "增加社会研究点数。",
	"engineering 5000":                "增加工程研究点数。",
	"max_resources":                   "把已解锁资源填到库存上限。未解锁的战略资源可能不加。",
	"resource sr_living_metal 5000":   "活金属。用于部分先进组件与巨构维修。",
	"resource sr_zro 5000":            "泽洛。灵能相关战略资源，虚境与部分飞地交易会用到。",
	"resource sr_dark_matter 5000":    "暗物质。用于暗物质反应器、护盾与跃迁核心。",
	"resource sr_volatile_motes 5000": "易爆微粒。战略资源，用于武器与精炼。",
	"resource sr_exotic_gases 5000":   "异星天然气。战略资源，用于护盾与高级建筑。",
	"resource sr_rare_crystals 5000":  "稀有水晶。战略资源，用于激光与奢侈品。",
	"resource nanites 5000":           "纳米机器人。L 星团与灰蛊相关，也可用于纳米飞升。",
	"minor_artifacts 500":             "次要文物。远古遗物 DLC，用于遗物决议与博物馆。",
	"menace 5000":                     "化身天灾的威慑度，用于危机升级与特殊舰船。",
	"imperial_authority 50":           "星海帝国的帝国权威，用于帝制联邦集权。",
	"astral_threads 5000":             "星界丝线。星界裂隙 DLC，用于裂隙探索与星界动作。",
	"advanced_logic 5000":             "高级逻辑。宇宙创生飞升的特殊资源，用于视界刺针。",
	"instant_build":                   "瞬间完成建造与升级。对玩家和 AI 同时生效，用完务必关闭。",
	"tweakergui instant_move":         "舰队右键即可瞬移到目标星系。测试用开关。",
	"invincible":                      "玩家舰队不受伤害。仅对舰队生效。",
	"damage 1000":                     "对当前选中舰队造成 1000 点伤害。",
	"damage 100000000":                "对当前选中舰队造成巨额伤害，通常直接摧毁。",
	"own":                             "把当前选中的舰队、星球或巨构所有权转给玩家。",
	"create_navy 1":                   "按玩家最新舰船设计把海军容量填满。",
	"create_navy 0.5":                 "按玩家最新设计生成半编制海军。",
	"research_all_technologies":       "立刻研究当前可解锁的全部常规科技，不含循环科技。",
	"research_all_technologies 1":     "含太空生物与天灾科技的全科技解锁。",
	"research_all_technologies 1 5":   "全科技外再叠 5 层循环科技。",
	"techupdate":                      "立刻刷新当前研究选项。",
	"finish_research":                 "完成当前三项研究，同时完成特殊项目。",
	"activate_all_traditions":         "解锁全部传统树。",
	"debugtooltip":                    "悬停显示星球、舰船、帝国、人口与领袖的内部 ID。舰船事件目标 ID 由此查看。",
	"play 0":                          "切回玩家帝国。玩家国家 ID 固定为 0。",
	"play 1":                          "接管帝国 1。常用于给堕落帝国或掠夺者触发觉醒事件。",
	"communications":                  "与银河中全部帝国、飞地和虚境实体建立通讯。",
	"contact":                         "接触全部已知帝国。",
	"intel":                           "开关全图视野。",
	"survey":                          "立刻调查全部星系，需至少一艘科研船。",
	"debug_yesmen":                    "AI 对玩家的所有提议一律同意。",
	"debug_nomen":                     "AI 对玩家的所有提议一律拒绝。",
	"observe":                         "进入观察者模式。操作前请存档。",
	"ai":                              "开关银河 AI。关闭后其它帝国停止行动。",
	"finish_special_projects":         "完成玩家当前全部特殊项目。",
	"complete_first_contact":          "完成当前选中的首次接触。",
	"complete_focus_tasks":            "完成当前帝国焦点任务。",
	"colonize 1":                      "立刻殖民选中星球。参数是殖民舰人口 ID，默认 1。",
	"planet_class pc_gaia":            "把选中星球改为盖亚世界，宜居度最高。",
	"planet_class pc_city":            "改为生态都市。需要天启相关科技或飞升支持才有完整区划。",
	"planet_class pc_hive":            "改为蜂巢星球，适合蜂巢思维。",
	"planet_class pc_machine":         "改为机械星球，适合机械智能。",
	"planet_size 25":                  "把选中星球大小改为 25。过大的体积会把星球推出星系。",
	"add_relic all":                   "把游戏中全部遗珍加入玩家收藏。",
	"annex":                           "吞并当前选中或指定 ID 的帝国。",
	"force_integrate":                 "强制整合指定附属国。",
	"election":                        "立刻开始一场选举。",
	"activate_gateways":               "激活银河中全部星门。",
}

var loreEvent = map[string]string{
	"crisis.199":                  "普雷索林虫群（Prethoryn Scourge）从银河边缘降临。这是原版三大终局天灾之一，蜂群会吞噬星球并繁殖。",
	"crisis.50":                   "哨兵组织崛起。哨兵是对抗虫群的跨银河联军，此事件可不依赖天灾单独触发。",
	"crisis.71":                   "哨兵向玩家捐赠一支舰队。不依赖天灾或哨兵已存在。",
	"crisis.105":                  "「虫群女王万岁」：在选中舰队处生成一艘驯化的普雷索林蜂后。需先选中玩家舰队。",
	"crisis.2400":                 "赛博勒克斯回归。赛博勒克斯是远古机械种族，常作为对抗肃正协议的盟友出现。",
	"crisis.4550":                 "焚天神兵开火。化身天灾的恒星杀手在当前星系引爆恒星，星系被摧毁。",
	"crisis.1000":                 "无形者（The Unbidden）从异次元裂隙入侵。异次元天灾会从锚点向外扩张。",
	"crisis.1100":                 "变异者（The Aberrant）来袭。与无形者同源的敌对次元势力，可不依赖无形者触发。",
	"crisis.1200":                 "狂暴者（The Vehement）来袭。需先有无形者，三方异次元会互相开战。",
	"crisis.2000":                 "肃正协议（The Contingency）启动。机械天灾会在银河多处苏醒，并试图灭绝有机生命。",
	"crisis.7516":                 "宇宙创生结局。玩家帝国将化作堕落形态并强制切换国家，进入视界刺针终点。",
	"crisis.8005":                 "合成女王塞塔娜（Cetana）起始事件。机械纪元 DLC 的第四天灾，以纳米与合成意识席卷银河。",
	"fallen_empires_awakening.1":  "沉睡的堕落帝国苏醒。需先用 play 选中目标堕落帝国，否则当前玩家会自己觉醒。",
	"fallen_machine_empire.1":     "远古看护者苏醒。机械堕落帝国，需先用 play 切到该国。",
	"war_in_heaven.1000":          "天堂之战：唤醒两个堕落帝国并让它们互相宣战，银河被迫选边。",
	"fallen_empires_tasks.1":      "堕落帝国向玩家赠送随机礼物，态度需足够patronizing。",
	"marauder.500":                "战鼓敲响：掠夺者联合为可汗。需先用 play 选中掠夺者，否则玩家自己成为可汗。",
	"marauder.85":                 "雇佣兵飞地变得可招募。每个掠夺者帝国需分别触发。",
	"galcom.16":                   "立即创建或加入银河共同体。",
	"action.99":                   "在首都星系建立银河市场。",
	"distar.11000":                "生成封闭的 L 星团（遥远繁星）。L 星门随后可通往该星团。",
	"distar.13000":                "生成带 L 龙结局的 L 星团，并打开全部 L 星门。",
	"graygoo.1":                   "灰蛊风暴爆发。纳米机器人开始吞噬 L 星团与银河，需已有 L 星团。",
	"graygoo.100":                 "L 星团内出现德萨努谐律（Dessanu Consonance），相对友善的纳米文明。",
	"graygoo.400":                 "「静谧漫步」：选中舰队后遇见小灰（Gray）。小灰可成为传奇科学家或纳米舰队。",
	"akx.8888":                    "地平线信号。需在黑洞星系中选中有领袖的舰船。沃塔-瓦鲁姆循环与虫洞文明线由此展开。",
	"anomaly.95":                  "旅行者一号事件。银河中需存在太阳系。",
	"anomaly.186":                 "林波：可把冷冻的外星人复活为新殖民地或新帝国。选中星球。",
	"story.107":                   "安抚太空变形虫，使其不再攻击玩家。",
	"story.207":                   "安抚水晶实体。",
	"origin.5605":                 "虚境导师：解锁虚境信标恒星基地建筑。",
	"utopia.3190":                 "天选之人。虚境中诞生灵能神谕领袖。",
	"syndaw.545":                  "一个问题：开启机械觉醒相关 AI 事件局势。",
	"enclave.220":                 "可选购买雇佣兵驻地建筑。需选中星球。",
	"leviathans.1040":             "击败无限神机后的奖励事件。无限神机是围绕黑洞运行的先驱者计算机。",
	"leviathans.2101":             "击败恩尼格玛要塞后的奖励。要塞是无法移动的古代防御设施。",
	"leviathans.3102":             "自动无畏舰被停用，可开始打捞。",
	"leviathans.3103":             "修复无畏舰并收为己用。需选中舰船，可附加目标 ID。",
	"distar.137":                  "击败拾荒者机器人后的奖励。",
	"distar.212":                  "泰扬基母兽之死。需选中舰队。",
	"distar.260":                  "野生真核生物：在选中星球生成前哨物种。",
	"distar.1001":                 "失乐园：生成盖亚世界与石器时代文明。选中星球。",
	"distar.2050":                 "异星实体：生成神秘容器（Enigmatic Cache）。",
	"distar.3014":                 "尼夫拉克以敌对态度出现，好感 -100。",
	"distar.3016":                 "尼夫拉克以友好态度出现，好感 +100。",
	"distar.3055":                 "打开异星盒子，随机获得一个免费物种特质。",
	"distar.5006":                 "虚空孳孽破壳。选中星球会变为破碎世界，并放出利维坦。",
	"distar.5012":                 "击败虚空孳孽后的巨化进化奖励。",
	"ancrel.4000":                 "石中低语：在选中星球生成考古遗址。",
	"ancrel.6130":                 "扎克兰的选民：允许殖民圣地星球。",
	"cstorms.195":                 "伊尼提亚商人飞地，解锁超空间细化器建筑。宇宙风暴 3.13。",
	"grand_archive.1002":          "原始呼唤姿态。博物天枢 3.14，可解锁多种姿态建筑。",
	"precursor.98":                "生成沃尔特劳姆先驱者母星系。",
	"precursor.598":               "生成尤特先驱者母星系。",
	"precursor.1098":              "生成第一联盟母星系。",
	"precursor.1598":              "生成伊拉西安母星系。",
	"precursor.2098":              "生成赛博勒克斯母星系。",
	"utopia.3304":                 "虚境契约：低语者。可选签订，换取灵能加成与风险。",
	"utopia.3305":                 "虚境契约：织缕者。改变物种特质。",
	"utopia.3306":                 "虚境契约：噬世者。以人口换取军事力量。",
	"utopia.3307":                 "虚境契约：欲望乐器。以资源换取享乐加成。",
	"utopia.3308":                 "轮回终末契约。五十年后清算之日会到来。",
	"utopia.3320":                 "清算之日。需先接受轮回终末契约，虚境实体将清算银河。",
	"shroud.4135":                 "暗影之幕版本的轮回终末契约。",
	"distar.172":                  "神经共生。选中舰船后可获得脑蛞蝓领袖特质。需目标舰船 ID。",
	"grand_archive.2300":          "虚空之虫瘟疫。先用 play 选中虚空之虫国家，否则玩家帝国可能崩溃。",
	"fallen_empires_awakening.10": "破碎片段苏醒。暗影之幕堕落帝国，需先 play 切过去。",
	"infernals.1055":              "星系热症结局：全恒星变为红巨星，非嗜热帝国灭亡。炼狱物种包。",
	"nomads.610":                  "游牧天命：可选择成为游牧部落。游牧 DLC。",
	"ancrel.4058":                 "在选中星球获得哨兵军队。",
	"ancrel.10050":                "解锁尤特净化决议。",
	"cyber.1505":                  "机械纪元：解锁义体权威。",
	"synth.1505":                  "机械纪元：解锁合成权威。",
	"bio.170":                     "生物创生：解锁纯净权威。",
	"bio.175":                     "生物创生：解锁克隆权威。",
	"bio.180":                     "生物创生：解锁变异权威。",
	"shroud.20":                   "暗影之幕：解锁灵能权威。",
}

var loreShip = map[string]string{
	"NAME_Constructor":                "工程船。用于建造恒星基地、采矿站、研究站与巨构。",
	"NAME_Science":                    "科研船。用于调查、研究异常、考古与首次接触。",
	"NAME_Colony_Ship":                "殖民船。把人口送到可殖民星球。",
	"NAME_Transport":                  "运输船。运送陆军。",
	"NAME_Corvette":                   "护卫舰全局设计。玩家自制设计名因帝国而异，失败时用 NAME_ 前缀。",
	"NAME_Destroyer":                  "驱逐舰全局设计。",
	"NAME_Cruiser":                    "巡洋舰全局设计。",
	"NAME_Battleship":                 "战列舰全局设计。",
	"NAME_Titan":                      "泰坦。天启 DLC，可装备光矛与光环。",
	"NAME_Colossus":                   "巨像。天启 DLC，可装备星球杀手武器。",
	"NAME_Juggernaut":                 "主宰。联邦 DLC 的可移动船坞巨舰。刷出后船坞功能可能异常。",
	"Corvette":                        "使用当前玩家帝国的护卫舰设计名。",
	"Destroyer":                       "使用当前玩家帝国的驱逐舰设计名。",
	"Cruiser":                         "使用当前玩家帝国的巡洋舰设计名。",
	"Battleship":                      "使用当前玩家帝国的战列舰设计名。",
	"Titan":                           "使用当前玩家帝国的泰坦设计名。",
	"NAME_Scholar":                    "失落帝国护卫（知识管理者）。",
	"NAME_Custodian":                  "失落帝国战列巡洋舰。",
	"NAME_Reclaimer":                  "失落帝国泰坦；拾荒者机器人也使用相近名称。",
	"NAME_Glory":                      "先祖荣耀级，传送门相关事件舰。",
	"NAME_Prototype":                  "原型科研船。",
	"NAME_Alien_Ship":                 "外星科研船。",
	"NAME_Hillos_ship":                "希尔洛斯幽灵船。",
	"NAME_enigmatic_cache":            "神秘容器。环绕银河的先驱者数据库。",
	"NAME_Nanite_Factory":             "纳米机器人工厂。灰蛊相关。",
	"NAME_Nanite_Mothership":          "纳米机器人母舰。",
	"NAME_Gray_Warship":               "灰蛊战列舰。",
	"NAME_Dessanu_Mothership":         "德萨努航空母舰。",
	"NAME_Nanite_Interdictor":         "纳米拦截舰（能量武器）。",
	"NAME_Dessanu_Interdictor":        "德萨努拦截舰（动能武器）。",
	"NAME_Shroud_Avatar":              "灵能化身。虚境实体。",
	"NAME_Corrupted_Avatar":           "堕恶化身。",
	"NAME_Warped_Consciousness":       "扭曲意识。",
	"NAME_Shroud_Manifestation":       "虚境显圣。",
	"NAME_Infinity_Machine":           "无限神机。直接刷出可能是中立彩箱，需 own 接管。",
	"NAME_Herald":                     "伟大先驱考古线的先驱者级泰坦。",
	"NAME_Synth_Queen_Titan":          "合成女王塞塔娜的泰坦。机械纪元。",
	"NAME_Synth_Queen_Harrier":        "合成女王猎手护卫。",
	"NAME_Synth_Queen_Royal_Guard":    "合成女王皇家卫队。",
	"NAME_Crisis_Star_Eater":          "焚天神兵。化身天灾用于摧毁恒星的舰船。",
	"NAME_Horizon_Needle":             "视界刺针。宇宙创生飞升的终局舰船。",
	"NAME_Paladin":                    "圣骑士舰。游牧/银河守护者。",
	"NAME_Swarmling":                  "普雷索林虫群幼体。",
	"NAME_Warrior":                    "普雷索林战士。",
	"NAME_Brood_Mother":               "普雷索林育母。",
	"NAME_Queen":                      "普雷索林蜂后。",
	"NAME_Revenant":                   "无形者护卫。",
	"NAME_Phantom":                    "无形者巡洋舰。",
	"NAME_Wraith":                     "无形者战列舰。",
	"NAME_Unbidden_Anchor":            "无形者次元锚点。",
	"NAME_Sterilizer":                 "肃正协议寻求者。",
	"NAME_Euthanizer":                 "肃正协议歼灭者。",
	"NAME_AI_STARBASE":                "肃正协议恒星基地。",
	"NAME_Paradox_Titan":              "悖论泰坦。宇宙创生危机舰。",
	"NAME_Behemoth_01":                "巨兽之怒 I 级。",
	"NAME_Behemoth_02":                "巨兽之怒 II 级。",
	"NAME_Behemoth_03":                "巨兽之怒 III 级。",
	"NAME_Behemoth_04":                "巨兽之怒 IV 级。",
	"NAME_Grand_Dragon":               "以太巨龙。利维坦，盘踞在宝藏星系。",
	"NAME_Sky_Dragon":                 "天龙。起源相关的幼龙成长线。",
	"NAME_Dragonspawn":                "以太幼龙。",
	"NAME_Sky_Dragon_Baby":            "幼天龙。",
	"NAME_Nanite_Dragon":              "L 星龙。纳米构成的巨龙。",
	"NAME_Shard_Dragon":               "裂爪龙。",
	"NAME_Stellarite":                 "噬星者。以恒星为食的利维坦。",
	"NAME_Dimensional_Horror":         "位面之魇。无法移动的次元裂隙实体。",
	"NAME_Voidspawn":                  "虚空孳孽。从行星卵中孵化的利维坦。",
	"NAME_Adopted_Amoeba_Centenarian": "被收养的百年变形虫「泡泡」。",
	"NAME_Spectral_Wraith_650THz":     "蓝色光谱幽魂。",
	"NAME_Spectral_Wraith_450THz":     "红色光谱幽魂。",
	"NAME_Spectral_Wraith_520THz":     "黄色光谱幽魂。",
	"NAME_Venomous_Deity":             "毒圣。骑士与领主 DLC。",
	"NAME_Ancient_Mining_Drone":       "古代采矿无人机。",
	"NAME_Ancient_Combat_Drone":       "古代战斗无人机。",
	"NAME_Ancient_Destroyer":          "古代驱逐舰。",
	"NAME_Crystal_Nidus":              "水晶巢穴。",
	"NAME_Large_Crystal_Entity_Blue":  "大型蓝色水晶实体。",
	"NAME_Medium_Crystal_Entity_Blue": "中型蓝色水晶实体。",
	"NAME_Small_Crystal_Entity_Blue":  "小型蓝色水晶实体。",
	"NAME_Tiyanki_Matriarch":          "泰扬基母兽。",
	"NAME_Dreadnought":                "自动无畏舰。可打捞修复。",
	"NAME_Enigmatic_Fortress":         "恩尼格玛要塞。通常无法移动。",
}

var loreMega = map[string]string{
	"dyson_sphere_5":                  "完成的戴森球。包裹恒星，提供巨额能量币。",
	"dyson_sphere_0":                  "戴森球基址。",
	"dyson_sphere_ruined":             "损坏的戴森球，可修复。",
	"ring_world_1":                    "环形世界基址。",
	"ring_world_3_intermediate":       "环形世界完成区段。重复执行可生成下一区段。",
	"ring_world_ruined":               "损坏的环形世界。",
	"think_tank_3":                    "完成的科学枢纽，提供巨额研究。",
	"think_tank_0":                    "科学枢纽基址。",
	"think_tank_ruined":               "损坏的科学枢纽。",
	"spy_orb_4":                       "完成的哨兵阵列，提供全图传感器。",
	"gateway_final":                   "完成的星门，可连接其它星门。",
	"gateway_ruined":                  "损坏的星门。",
	"lgate_base":                      "L 星门。通往 L 星团，遥远繁星 DLC。",
	"habitat_0":                       "3.9 前的旧轨道居住站。",
	"habitat_central_complex":         "居住站核心。3.9 居住站重做后的主体。",
	"matter_decompressor_4":           "完成的物质解压器。需黑洞星系，提供巨额矿物。",
	"strategic_coordination_center_3": "战略指挥中心。寰宇企业，提升海军容量与舰船速度。",
	"mega_art_installation_4":         "巨型艺术设施，提供凝聚力与幸福。",
	"interstellar_assembly_4":         "星际集会，提供外交权重。",
	"mega_shipyard_3":                 "巨型船坞。纵横捭阖，加快舰船建造。",
	"hyper_relay":                     "超空间中继器。纵横捭阖，加快帝国内部旅行。",
	"quantum_catapult_3":              "量子弹弓。需中子星或脉冲星，可把舰队弹射到远方。",
	"orbital_ring":                    "星环基地。需选中行星，错误时可能生成在恒星上。",
	"crisis_sphere_4":                 "完成的以太相引擎。化身天灾终局巨构。",
	"dyson_swarm_3":                   "完成的戴森云。机械纪元，环绕恒星的能量阵列。",
	"orbital_arc_furnace_4":           "完成的电弧熔炉。需熔岩行星。",
	"cosmogenesis_world_0":            "突触凝练机。宇宙创生飞升的人口转化巨构。",
	"cosmogenesis_needle":             "视界刺针巨构。",
	"shroud_seal":                     "虚境封印。",
	"grand_archive_0":                 "博物天枢。寰宇博物 DLC 的标本博物馆。",
	"deep_space_citadel_0":            "深空擎天堡基址。",
	"behemoth_egg":                    "噬天巨兽卵。万物竞发。",
	"behemoth_chrysalis":              "巨兽蛹。",
	"galactic_crucible_4":             "完成的星系熔炉。炼狱物种包。",
	"dyson_gun_0":                     "星际大炮基址。游牧 DLC，需中子星。",
	"dyson_gun_3":                     "完成的星际大炮。",
	"civilian_arkship_megastructure":  "民用方舟。游牧。",
	"military_arkship_megastructure":  "军用方舟。",
	"science_arkship_megastructure":   "科研方舟。",
	"waystation_megastructure":        "驿站。",
}
