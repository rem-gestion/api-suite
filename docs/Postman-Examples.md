# 📋 Ejemplos de Payloads para REM API

Este archivo contiene ejemplos de payloads JSON para usar en las peticiones de Postman.

## 🔐 Authentication

### Register Individual User
```json
{
  "email": "nuevo.usuario@example.com",
  "password": "password123",
  "person": {
    "type": "individual",
    "first_name": "Nuevo",
    "last_name": "Usuario",
    "dni": "98765432",
    "sexo": "masculino",
    "address": {
      "street": "Av. Corrientes",
      "number": 1500,
      "city": "Buenos Aires",
      "state": "CABA",
      "zip": "1043",
      "country": "AR"
    }
  }
}
```

### Register Company User
```json
{
  "email": "admin@nuevaempresa.com",
  "password": "password123",
  "person": {
    "type": "company",
    "legal_name": "Nueva Empresa SA",
    "cuit": "30123456789",
    "society_type": "S.A.",
    "address": {
      "street": "Florida",
      "number": 1000,
      "floor": "12",
      "unit": "A",
      "city": "Buenos Aires",
      "state": "CABA",
      "zip": "1005",
      "country": "AR"
    }
  }
}
```

### Login
```json
{
  "email": "juan.perez@example.com",
  "password": "password"
}
```

## 👥 Persons

### Create Individual Person (Complete)
```json
{
  "type": "individual",
  "first_name": "Roberto",
  "last_name": "García",
  "dni": "87654321",
  "sexo": "masculino",
  "avatar_url": "https://avatar.example.com/roberto.jpg",
  "address": {
    "street": "Av. Rivadavia",
    "number": 5000,
    "floor": "3",
    "unit": "B",
    "city": "Buenos Aires",
    "state": "CABA",
    "zip": "1424",
    "country": "AR"
  },
  "contacts": [
    {
      "type": "email",
      "value": "roberto.garcia@example.com",
      "is_primary": true
    },
    {
      "type": "phone",
      "value": "+54911234567",
      "is_primary": true
    },
    {
      "type": "phone",
      "value": "+541144556677",
      "is_primary": false,
      "notes": "Teléfono del trabajo"
    }
  ]
}
```

### Create Company (Complete)
```json
{
  "type": "company",
  "legal_name": "Innovación Tech SRL",
  "cuit": "30876543210",
  "society_type": "S.R.L.",
  "address": {
    "street": "Av. Libertador",
    "number": 2500,
    "floor": "10",
    "unit": "A",
    "city": "Buenos Aires",
    "state": "CABA",
    "zip": "1425",
    "country": "AR"
  },
  "contacts": [
    {
      "type": "email",
      "value": "info@innovaciontech.com",
      "is_primary": true
    },
    {
      "type": "email",
      "value": "ventas@innovaciontech.com",
      "is_primary": false
    },
    {
      "type": "phone",
      "value": "+541144556677",
      "is_primary": true
    },
    {
      "type": "website",
      "value": "https://www.innovaciontech.com",
      "is_primary": false
    }
  ]
}
```

### Update Person
```json
{
  "first_name": "Roberto Carlos",
  "last_name": "García Pérez",
  "avatar_url": "https://avatar.example.com/roberto-new.jpg"
}
```

### Bulk Create Persons
```json
{
  "persons": [
    {
      "type": "individual",
      "first_name": "Pedro",
      "last_name": "López",
      "dni": "11111111",
      "sexo": "masculino"
    },
    {
      "type": "individual", 
      "first_name": "Carla",
      "last_name": "Sánchez",
      "dni": "22222222",
      "sexo": "femenino"
    },
    {
      "type": "company",
      "legal_name": "Startup ABC SAS",
      "cuit": "30999888777",
      "society_type": "S.A.S."
    }
  ]
}
```

## 📱 Contacts

### Add Email Contact
```json
{
  "type": "email",
  "value": "nuevo.email@example.com",
  "is_primary": false,
  "notes": "Email secundario"
}
```

### Add Phone Contact
```json
{
  "type": "phone",
  "value": "+54911999888",
  "is_primary": false,
  "notes": "Teléfono personal"
}
```

### Add Website Contact
```json
{
  "type": "website",
  "value": "https://www.ejemplo.com",
  "is_primary": false,
  "notes": "Sitio web corporativo"
}
```

### Update Contact
```json
{
  "value": "+54911888999",
  "is_primary": true,
  "notes": "Teléfono principal actualizado"
}
```

## 🏠 Addresses

### Create Address (Departamento)
```json
{
  "street": "Av. Libertador",
  "number": 2500,
  "floor": "15",
  "unit": "B",
  "city": "Buenos Aires",
  "state": "CABA", 
  "zip": "1425",
  "country": "AR"
}
```

### Create Address (Casa)
```json
{
  "street": "San Martín",
  "number": 1234,
  "city": "San Isidro",
  "state": "Buenos Aires",
  "zip": "1642",
  "country": "AR"
}
```

### Create Address (Oficina)
```json
{
  "street": "Florida",
  "number": 850,
  "floor": "12",
  "unit": "A",
  "city": "Buenos Aires", 
  "state": "CABA",
  "zip": "1005",
  "country": "AR"
}
```

### Update Address
```json
{
  "street": "Av. Libertador",
  "number": 2500,
  "floor": "16", 
  "unit": "A",
  "city": "Buenos Aires",
  "state": "CABA",
  "zip": "1425",
  "country": "AR"
}
```

## 🔧 Query Parameters Examples

### List Persons with Filters
```
GET /api/persons?page=1&per_page=10&type=individual&with_contacts=true&search=Juan
```

### List Users with Pagination
```
GET /api/users?page=1&per_page=20
```

## 📝 Headers Examples

### Authenticated Request
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

### Public Request
```
Content-Type: application/json
```

## 🔍 IDs de Prueba (Seed Data)

Estos IDs están disponibles en los datos de seed para pruebas:

### Person IDs
- `11111111-1111-1111-1111-111111111111` - Juan Pérez
- `22222222-2222-2222-2222-222222222222` - María González  
- `66666666-6666-6666-6666-666666666666` - ACME Corporation

### Address IDs
- `11111111-1111-1111-1111-111111111111` - Av. Corrientes 1234
- `22222222-2222-2222-2222-222222222222` - Av. Santa Fe 2567
- `88888888-8888-8888-8888-888888888888` - Av. Leandro N. Alem 456

### User IDs
- `11111111-1111-1111-1111-111111111111` - juan.perez@example.com
- `22222222-2222-2222-2222-222222222222` - maria.gonzalez@example.com
- `66666666-6666-6666-6666-666666666666` - admin@acmecorp.com

## 🚀 Variables de Postman

En tus requests puedes usar estas variables:

- `{{base_url}}` - http://localhost:8081
- `{{access_token}}` - Se llena automáticamente al hacer login
- `{{user_id}}` - Se llena automáticamente al hacer login
- `{{person_id}}` - Se llena automáticamente al listar/crear personas
- `{{address_id}}` - Se llena automáticamente al crear direcciones
- `{{contact_id}}` - Se llena automáticamente al listar/crear contactos

---

💡 **Tip**: Copia y pega estos payloads en el body de tus requests en Postman para pruebas rápidas.
