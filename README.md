# 群星控制台面板 Stellaris Console Panel

Windows 桌面工具：把《群星》原版与 MOD 的控制台命令做成分类按钮，按下后发送到正在运行的游戏控制台。

适配版本：3.0.0 至 4.4.6（含 4.5.0 条目）。

## 功能

- 每条命令一个独立按钮，点击即发送
- 分类分页：资源、舰船、天灾舰、利维坦、巨构、事件、科技、遗珍、行星、帝国、调试、MOD
- 扫描游戏根目录、创意工坊与额外 MOD 目录，只保留有中文名或可实际刷出的条目
- 扫描结果自动与内置库去重，同一条命令只保留一个按钮
- 按游戏版本过滤命令（例如女王泰坦自 3.12 起，星际大炮自 4.4 起）
- 按钮标注前置条件：`选舰队`、`选星球`、`先选目标`、`需 DLC`；扫描游戏目录后会自动校验事件 ID 与舰船设计是否存在，找不到的标为 `未找到`
- 管理员权限下向群星窗口注入按键，并发点击会串行发送
- 设置自动保存，界面底部显示当前程序版本

示例：

```
女王泰坦 / add_ship NAME_Synth_Queen_Titan
戴森球 / create_megastructure dyson_sphere_5
合成女王 / event crisis.8005
```

## 安装（Windows）

只需要 `Setup.exe`（自包含安装包，已嵌入主程序）。

1. 双击 `Setup.exe`，在 UAC 提示选择「是」
2. 确认安装目录（默认 `C:\Program Files\StellarisConsolePanel`）
3. 选择是否创建开始菜单 / 桌面快捷方式
4. 安装完成后从开始菜单打开「群星控制台面板」
5. 启动群星（关闭铁人与多人），点「打开控制台」或手动按 `~`，再点命令按钮

静默安装：

```
Setup.exe /S
Setup.exe /S /DIR="D:\Games\StellarisConsolePanel"
```

卸载：

- 「设置 - 应用」里卸载「群星控制台面板」
- 或开始菜单「卸载 群星控制台面板」
- 或 `Uninstall.exe /S`

### 可选：校验命令是否适用于当前游戏

首次扫描游戏根目录后，面板会读取 `events/`、`common/global_ship_designs`、`common/ship_designs`：
- 事件会按文件中的 `country_event` / `planet_event` / `ship_event` 自动标注作用域
- 内置事件 ID 或 `NAME_` 舰船设计在当前版本不存在时标为 `未找到`
- 勾选「隐藏不可用」可只显示校验通过的按钮

## Linux 预览

当前环境可运行 Linux 预览版查看界面（不会注入 Windows 按键）：

```
# 运行测试
go test ./internal/catalog ./internal/scan

# 启动预览
./dist/linux/StellarisConsolePanel
```

## 从源码打包 Windows EXE

需要 Go 1.22+。在任意系统交叉编译：

```
# 测试并交叉编译 Windows 可执行文件
bash scripts/build.sh
```

版本号统一写在仓库根 `VERSION`，构建时注入主程序与安装器，无需改代码。

产物：

- `dist/windows/Setup.exe` 标准安装包（双击安装）
- `dist/windows/StellarisConsolePanel.exe` 便携主程序
- `dist/linux/StellarisConsolePanel` 预览

## 使用注意

- 控制台在铁人模式和多人游戏中不可用
- `event` 命令按事件作用域生效：国家事件可直接触发；星球事件要先选中星球；舰船事件必须填「目标 ID」（控制台输入 `debugtooltip`，悬停舰船查看 ID），例如 `event crisis.105 <舰船ID>`
- `add_ship` 会在当前选中的舰队、星球或星系旁生成舰船，先选中一个目标再点按钮
- 部分事件和舰船需要 DLC 或特定局势；按钮右上角会标出 `需 DLC` 与 `未找到`
- 首次使用先点「打开控制台」，确认控制台出现后再点命令；之后保持「已打开，只输入回车」
- 若命令没进去：先点「检测窗口」确认找到 `stellaris.exe`，再点一下群星画面让窗口置于最前
- 命令含特殊字符时自动改用剪贴板粘贴，会覆盖系统剪贴板
- `instant_build` 对 AI 同样生效，用完再点一次关闭
- 部分事件需要先选中星球、舰队或局势
- MOD 扫描结果进入「MOD扫描」分类，重复命令会被自动跳过
