# 群星控制台面板 Stellaris Console Panel

Windows 桌面工具：把《群星》原版与 MOD 的控制台命令做成分类按钮，按下后发送到正在运行的游戏控制台。

适配版本：3.0.0 至 4.4.6（含 4.5.0 条目）。

## 功能

- 每条命令一个独立按钮，点击即发送
- 分类分页：资源、舰船、天灾舰、利维坦、巨构、事件、科技、遗珍、行星、帝国、调试、MOD
- 扫描游戏根目录与创意工坊 MOD，只保留有中文名或可实际刷出的条目
- 按游戏版本过滤命令（例如女王泰坦自 3.12 起，星际大炮自 4.4 起）
- 管理员权限下向群星窗口注入按键

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

产物：

- `dist/windows/Setup.exe` 标准安装包（双击安装）
- `dist/windows/StellarisConsolePanel.exe` 便携主程序
- `dist/linux/StellarisConsolePanel` 预览

## 使用注意

- 控制台在铁人模式和多人游戏中不可用
- 若命令没进去：先点「打开控制台」，确认控制台出现后再点命令
- `instant_build` 对 AI 同样生效，用完再点一次关闭
- 部分事件需要先选中星球、舰队或局势
- MOD 扫描结果进入「MOD扫描」分类
