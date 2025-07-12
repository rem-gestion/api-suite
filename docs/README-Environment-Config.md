# 🔧 Configuración de Variables de Entorno

## 📁 Estructura de archivos .env

```
services/
├── .env.development          # ✅ Desarrollo (se commitea)
├── address-svc/
│   ├── .env.development      # ✅ Desarrollo (se commitea)
│   ├── .env.production.example # ✅ Plantilla (se commitea)
│   └── .env.production       # ❌ Producción real (NO se commitea)
├── person-svc/
│   ├── .env.development      # ✅ Desarrollo (se commitea)
│   ├── .env.production.example # ✅ Plantilla (se commitea)
│   └── .env.production       # ❌ Producción real (NO se commitea)
└── auth-identity-svc/
    ├── .env.development      # ✅ Desarrollo (se commitea)
    ├── .env.production.example # ✅ Plantilla (se commitea)
    └── .env.production       # ❌ Producción real (NO se commitea)
```

## 🔍 Prioridad de carga

El sistema carga variables de entorno en este orden (se detiene en el primero que encuentra):

1. `.env.{environment}.local` (overrides locales)
2. `.env.{environment}` (configuración por entorno)
3. `.env.local` (overrides globales locales)
4. `.env` (configuración global)

## 🛠️ Configuración por Entorno

### 💻 Desarrollo
- **Archivos**: `.env.development`
- **Base de datos**: Una sola BD compartida (`rem_development`)
- **Credenciales**: Simples y conocidas
- **API Keys**: De desarrollo (no secretas)

### 🚀 Producción
- **Archivos**: `.env.production` (crear desde `.env.production.example`)
- **Base de datos**: BD individual por servicio
- **Credenciales**: Seguras y únicas
- **API Keys**: Producción (secretas)

## 📋 Para configurar producción:

1. Copia las plantillas:
   ```bash
   cp .env.production.example .env.production
   ```

2. Edita `.env.production` con valores reales:
   - Cambia todas las credenciales
   - Usa bases de datos separadas
   - Configura SSL/TLS
   - Usa API keys seguras

## 🔒 Seguridad

- ✅ **Se commitean**: `.env.development`, `.env.*.example`
- ❌ **NO se commitean**: `.env.production`, `.env.local`, `.env`
- 🔑 **Secretos reales**: Solo en archivos excluidos del git
