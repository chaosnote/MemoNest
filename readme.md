# 🪶 MemoNest

MemoNest 是一款以 Go 語言打造的模組化筆記系統，支援使用者註冊、登入、文章分類、富文字編輯與圖片上傳。系統採用 Hexagonal Architecture，整合 MariaDB、Redis、MongoDB、NATS 等多種後端服務，並以 Uber fx 進行依賴注入與模組管理，具備高可維護性與擴充性。

## 📦 專案架構

MemoNest 是一個模組化、可擴充的筆記系統，支援使用者註冊、登入、文章分類、富文字編輯與圖片上傳等功能。系統採用 Go 語言開發，整合 MariaDB、Redis、MongoDB、NATS 等多種後端服務，並以 Hexagonal Architecture 與 Uber fx 進行依賴注入與模組管理。

``` dir
.
├── adapter/           # 外部介面與基礎設施實作（HTTP、DB、Redis、NATS、Mongo）
├── api/http/          # HTTP 路由與處理器（Gin）
├── application/       # Usecase 應用邏輯層
├── config/            # 設定檔與 CLI 參數解析
├── domain/            # 核心領域模型與介面定義
├── utils/             # 工具函式（加解密、時間、檔案、模板等）
├── web/templates/     # HTML 模板（Bootstrap + Quill 編輯器）
├── doc/db/            # SQL 建表與測試資料
├── assets/            # 設定檔（如 config.json）
├── dist/              # 日誌與圖片上傳目錄
├── cmd/main.go        # 程式進入點
├── go.mod             # Go module 設定
└── start.sh           # 啟動腳本
```

🚀 快速開始

搭配 WSL

``` shell
wsl -d Ubuntu-24.04 --cd "docker 目錄位置"
./start.sh

wsl -d Ubuntu-24.04 --cd "golang 專案位置"
./start.sh
```

``` ConEmu
%windir%\system32\wsl.exe -cur_console:t:Linux -d Ubuntu-24.04 --cd "...\MemoNest\docker\"
%windir%\system32\wsl.exe -cur_console:t:Golang -d Ubuntu-24.04 --cd "...\MemoNest\work\"
```

開啟瀏覽器前往 http://localhost:8080

🔐 功能特色

1. 使用者註冊 / 登入 / 登出（支援 AES 加密與「記住我」功能）
1. 分類節點管理（支援巢狀分類、拖曳搬移、編輯、刪除）
1. 文章 CRUD（支援 Quill 富文字編輯器）
1. 圖片上傳與 Base64 轉檔儲存
1. Redis Session 管理與 IP 驗證
1. NATS、MongoDB、MariaDB 整合（可擴充）
1. 模組化架構，支援 fx DI 與 Hexagonal Architecture

🧱 架構設計

1. Hexagonal Architecture：明確區分 adapter、application、domain 三層，提升可測試性與可維護性。
1. fx DI：使用 Uber fx 進行依賴注入與生命週期管理。
1. Gin Framework：作為 HTTP 框架，搭配 middleware 管理 session、IP 驗證與錯誤攔截。
1. Template Engine：使用 html/template 搭配 Bootstrap 與 Quill 提供簡潔 UI。

🧪 測試資料
可使用 doc/db/ 中的 SQL 檔案建立測試資料表與預設帳號。
