# [Deeli Ai] Golang 迷你專案實作作業   

## 作業目標   
設計一個後端 API 服務，讓使用者可以儲存感興趣的文章連結，系統會自動抓取文章的 Metadata，根據使用者行為進行推薦，並在爬取失敗時自動重試。   
   
### 功能需求   
1. 使用者註冊與登入   
    | POST /signup | 註冊新帳號               |      |
    | :----------- | :----------------------- | :--- |
    | POST /login  | 登入，取得 Token         |      |
    | GET /me      | 回傳目前登入的使用者資訊 |      |

2. 文章收藏與 Metadata 抓取   
    | POST /articles       | 傳入 URL，儲存連結並嘗試抓取 title/description/image |                        |
    | :------------------- | :--------------------------------------------------- | :--------------------- |
    | GET /articles        | 取得使用者儲存的所有文章清單（支援簡單分頁）         | 只看得到自己儲存的文章 |
    | DELETE /articles/:id | 刪除使用者的某篇文章                                 |                        |

    - 當使用者儲存文章時，系統會擷取網站 Metadata (OpenGraph 或其他 HTML metadata)   
    - 若失敗，需有 background worker 重試機制   
3. 標籤與評分   
    | POST /articles/:id/rate   | 為該使用者收藏的文章設定評分（1–5） | 只能評分自己的文章 |
    | :------------------------ | :---------------------------------- | :----------------- |
    | GET /articles/:id/rate    | 取得該使用者對文章的評分            |                    |
    | DELETE /articles/:id/rate | 取消評分                            |                    |

4. 推薦文章   
    | GET /recommendations | 根據使用者評分與標籤推薦其他相關文章 | 邏輯需使用到評分，算法自訂即可 |
    | :------------------- | :----------------------------------- | :----------------------------- |

5. Background Worker   
    1. 設計一個背景 worker，每隔 5 分鐘檢查失敗的 metadata 爬取紀錄並重試   
    2. 最多重試三次   
   
   
### 補充說明   
1. 可以根據需求新增或調整 API   
2. 需考慮未來擴展性，請盡量在實作中區分不同層級的邏輯責任（如資料存取、業務邏輯、HTTP 請求處理等）。不限制使用哪種架構風格（如 MVC、DDD、Hexagonal 等），但需注意是否有良好模組劃分與程式結構清晰度   
   
   
## 交付內容   
1. Git repo 或 zip   
2. README.md 需包含以下內容   
    1. 系統說明與設計概念   
    2. 如何 Setup   
3. API 文件


### 1. 系統說明與設計概念
DeeliAI 是一個基於 Go 語言與 PostgreSQL 打造的後端服務，專為使用者提供文章收藏與個人化推薦功能。本專案的設計核心在於高內聚、低耦合，將各個功能模組化，確保程式碼的易讀性、可維護性與可擴展性。

#### 設計架構
##### Web Server 
採用經典的三層式架構 (Three-Tier Architecture)，每個層級負責不同的職責：
- Handler (API 層): 處理 HTTP 請求與回應，負責資料驗證（binding）和呼叫服務層的業務邏輯。所有 API 都透過 JWT 進行身份驗證。
- Service (服務層): 實現核心業務邏輯。它協調不同 Repository 的操作，並處理資料轉換、驗證和複雜的演算法，例如文章推薦。
- Repository (資料存取層): 負責與資料庫的互動。每個 Repository 都封裝了對單一資料表的 CRUD 操作，將資料庫細節與服務層隔離。

##### Worker
使用 Goroutine & Channel 實作一個高效的 Worker Pool，專門處理耗時的爬取任務，確保 API 服務的響應速度。

##### Scheduler 
負責執行定時任務，例如 Metadata 爬取失敗的重試機制。

##### API 文件
使用 swaggo/gin-swagger 處理 API 文件。


### 2. 如何 Setup
#### 環境要求
- Go: 1.23 或以上版本
- Docker: 用於啟動 PostgreSQL 服務
- migrate: Go 語言的資料庫 Migration 工具
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
- swag: 自動生成 Swagger API 文件的工具
```
go install github.com/swaggo/swag/cmd/swag@latest
```

#### 設定步驟
1. 複製專案
```
git clone https://github.com/your-username/deeli-ai.git
cd deeli-ai
```

2. 啟動 PostgreSQL 資料庫
```
docker-compose up -d postgres
```

3. 執行資料庫 Migration
```
migrate -path migrations -database "postgres://user:password@localhost:5432/deeliai?sslmode=disable" up
```
或是利用 PostgreSQL GUI 或其他工具***依序***手動執行 migrations 資料夾下的 *.up.sql

4. 運行專案
```
go run ./cmd/server/main.go
```

### 3. API 文件
在 server ***啟動後***，在瀏覽器中打開 http://localhost:8080/swagger/index.html，即可瀏覽完整的 API 文件並進行測試。

如果有 route 有更新，可運行 swag 工具來重新生成 API 文件
```
swag init -g ./cmd/server/main.go -o ./docs
```
   