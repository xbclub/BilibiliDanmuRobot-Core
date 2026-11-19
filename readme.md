# BilibiliDanmuRobot-Core

B 站直播弹幕机器人的**核心库**，被 [BilibiliDanmuRobot](https://github.com/xbclub/BilibiliDanmuRobot)（Wails GUI）作为 Go module 依赖引入，本仓库自身不直接启动运行。

- 环境：go 1.24+
- 鸣谢：https://github.com/Akegarasu/blivedm-go （内置于 `blivedm-go/` 的本地 fork）

## 目录结构

```
biliapi/     B站直播开放接口客户端 + 登录态（cookie/扫码）管理
qrcode/      终端二维码渲染（扫码登录用）
handler/     对外门面：WsHandler 生命周期、装配、配置热加载
logic/       中央弹幕发送队列（PushToBulletSender）与切分工具
internal/features/   内部功能模块（不对外暴露）
  welcome/     进场欢迎（特效/INTERACT_WORD/天选互锁）
  thanksgift/  礼物/舰长/红包答谢与盲盒盈亏播报
  pk/          PK 对手信息播报
  robot/       AI 聊天（ChatGPT/青云客）
  danmuproc/   弹幕解析分发：签到/抽签/关键词/弹幕计数/盲盒查询/主播指令
  notice/      禁言提醒等杂项
svc/         ServiceContext（共享状态与数据模型）
config/      配置定义（go-zero conf，yaml 加载）
entity/      B站消息与接口的数据结构
model/       SQLite 数据模型（签到/弹幕计数/盲盒统计）
blivedm-go/  弹幕 WebSocket 客户端（fork）
```

## 开发调试方式

```bash
go mod tidy
go run test/test.go
```

调试入口在 `test/test.go`，需要先在 `token/` 下放置登录凭证（`bili_token.txt` / `bili_token.json`）、`etc/` 下放置配置文件。配置模板见 [etc/bilidanmaku-api.example.yaml](etc/bilidanmaku-api.example.yaml)，复制为 `etc/bilidanmaku-api.yaml` 即可使用。

### dry-run 模式（GUI 调试）

弹幕发送出口提供进程内 dry-run 开关（**不暴露为配置项**，配置文件不会出现）：`NewWsHandler` 之前调用一次

```go
biliapi.EnableDryRun()   // 或 DisableDryRun() 关闭；DryRunEnabled() 查询
```

开启后所有弹幕（欢迎/感谢/机器人/指令/下播告别）不请求 B 站，只打 `[dry-run]` 前缀日志，`Send` 按成功返回，发送队列按正常 3.5s 节奏排空，可完整模拟运行时序。

## 三档进房欢迎（PR #33）

在原有 `InteractWord` 欢迎之上新增三档分流，默认关闭，旧配置行为不变：

```yaml
InteractWord: true
WelcomeTier: true          # 三档欢迎开关
WelcomeFanMedalLevel: 0    # 老粉最低粉丝牌等级，0 不限
WelcomeNew:                # 第一次来的观众
  - "欢迎呀~{at}恭喜你发现了{anchor}"
WelcomeBack:               # 来过但没戴本直播间粉丝牌
  - "哈喽呀~{at}等你很久喽~"
WelcomeFan:                # 戴着本直播间粉丝牌（老粉）
  - "哈喽{at}今天有没有更想{anchor}呀~"
```

- 匹配顺序：PK 串门 → 指定 UID → 老粉 → 新人 → 回头客；某档模板为空时跳过该档，全部为空回退 `WelcomeDanmu` / `WelcomeDanmuByTime`
- 占位符：`{at}`（@ 模式下自动省略）、`{user}`、`{anchor}`、`{medal}`、`{medal_level}`
- 进场记录存于 SQLite `visit_<房间号>` 表（`uid` / `first_seen` / `last_seen` / `visit_cnt`），90 天未访问自动清理；V1/V2 双事件 10 秒短窗去重

## 依赖方（GUI）升级指引

v1.1.21 之后的版本有两个包路径变更，GUI 仓库需同步调整 import：

| 旧 | 新 |
|---|---|
| `github.com/xbclub/BilibiliDanmuRobot-Core/http` | `github.com/xbclub/BilibiliDanmuRobot-Core/biliapi` |
| `github.com/xbclub/BilibiliDanmuRobot-Core/utiles` | `github.com/xbclub/BilibiliDanmuRobot-Core/qrcode` |

调用点替换：`http.Xxx(...)` → `biliapi.Xxx(...)`；`utiles.GenerateQr(msg)` → `qrcode.Generate(msg)`。
其余 import（`handler`、`logic`、`svc`、`config`、`entity`）路径不变，函数签名不变。
