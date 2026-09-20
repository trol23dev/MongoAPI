# MongoDB REST API in Go (Hexagonal Architecture + Fiber)

> API RESTful de ultra-alto rendimiento para MongoDB desarrollada en **Go (Golang)** bajo los principios de **Arquitectura Hexagonal (Ports & Adapters)**, utilizando el framework **Fiber v2**, con un **CRUD completo**, pool de conexiones persistente y soporte dinámico para múltiples bases de datos y colecciones.

---

## 👨‍💻 Autor

- [Ramón L. (trol23dev)](https://www.github.com/trol23dev)

---

## 🏛️ Arquitectura Hexagonal (Ports & Adapters)

El proyecto está diseñado siguiendo estrictamente la separación de responsabilidades:

```
├── cmd/
│   └── api/
│       └── main.go                 # Composition Root (Inyección de Dependencias)
├── internal/
│   ├── core/                       # Núcleo de Dominio & Casos de Uso
│   │   ├── domain/                 # Entidades y errores agnósticos
│   │   │   ├── document.go
│   │   │   ├── filter.go
│   │   │   └── errors.go
│   │   ├── ports/                  # Interfaces de Puertos
│   │   │   ├── inbound.go          # Driving Port: DocumentUseCase
│   │   │   └── outbound.go         # Driven Port: DocumentRepository
│   │   └── services/               # Implementación de casos de uso (Business Logic)
│   │       └── document_service.go
│   ├── adapters/                   # Adaptadores
│   │   ├── http/                   # Adaptador Primario (Driving): Fiber REST API
│   │   │   ├── handler.go
│   │   │   └── router.go
│   │   └── repository/             # Adaptador Secundario (Driven): MongoDB Driver v2
│   │       └── mongo_repository.go
│   └── infrastructure/             # Infraestructura
│       ├── config/                 # Carga de variables de entorno (.env)
│       ├── db/                     # Cliente y pool de conexión MongoDB
│       └── logger/                 # Logging con emojis y colores
├── bin/
│   └── api.exe                     # Binario compilado
├── .env
├── example.env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 Características

- **Arquitectura Hexagonal**: Máxima desacoplación entre el dominio, los adaptadores de transporte (HTTP) y la base de datos (MongoDB).
- **Framework Web Fiber v2 (`gofiber/fiber`)**: Construido sobre Fasthttp (el motor HTTP más rápido en Go), con middleware de recuperación, CORS y Monitor integrado.
- **Driver Oficial MongoDB v2**: Pool de conexiones persistente (`go.mongodb.org/mongo-driver/v2`) para máxima eficiencia.
- **CRUD Completo**:
  - `GET /api/:database/:collection` (Listar todos con soporte para `limit`, `skip`, `sort`, `filter`)
  - `GET /api/:database/:collection/:id` (Obtener documento por ID)
  - `POST /api/:database/:collection` (Crear documento individual o inserción por lotes)
  - `PUT /api/:database/:collection/:id` (Reemplazo / actualización completa)
  - `PATCH /api/:database/:collection/:id` (Actualización parcial `$set`)
  - `DELETE /api/:database/:collection/:id` (Eliminación por ID)
  - `GET /health` (Estado del servidor y conexión con MongoDB)
  - `GET /monitor` (Dashboard visual en tiempo real del servidor)
- **Logging Personalizado**: Salida con emojis y colores para depuración clara (`🔧 SYSTEM`, `⚠️ WARNING`, `✅ SUCCESS`, `❌ ERROR`).

---

## ⚙️ Variables de Entorno

Crea o edita tu archivo `.env` en la raíz del proyecto:

```env
PORT=3000
MONGO_URI="mongodb+srv://<usuario>:<password>@<cluster>.mongodb.net/?retryWrites=true&w=majority"
```

| Variable | Descripción | Valor por Defecto |
| :--- | :--- | :--- |
| `PORT` | Puerto donde escuchará la API | `3000` |
| `MONGO_URI` | Cadena de conexión de MongoDB | *(Requerido)* |

---

## 📦 Ejecución y Despliegue

### Modo Desarrollo
```bash
go run cmd/api/main.go
```

### Compilar Binario
```bash
go build -o bin/api.exe ./cmd/api
```

### Ejecutar Binario
```bash
./bin/api.exe
```

---

## 📖 Referencia de la API (API Reference)

### 1. Health Check & Monitor
```http
GET /health
```
**Respuesta:**
```json
{
  "status": "ok",
  "database": "connected",
  "time": "2026-08-15T13:50:00Z"
}
```

```http
GET /monitor
GET /metrics
```
> **Fiber Monitor Dashboard**: Abre `http://localhost:3000/monitor` en tu navegador para ver un panel visual en tiempo real con el consumo de CPU, RAM, uso de red, peticiones por segundo y estadísticas del servidor.

---

### 2. Obtener Todos los Documentos
```http
GET /api/:database/:collection
```

#### Parámetros Query Opcionales:
- `limit`: Número máximo de documentos a devolver (ej. `?limit=20`)
- `skip`: Número de documentos a omitir para paginación (ej. `?skip=40`)
- `sort`: Campo para ordenar. Anteponer `-` para orden descendente (ej. `?sort=-createdAt` o `?sort=name`)
- `filter`: Filtro JSON para consultas complejas (ej. `?filter={"status":"active"}`)

---

### 3. Obtener Documento por ID
```http
GET /api/:database/:collection/:id
```
- Soporta tanto `ObjectID` estándar de MongoDB de 24 caracteres como identificadores de cadena personalizados.

---

### 4. Crear Documento(s)
```http
POST /api/:database/:collection
```
- **Insertar un documento único**: Enviar un objeto JSON `{ "nombre": "Producto 1", "precio": 100 }`.
- **Insertar múltiples documentos (Bulk)**: Enviar un array JSON `[{"nombre": "A"}, {"nombre": "B"}]`.

---

### 5. Reemplazar Documento Completo (PUT)
```http
PUT /api/:database/:collection/:id
```
- Reemplaza todo el contenido del documento manteniendo el mismo `_id`.

---

### 6. Actualizar Parcialmente Documento (PATCH)
```http
PATCH /api/:database/:collection/:id
```
- Aplica una actualización `$set` únicamente sobre los campos incluidos en el cuerpo JSON:
```json
{
  "precio": 120,
  "stock": 15
}
```

---

### 7. Eliminar Documento (DELETE)
```http
DELETE /api/:database/:collection/:id
```
