# 🔗 URL-4-U

> **URL-4-U** is a modern, fast, and lightweight URL shortener built with **Go** and **React + Vite**, featuring instant QR code generation and real-time analytics.
>
> 🌐 **Live Demo:** [https://url-4-u.vercel.app](https://url-4-u.vercel.app)

> [!NOTE]
> **Free Tier Note (Backend Cold Starts):** The backend API is hosted on Render's free tier. Free instances spin down during periods of inactivity. If the service is sleeping, the first request may take ~30-50 seconds to wake up; subsequent requests are immediate.

---

## 🚀 Key Features

- ✂️ **Smart Shortener**: Generates unique Base62 short IDs with custom alias support.
- 📱 **QR Code Generator**: Instant export of scan-ready and downloadable QR codes.
- 📊 **Analytics Dashboard**: Real-time click tracking, timeline trends, referrers, browsers, and device types.
- ⚡ **High Performance with Go**: Clean architecture with a concurrent, resource-efficient backend engine.
- 🎨 **Modern & Reactive UI**: Built with React, Vite, Tailwind CSS, and Lucide Icons.
- 🌐 **Automatic Internationalization (i18n)**: Auto-detects device/browser language (English / Spanish) with instant manual switching and local persistence (`i18next`).
- 🐳 **Docker-Ready**: Launch the entire stack with a single command without installing local dependencies.

---

## 🛠️ Tech Stack

| Layer | Technology | Description |
| :--- | :--- | :--- |
| **Backend** | [Go (Golang) 1.22+](https://go.dev/) | RESTful API adhering to the Standard Go Project Layout |
| **Frontend** | [React 18](https://react.dev/) + [Vite](https://vitejs.dev/) | Fast, optimized Single Page Application |
| **Styling** | [Tailwind CSS](https://tailwindcss.com/) | Responsive, minimalist dark mode design |
| **Database** | [SQLite](https://sqlite.org/) | Embedded zero-config database (no external servers required) |
| **Containers** | [Docker](https://www.docker.com/) & Docker Compose | Reproducible, isolated environment for development and deployment |

---

## 🏛️ Project Architecture

The backend follows **Standard Go Layout** best practices and Clean Architecture principles:

```
url-4-u/
├── backend/
│   ├── cmd/
│   │   └── api/                  # Application entrypoint (main.go)
│   ├── internal/                 # Private application code
│   │   ├── config/               # Configuration & environment variables
│   │   ├── handler/              # HTTP controllers (REST Handlers)
│   │   ├── model/                # Domain entities & DTOs
│   │   ├── repository/           # Database access (SQLite/SQL)
│   │   └── service/              # Business logic (Base62 hashing, QR, analytics)
│   ├── pkg/                      # Reusable utility packages
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/           # Reusable UI components
│   │   ├── i18n/                 # Internationalization configurations & locales
│   │   ├── services/             # HTTP API client services
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
├── .gitignore
└── README.md
```

---

## 📡 REST API Specification

| Method | Endpoint | Description | Body / Parameters |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Service health check | - |
| `POST` | `/api/v1/urls` | Creates a new shortened URL | `{ "original_url": "...", "custom_alias": "..." }` |
| `GET` | `/r/:code` | Redirects to original URL and tracks analytics | HTTP `302 Found` redirection |
| `GET` | `/api/v1/urls/:code` | Retrieves details for a short link | Returns link metadata |
| `GET` | `/api/v1/urls/:code/qr` | Generates and serves QR code (PNG) | Returns QR code image |
| `GET` | `/api/v1/urls/:code/analytics` | Retrieves click analytics and metrics | Total clicks, devices, referrers, timeline |
| `GET` | `/api/v1/urls` | Lists recently created links | Pagination query params (`limit`, `offset`) |

---

## 💾 Database Strategy (Zero-Config)

Designed for easy evaluation and zero setup friction:

- Uses **SQLite** (via pure Go SQLite driver).
- **Advantage**: No database server installation or configuration required; the database is automatically created in a local file (`url4u.db`).
- **Clean Abstraction**: The data repository implements Go interfaces, making it straightforward to switch to PostgreSQL or MySQL in production with a single connection string change.

---

## 🚦 Getting Started

### Option 1: Docker Compose (Recommended)

```bash
docker compose up --build
```
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`

### Option 2: Local Development

**Backend:**
```bash
cd backend
go run cmd/api/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

---

## 🌐 Free Cloud Deployment

### 1. Backend on [Render](https://render.com/) (Web Service)
- **Root Directory**: `backend`
- **Environment**: `Go`
- **Build Command**: `go build -o server ./cmd/api/main.go`
- **Start Command**: `./server`
- **Environment Variables**:
  - `PORT`: `10000` (or leave default assigned by Render)
  - `BASE_URL`: `https://your-service-name.onrender.com`
  - `ALLOWED_ORIGINS`: `*` (or your Vercel URL)

### 2. Frontend on [Vercel](https://vercel.com/)
- **Root Directory**: `frontend`
- **Framework Preset**: `Vite`
- **Environment Variables**:
  - `VITE_API_URL`: `https://your-service-name.onrender.com`
