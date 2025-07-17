
# 🏠 Property Service

Servicio de gestión de propiedades inmobiliarias para REM Platform.

## Descripción
Permite la administración de propiedades (casas, departamentos, terrenos, comerciales), sus características, historial y relación con personas y direcciones.

---

## Endpoints principales

| Método | Ruta                | Descripción                  |
|--------|---------------------|------------------------------|
| POST   | /properties         | Crear propiedad              |
| GET    | /properties         | Listar propiedades           |
| GET    | /properties/:id     | Obtener propiedad por ID     |
| PUT    | /properties/:id     | Actualizar propiedad         |
| DELETE | /properties/:id     | Eliminar propiedad           |

---

## Ejemplo de request (POST /properties)
```json
{
  "type": "casa",
  "address_id": 1,
  "owner_id": 1,
  "title": "Casa en el centro",
  "description": "Hermosa casa",
  "area": 120,
  "rooms": 4,
  "bathrooms": 2,
  "amenities": "jardin,pileta"
}
```

## Ejemplo de response (GET /properties/1)
```json
{
  "id": 1,
  "type": "casa",
  "address_id": 1,
  "owner_id": 1,
  "title": "Casa en el centro",
  "description": "Hermosa casa",
  "area": 120,
  "rooms": 4,
  "bathrooms": 2,
  "amenities": "jardin,pileta",
  "created_at": "2025-07-16T21:30:00Z",
  "updated_at": "2025-07-16T21:30:00Z"
}
```

---

## Validaciones
- Campos obligatorios: `type`, `title`, `address_id`, `owner_id`
- `area`, `rooms`, `bathrooms` no pueden ser negativos
- Mensajes de error claros en formato JSON

---

## Migraciones
- `migrations/pg/0001_initial.up.sql`: crea la tabla `properties`
- `migrations/pg/0001_initial.down.sql`: elimina la tabla

---

## Testing rápido (PowerShell)
```powershell
# Crear propiedad
Invoke-WebRequest -Uri http://localhost:4003/properties -Method POST -Body '{"type":"casa","address_id":1,"owner_id":1,"title":"Casa en el centro","description":"Hermosa casa","area":120,"rooms":4,"bathrooms":2,"amenities":"jardin,pileta"}' -ContentType "application/json"

# Listar propiedades
Invoke-WebRequest -Uri http://localhost:4003/properties -Method GET

# Obtener por ID
Invoke-WebRequest -Uri http://localhost:4003/properties/1 -Method GET

# Actualizar
Invoke-WebRequest -Uri http://localhost:4003/properties/1 -Method PUT -Body '{"title":"Casa remodelada"}' -ContentType "application/json"

# Eliminar
Invoke-WebRequest -Uri http://localhost:4003/properties/1 -Method DELETE
```

---

## Por hacer / Mejoras futuras
- Integración real con address-svc y person-svc
- Health check endpoint
- Tests automáticos
- Seed de datos
- Documentación OpenAPI

---

*Sigue la arquitectura y patrones de los otros microservicios REM.*
