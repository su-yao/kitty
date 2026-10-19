# Requirements Document

## Introduction

Stellaris Console Panel 是一款 Windows 桌面工具。它收集《群星》游戏目录与 MOD 中的控制台命令（舰船、巨构、事件等），按版本与分类分页展示为独立按钮；按下按钮后将命令发送到正在运行的群星控制台执行。

## Glossary

- **System**: Stellaris Console Panel
- **Game Root**: 群星安装根目录（含 `stellaris.exe`、`common/`、`events/`）
- **Mod Root**: 创意工坊或用户 MOD 目录
- **Console Command**: 可在群星 `~` 控制台中执行的一行文本
- **Version Pack**: 适配指定游戏版本区间的命令集合
- **Injector**: 将按键/文本发送到群星窗口的模块

## Requirements

### Requirement 1

**User Story:** AS 玩家, I want 带界面的 Windows 安装包, so that 我可以一键安装并使用本工具

#### Acceptance Criteria

1. WHEN 用户运行安装程序, THE System SHALL 请求管理员权限并将主程序安装到用户选择的目录
2. WHEN 安装进行中, THE System SHALL 创建开始菜单与桌面快捷方式，并写入「设置-应用」卸载项
3. WHEN 用户请求卸载, THE System SHALL 移除已安装文件、快捷方式与注册表卸载项
4. WHEN 用户以 `/S` 参数运行, THE System SHALL 执行无人值守静默安装
5. WHEN 用户指定安装目录, THE System SHALL 兼容带引号与不带引号的 `/DIR=路径`
6. WHILE 程序运行中, THE System SHALL 只允许一个面板实例，避免重复注入按键
7. WHEN 卸载完成, THE System SHALL 只清空安装目录内的文件，保留用户放入的其他文件

### Requirement 2

**User Story:** AS 玩家, I want 每条控制台命令对应一个独立按钮, so that 我按下即可执行

#### Acceptance Criteria

1. WHEN 用户点击某条命令按钮, THE System SHALL 把该按钮绑定的完整命令发送到群星控制台
2. WHILE 发送过程进行中, THE System SHALL 在界面显示发送状态与命令原文，并暂时禁用其他命令按钮
3. IF 未找到群星窗口, THE System SHALL 提示用户先启动游戏并保持非铁人模式
4. WHEN 命令包含当前键盘布局无法直接键入的字符, THE System SHALL 整条改用剪贴板粘贴，避免重复输入
5. WHEN 多个发送请求同时到达, THE System SHALL 串行执行，按键不得交错
6. WHEN 命令需要选中对象或目标 ID, THE System SHALL 在按钮上标注前置条件
7. WHEN 舰船事件缺少目标 ID, THE System SHALL 阻止发送并提示先获取舰船 ID

### Requirement 3

**User Story:** AS 玩家, I want 命令按分类分页, so that 我能快速找到舰船、巨构、事件等命令

#### Acceptance Criteria

1. WHEN 用户打开主界面, THE System SHALL 按分类展示命令（资源、舰船、巨构、事件、科技、遗珍、调试、MOD）
2. WHEN 单页命令超过 24 条, THE System SHALL 提供分页控件
3. WHEN 用户在搜索框输入关键字, THE System SHALL 按中文名、英文名、命令文本过滤结果
4. WHEN 窗口尺寸或缩放变化, THE System SHALL 按 27 寸 2K 与 31.5 寸 4K、100% 与 150% 重新计算平铺网格，使按键铺满一屏

### Requirement 4

**User Story:** AS 玩家, I want 从 3.0.0 到 4.4.6 及更新版本分版本适配, so that 当前游戏版本只显示可用命令

#### Acceptance Criteria

1. WHEN 用户选择或检测到游戏版本, THE System SHALL 只展示该版本区间内有效的命令
2. WHEN 某命令在该版本不可用, THE System SHALL 隐藏该命令
3. WHEN 游戏版本无法自动检测, THE System SHALL 允许用户手动选择版本

### Requirement 5

**User Story:** AS 玩家, I want 扫描游戏根目录与 MOD, so that 原版与模组舰船/巨构/事件都能变成按钮

#### Acceptance Criteria

1. WHEN 用户指定游戏根目录并点击扫描, THE System SHALL 解析 `common/global_ship_designs`、`common/megastructures`、`events` 生成命令
2. WHEN 用户指定 MOD 目录并扫描, THE System SHALL 将 MOD 内同类定义合并到「MOD」分类
3. WHEN 扫描结果与内置命令重复, THE System SHALL 跳过该结果，保证每条命令只有一个按钮
4. WHEN 事件没有本地化标题, THE System SHALL 不展示该事件，避免出现无法辨认的内部 ID
5. IF 扫描失败, THE System SHALL 保留内置命令库并显示失败原因
6. WHEN 用户只填写 MOD 目录而未填写游戏根目录, THE System SHALL 仍允许扫描 MOD
7. WHEN 扫描到游戏根目录, THE System SHALL 按事件文件的作用域标注星球/舰队/舰船事件，并校验内置事件 ID 与舰船设计是否存在
8. WHEN 内置事件或舰船设计在当前游戏文件中不存在, THE System SHALL 标记为不可用并允许一键隐藏

### Requirement 6

**User Story:** AS 玩家, I want 以管理员权限把命令送进群星控制台, so that 全屏或后台时仍能执行

#### Acceptance Criteria

1. WHEN 用户点击发送, THE System SHALL 激活群星窗口、打开控制台、输入命令并回车
2. WHILE 以管理员身份运行, THE System SHALL 使用 Win32 SendInput 或剪贴板粘贴完成输入
3. IF 键盘布局导致 `~` 无法打开控制台, THE System SHALL 允许用户自定义开控制台按键
4. WHEN 用户保存设置, THE System SHALL 原子写入配置并在重启后恢复，未提交的字段保持原值
5. WHEN 界面加载完成, THE System SHALL 显示当前程序版本
