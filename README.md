# DesktopVoc - 桌面背单词

基于 **Go + Wails2** 构建的桌面背单词应用，以右下角浮窗形式展示单词。

## 功能特点

- 🪟 **右下角浮窗** — 启动后自动定位到屏幕右下角
- 📖 **本地词库** — 内置 110 个常用英语单词，SQLite 存储
- 🎲 **随机展示** — 每次点击「下一个」随机切换单词
- 🔍 **点击查看释义** — 点击单词区域切换显示/隐藏翻译和例句
- 📌 **窗口置顶** — 点击 📌 按钮切换窗口置顶
- 🖱️ **拖拽移动** — 拖动顶部区域可自由移动窗口

## 技术栈

- **后端**: Go + Wails v2
- **前端**: Vanilla JS + Vite
- **数据库**: SQLite (modernc.org/sqlite)
- **UI**: 无边框透明浮窗

## 项目结构

```
DesktopVoc/
├── main.go                 # 应用入口 & Wails 配置
├── app.go                  # 应用逻辑 (API 绑定)
├── database/
│   ├── database.go         # SQLite 数据库操作
│   ├── words.go            # 单词模型 & 接口定义
│   └── default_words.go    # 默认词库 (110个单词)
├── services/
│   └── dictionary.go       # 词典查询接口 (预留扩展)
├── frontend/
│   ├── index.html
│   ├── src/
│   │   ├── main.js         # 前端逻辑 (浮窗交互)
│   │   └── style.css       # 浮窗样式
│   └── package.json
└── build/
    └── bin/
        └── DesktopVoc.exe  # 编译后的可执行文件
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 16+
- Wails CLI v2

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式运行

```bash
cd DesktopVoc
wails dev
```

### 构建生产版本

```bash
cd DesktopVoc
wails build
```

构建后的可执行文件在 `build/bin/` 目录下，可直接运行。

## 后端 API

| 方法 | 说明 |
|------|------|
| `GetRandomWord()` | 获取一个随机单词 |
| `GetWordByID(id)` | 根据 ID 获取单词 |
| `GetWordsByRange(offset, limit)` | 分页获取单词列表 |
| `GetWordCount()` | 获取单词总数 |
| `LookupWordDetail(word)` | 查询单词详细信息 **(预留接口)** |

## 扩展：接入在线词典

`services/dictionary.go` 中定义了 `DictionaryService` 接口，实现它即可接入在线词典 API：

```go
type DictionaryService interface {
    Lookup(word string) (*WordDetail, error)
}
```

实现后在 `app.go` 的 `startup` 中将 `a.dict` 替换为你的实现即可。

## 自定义词库

编辑 `database/default_words.go`，在 `GetDefaultWords()` 中添加或修改单词，然后重新构建。
