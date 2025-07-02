## 🛠️ Helper de migraciones `RunMigrations` – Guía rápida

`RunMigrations`, definido en **`rem-common/db/migrate.go`**, aplica o revierte migraciones para **PostgreSQL, MySQL o MongoDB** desde cualquier microservicio.

---

### 1. Variables de entorno necesarias

```env
# Motor de base de datos que usa este servicio
REM_DB_DRIVER=postgres          #  postgres | mysql | mongo

# --- Postgres (solo si REM_DB_DRIVER=postgres) ---
REM_POSTGRES_HOST=localhost
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=usuario
REM_POSTGRES_PASSWORD=clave
REM_POSTGRES_DB=base
REM_POSTGRES_SSLMODE=disable    # disable | require

# --- MySQL (solo si REM_DB_DRIVER=mysql) ---
REM_MYSQL_HOST=localhost
REM_MYSQL_PORT=3306
REM_MYSQL_USER=usuario
REM_MYSQL_PASSWORD=clave
REM_MYSQL_DB=base
REM_MYSQL_CHARSET=utf8mb4

# --- MongoDB (solo si REM_DB_DRIVER=mongo) ---
REM_MONGO_URI=mongodb://user:pass@localhost:27017   # opcional
REM_MONGO_DB=nombre_db
```

### 2. Estructura de migraciones

```bash
services/tu-servicio/
├─ migrations/
│   ├─ 0001_initial.up.sql      # crea tablas o colecciones
│   ├─ 0001_initial.down.sql    # revierte lo anterior
│   ├─ 0002_add_index.up.sql
│   └─ 0002_add_index.down.sql
└─ cmd/
    └─ migrate/
        └─ main.go              # corre RunMigrations
```

Para MongoDB los scripts se nombran igual pero con extensión `.js`.

### 3. cmd/migrate/main.go

```go
package main

import (
    "log"
    "os"

    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/db"
)

func main() {
    cfg := config.Load()          // lee todas las REM_*
    dir := "up"                   // por defecto aplica todo
    if len(os.Args) > 1 {
        dir = os.Args[1]          // "up" | "down" | "steps -1"
    }
    if err := db.RunMigrations(cfg.DriverRelacional, cfg, "./migrations", dir); err != nil {
        log.Fatal(err)
    }
}
```

### 4. Comandos de uso

```bash
# Aplica todas las migraciones pendientes
go run ./cmd/migrate                 # equivale a "up"

# Revierte todas (¡cuidado!)
go run ./cmd/migrate down

# Retrocede exactamente una migración
go run ./cmd/migrate "steps -1"

# Avanza dos migraciones
go run ./cmd/migrate "steps 2"
```

En CI/CD agregá un paso antes de levantar el contenedor:

```bash
go run ./cmd/migrate
```

### 5. Qué hace internamente

1. Construye el DSN según `REM_DB_DRIVER` y tus variables.

2. Crea el driver golang-migrate correspondiente:
   - `postgres.WithInstance` → Postgres
   - `mysql.WithInstance` → MySQL
   - `mongodb.WithInstance` → MongoDB

3. Lee los scripts de `migrations/`.

4. Ejecuta Up, Down o Steps N y registra el progreso en `schema_migrations` (o colección `.migrations` en Mongo).

### 6. ¿Por qué usarlo?

- Histórico versionado y reversible del esquema.
- Soporta SQL y MongoDB con un solo helper.