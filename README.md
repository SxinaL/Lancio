# Lancio - 桌面背单词

基于 **Go + Wails v2** 构建的桌面背单词应用，以右下角浮窗形式展示单词，支持全局划词翻译。

## 功能特点

- 🪟 **右下角浮窗** — 启动后自动定位到屏幕右下角，无边框透明窗口
- 📖 **词库管理** — 支持创建/删除自定义词库，分类管理单词
- ✏️ **单词管理** — 完整的单词 CRUD（增删改查），支持收藏/取消收藏
- 🎲 **随机展示** — 按词库随机切换单词，展示音标、释义、例句
- 🔍 **全局划词翻译** — 按住 `Ctrl` + 鼠标划词，自动获取选中文本并翻译
- 🌐 **在线词典** — 短文本走有道词典（音标+释义+例句），长文本走百度翻译
- ⚡ **三级缓存** — 内存缓存 → SQLite 缓存 → 在线 API，查询毫秒级响应
- 📌 **窗口置顶/置底** — 一键切换浮窗桌面工具模式
- 🎚️ **透明度调节** — 自由调节浮窗透明度
- 🖱️ **拖拽缩放** — 支持拖拽移动和边缘缩放窗口
- 📋 **系统托盘** — 最小化到托盘，后台静默运行

## 技术栈

- **后端**: Go + Wails v2
- **前端**: Vue 3 + Vue Router + Vite 7
- **数据库**: SQLite（glebarez/sqlite + GORM）
- **系统交互**: Windows API（gonutz/w32）
- **在线 API**: 有道词典 API + 百度翻译 API

## 架构设计

Lancio 采用**双窗口多进程**架构：

- **主窗口（Main Window）** — 800×600，负责词库管理、单词管理、浮窗配置
- **浮窗（Float Window）** — 可配置尺寸，无边框透明，负责单词展示与划词翻译

主窗口启动时根据配置自动拉起浮窗子进程，关闭主窗口时终止所有浮窗。

## 项目结构

```
Lancio/
├── main.go                         # 应用入口，主窗/浮窗双进程启动
├── go.mod / go.sum                 # Go 依赖
├── wails.json                      # Wails 构建配置
├── api/                            # 外部 API 调用
│   ├── request_baidu.go            # 百度翻译 API
│   └── request_youdao.go           # 有道词典 API
├── database/
│   ├── entity/                     # 数据模型定义
│   │   ├── dictionary.go           # 词典缓存实体
│   │   ├── vocabulary.go           # 单词实体 & 接口
│   │   └── vocabularybank.go       # 词库实体 & 接口
│   └── method/                     # 数据库操作
│       ├── db_method.go            # 数据库初始化 & 表迁移
│       ├── default_words.go        # 默认种子数据
│       ├── dictionary_method.go    # 词典缓存 CRUD
│       ├── vocabulary_method.go    # 单词 CRUD
│       └── vocabularybank_method.go # 词库 CRUD
├── handler/                        # Wails 绑定层（暴露给前端）
│   ├── dictionary_handler.go       # 在线词典查询
│   ├── taskbar_handler.go          # 系统托盘（最小化/恢复）
│   ├── vocabulary_handler.go       # 单词增删改查
│   ├── vocabularybank_handler.go   # 词库管理 & 随机取词
│   └── window_handler.go           # 窗口控制（置顶/透明度/位置）
├── model/                          # 共享模型 & 全局状态
│   ├── app.go                      # App 全局实例
│   ├── config.go                   # 应用配置结构
│   ├── dictionary.go               # 词典查询结果模型
│   ├── hook.go                     # 划词钩子相关模型
│   └── window.go                   # 窗口实例
├── services/                       # 业务逻辑层
│   ├── app_service.go              # 应用生命周期（Startup/DomReady/Shutdown）
│   ├── dictionary_services.go      # 词典查询（三级缓存 + 在线）
│   └── hook_service.go             # 全局鼠标钩子 & 划词翻译
├── utils/                          # 工具函数
│   ├── cache_utils.go              # 内存缓存
│   ├── clipboard_utils.go          # 剪贴板操作
│   ├── config_utils.go             # 配置文件读写
│   ├── data_utils.go               # 数据处理
│   ├── dictionary_utils.go         # 词典工具
│   ├── file_utils.go               # 文件路径工具
│   ├── hook_utils.go               # Windows 全局钩子实现
│   ├── key_utils.go                # 键盘工具
│   ├── net_utils.go                # 网络工具
│   ├── taskbar_utils.go            # 系统托盘工具
│   └── window_utils.go             # 窗口操作工具
├── frontend/
│   ├── package.json                # Vue 3 + Vite 7 + Vue Router
│   └── src/
│       ├── main.js                 # Vue 应用入口
│       ├── App.vue                 # 根组件（根据窗口模式动态路由）
│       ├── store.js                # 全局状态管理
│       ├── router/index.js         # Vue Router 路由配置
│       ├── views/                  # 页面视图
│       │   ├── MainWin.vue         # 主窗口布局
│       │   ├── BankManager.vue     # 词库管理页
│       │   ├── VocManager.vue      # 单词管理页
│       │   ├── FloatWinManager.vue # 浮窗管理页
│       │   └── VocFloatWin.vue     # 浮窗展示页
│       ├── components/             # 组件
│       │   ├── mainwin/            # 主窗口组件（侧边栏、工具栏）
│       │   ├── bankmanager/        # 词库管理组件（添加/删除面板）
│       │   ├── vocmanager/         # 单词管理组件（表格/编辑/删除）
│       │   ├── vocfloatwin/        # 浮窗组件（单词卡片/操作栏/设置）
│       │   └── floatwinmanager/    # 浮窗管理组件
│       └── template/               # 公共模板组件
└── build/
    ├── appicon.png                 # 应用图标
    ├── darwin/                     # macOS 构建配置
    ├── windows/                    # Windows 构建配置 & NSIS 安装器
    └── bin/                        # 编译输出目录
```

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- Wails CLI v2

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式运行

