# Artela Invitation API 💌

Backend Service untuk aplikasi Undangan Pernikahan Digital Artela.
Dibangun dengan prinsip **Clean Architecture** menggunakan **Go (Golang)**, **Fiber**, dan **PostgreSQL**.

## 🚀 Tech Stack
- **Language:** Go (1.20+)
- **Framework:** Fiber v2 (Fast HTTP)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Architecture:** Clean Architecture (Entity, Repository, Service, Handler)

## 📂 Struktur Folder
```text
.
├── cmd/                # Entry point aplikasi
├── internal/
│   ├── config/         # Konfigurasi Database & Env
│   ├── entity/         # Struktur Data (Model DB & JSON)
│   ├── handler/        # Controller HTTP
│   ├── repository/     # Query Database (SQL)
│   └── service/        # Business Logic
├── .env                # Environment Variables (Tidak dicommit)
└── .gitignore
```

## 🛠️ Cara Menjalankan (Local)
1. Clone Repo
```Bash
git clone <repo_url>
cd artela-api
```

2. Setup Database <br>
Pastikan PostgreSQL sudah berjalan dan buat database:
```SQL
CREATE DATABASE "artela-db";
```
3. Setup Environment <br>
Copy .env.example ke .env dan isi kredensial DB:
```Code snippet
APP_PORT=3000
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=artela-db
```
4. Run App
```Bash
go mod tidy
go run cmd/main.go
```

## 🌐 API Endpoints
Method|Endpoint|Deskripsi
---|---|---
GET|/health|Cek status server & DB
GET|/api/invitation/:slug|Ambil data undangan lengkap
POST|/api/admin/create|Buat undangan baru (Admin)

## 📦 Deployment (VPS Linux)
Karena dikembangkan di Mac/Windows tapi deploy di Linux, gunakan Cross Compile:
1. Build Binary
```Bash
GOOS=linux GOARCH=amd64 go build -o artela-api cmd/main.go
```
2. Upload file artela-api dan .env ke VPS.
3. Jalankan dengan PM2
```Bash
chmod +x artela-api
pm2 start ./artela-api --name "artela-backend"
```
