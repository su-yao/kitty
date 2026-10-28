# User Instruction Memory

此文件记录用户指令、偏好，以及执行任务中沉淀的项目知识。

## Entries

[Project Knowledge Summary]
- Date: 2026-10-18
- Context: Discovered by Agent while building and publishing the Windows installer
- Category: Build Methods / Operations & Deployment
- Instructions:
  - 构建：`bash scripts/build.sh`，读取仓库根 `VERSION` 注入主程序与安装器，并执行 go test、go vet（含 windows 交叉）与 -trimpath
  - 本地预览：`PORT=<端口> XDG_CONFIG_HOME=/tmp/<dir> ./dist/linux/StellarisConsolePanel-<VERSION>`
  - 发布：`GH_TOKEN=$(printf 'protocol=https\nhost=github.com\n\n' | git credential fill | sed -n 's/^password=//p') gh release create <tag> dist/windows/StellarisConsolePanel-Setup-<VERSION>.exe dist/windows/StellarisConsolePanel-<VERSION>.exe --repo su-yao/kitty`
  - 分发产物文件名在后缀前带版本号：`StellarisConsolePanel-<VERSION>.exe`、`StellarisConsolePanel-Setup-<VERSION>.exe`、Linux `StellarisConsolePanel-<VERSION>`
  - 安装后的程序名与快捷方式保持 ASCII 无版本号：`StellarisConsolePanel.exe`、`StellarisConsolePanel.lnk` / `Uninstall-StellarisConsolePanel.lnk`
  - 发布流程：先提交并 `git push origin master`，再 `gh release create <tag> dist/windows/StellarisConsolePanel-Setup-<VERSION>.exe dist/windows/StellarisConsolePanel-<VERSION>.exe --repo su-yao/kitty`
  - 环境无 wine/NSIS/Inno，Windows 侧只能靠交叉编译、`GOOS=windows go vet`、PE 头校验与单测验证

[Project Knowledge Summary]
- Date: 2026-10-18
- Context: Discovered by Agent while editing the frontend
- Category: Workflow & Collaboration
- Instructions:
  - 前端源文件在 `web/`，`cmd/panel/web/` 是 embed 副本；改完必须同步，`TestEmbeddedWebMatchesSource` 会校验两者一致