```bash
cd Lancio
wails dev
```

### 构建生产版本

```bash
wails build
```

构建后的可执行文件在 `build/bin/` 目录下。

## 后端 API

### 单词管理（VocabularyHandler）

| 方法 | 说明 |
|------|------|
| `GetWordByID(id)` | 根据 ID 获取单词 |
| `GetWordsByRange(offset, limit)` | 分页获取单词列表 |
| `GetWordCount()` | 获取单词总数 |
| `GetWordsByVocabularyBankID(vocabularyBankID)` | 根据词库 ID 获取单词列表 |
| `AddVocabulary(vocabulary)` | 新增单词 |
| `UpdateVocabulary(vocabulary)` | 更新单词 |
| `DeleteVocabulary(id)` | 删除单词 |
| `GetStarStatus(word)` | 获取单词收藏状态 |

### 词库管理（VocabularyBankHandler）

| 方法 | 说明 |
|------|------|
| `GetAllVocabularyBanks()` | 获取所有词库 |
| `CreateVocabularyBank(name, description)` | 创建词库 |
| `DeleteVocabularyBank(id)` | 删除词库 |
| `GetVocabularyCountByVocabularyBankID(id)` | 获取词库下单词数量 |
| `GetRandomVocabulary()` | 随机获取单词（按选中词库） |
| `GetSelectedVocabularyBankID()` | 获取当前选中词库 ID |
| `SetSelectedVocabularyBankID(id)` | 设置当前选中词库 |

### 窗口控制（WindowHandler）

| 方法 | 说明 |
|------|------|
| `SetWindowOpacity(opacity)` | 设置窗口透明度 |
| `GetWindowOpacity()` | 获取窗口透明度 |
| `PinFloatWindow()` / `UnpinFloatWindow()` | 切换窗口置顶/置底 |
| `ToggleFloatWindow(windowName)` | 开关浮窗子进程 |
| `SaveFloatWindowPositionAndSize()` | 保存浮窗位置和尺寸 |

### 系统托盘（TaskbarHandler）

| 方法 | 说明 |
|------|------|
| `MinimizeToTray()` | 最小化到系统托盘 |
| `RestoreFromTray()` | 从系统托盘恢复 |

### 在线词典（DictionaryHandler）

| 方法 | 说明 |
|------|------|
| `LookupWord(query, language)` | 在线查询单词/文本（含三级缓存） |

## 配置说明

应用配置存储在 `config.json`（位于可执行文件同目录），支持以下配置项：

```json
{
  "windowFloatConfig": {
    "LancioFloat": {
      "windowOpacity": 0.4,
      "windowSizeW": 280,
      "windowSizeH": 300,
      "selectedVocabularyBankID": 0,
      "isPinned": false,
      "launch": true
    }
  },
  "apiConfigKey": {
    "appid": "",
    "key": ""
  }
}
```

- `selectedVocabularyBankID`: 0 表示从所有词库随机取词；> 0 表示从指定词库取词
- `apiConfigKey`: 百度翻译 API 密钥（用于长文本翻译），有道词典无需密钥

## 自定义词库

1. 启动应用，进入主窗口 → 词库管理
2. 点击「添加」创建新词库
3. 进入对应词库，通过「添加单词」手动录入或使用划词翻译自动收录

## 许可证

MIT