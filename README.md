# 🔗 URL-4-U

> **URL-4-U** es un acortador de URLs moderno, rápido y ligero desarrollado en **Go** y **React + Vite**, con soporte para códigos QR instantáneos y un panel de analíticas en tiempo real.

---

## 🚀 Características Principales

- ✂️ **Acortador Inteligente**: Generación de identificadores cortos únicos en Base62 con soporte para alias personalizados.
- 📱 **Generador de Códigos QR**: Exportación inmediata de códigos QR listos para escanear y descargar.
- 📊 **Panel de Analíticas**: Métricas de clics en tiempo real, desglose temporal, fuentes de referencia (referrers), navegadores y tipos de dispositivo.
- ⚡ **Alto Rendimiento en Go**: Arquitectura limpia con backend concurrente y eficiente en recursos.
- 🎨 **UI Moderna y Reactiva**: Frontend desarrollado con React, Vite y Tailwind CSS.
- 🐳 **Docker-Ready**: Puesta en marcha en un solo comando sin necesidad de instalar dependencias de sistema.

---

## 🛠️ Stack Tecnológico

| Capa | Tecnología | Descripción |
| :--- | :--- | :--- |
| **Backend** | [Go (Golang) 1.22+](https://go.dev/) | API REST construida siguiendo el Standard Go Project Layout |
| **Frontend** | [React 18](https://react.dev/) + [Vite](https://vitejs.dev/) | Single Page Application rápida y optimizada |
| **Estilos** | [Tailwind CSS](https://tailwindcss.com/) | Diseño responsivo, minimalista y modo oscuro |
| **Base de Datos** | [SQLite](https://sqlite.org/) | Base de datos embebida zero-config (sin servidores externos) |
| **Contenedores**| [Docker](https://www.docker.com/) & Docker Compose | Entorno reproducible y aislado para desarrollo y demo |

---

## 🏛️ Arquitectura del Proyecto

El backend sigue las mejores prácticas del **Standard Go Layout** y principios de arquitectura en capas (Clean Architecture):

```
url-4-u/
├── backend/
│   ├── cmd/
│   │   └── api/                  # Punto de entrada de la aplicación (main.go)
│   ├── internal/                 # Código privado de la aplicación
│   │   ├── config/               # Configuración y variables de entorno
│   │   ├── handler/              # Controladores HTTP (REST Handlers)
│   │   ├── model/                # Entidades de dominio y DTOs
│   │   ├── repository/           # Acceso a base de datos (SQLite/SQL)
│   │   └── service/              # Lógica de negocio (Hashing Base62, QR, etc.)
│   ├── pkg/                      # Paquetes utilitarios reutilizables
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/           # Componentes UI reutilizables
│   │   ├── pages/                # Vistas principales (Home, Analytics, etc.)
│   │   ├── services/             # Clientes de API HTTP
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
├── .gitignore
└── README.md
```

---

## 📡 Especificación de la API REST

| Método | Endpoint | Descripción | Body / Parámetros |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Healthcheck del servicio | - |
| `POST` | `/api/v1/urls` | Crea una nueva URL corta | `{ "original_url": "...", "custom_alias": "..." }` |
| `GET` | `/r/:code` | Redirige a la URL original y registra métricas | Redirección HTTP `302 Found` |
| `GET` | `/api/v1/urls/:code` | Obtiene detalles de un enlace corto | Retorna metadatos del link |
| `GET` | `/api/v1/urls/:code/qr` | Genera y sirve el código QR (PNG/SVG) | Retorna imagen del QR |
| `GET` | `/api/v1/urls/:code/analytics` | Obtiene analíticas y métricas de clics | Total de visitas, dispositivos, referrers, timeline |
| `GET` | `/api/v1/urls` | Lista de enlaces creados recientemente | Parámetros de paginación (`limit`, `offset`) |

---

## 💾 Estrategia de Base de Datos (Zero-Config)

Para facilitar la evaluación por parte de recruiters y evitar la instalación obligatoria de motores pesados como PostgreSQL o MySQL:

- Se utiliza **SQLite** (a través de un driver puro Go como `modernc.org/sqlite` o `mattn/go-sqlite3`).
- **Ventaja**: No requiere instalar ni configurar ningún gestor de base de datos; la base se crea automáticamente en un archivo local (`url4u.db`).
- **Abstracción**: El repositorio de datos implementa interfaces en Go, permitiendo migrar a PostgreSQL en producción con solo cambiar la cadena de conexión.

---

## 🚦 Cómo Iniciar el Proyecto

### Opción 1: Con Docker Compose (Recomendada para pruebas rápidas)

```bash
docker compose up --build
```
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`

### Opción 2: Modo Local

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
