# 🚀 Guía de Onboarding Rápido - REM Platform

## 👋 ¡Bienvenido al equipo!

Esta guía te ayudará a tener todo funcionando en **menos de 10 minutos**.

---

## ⚡ **Quick Start (5 minutos)**

### **1. Prerrequisitos** *(2 min)*
```bash
# Verificar que tienes todo instalado
docker --version     # ✅ Debe funcionar
go version           # ✅ Debe mostrar Go 1.21+
git --version        # ✅ Debe funcionar
```

### **2. Iniciar Desarrollo** *(2 min)*
```bash
# Desde la carpeta services/
.\scripts\dev.bat
```
**¿Qué hace?** 
- 🐳 Levanta infraestructura (PostgreSQL + API Gateway)
- 🔄 Ejecuta migraciones automáticamente
- 📊 **Te pregunta si poblar con datos de prueba** → Responde **"S"**
- 🔥 Inicia hot-reload en 3 servicios

### **3. Verificar** *(1 min)*
```bash
# En tu navegador
http://localhost:8081/health
```
**Debe responder**: `{"status":"ok","service":"api-gateway","version":"1.0.0"}`

---

## 📡 **Testing con Postman (3 minutos)**

### **1. Importar Collections**
- Abrir **Postman**
- **Import** → Arrastrar estos archivos:
  - `REM-API-Collection.postman_collection.json`
  - `REM-Development.postman_environment.json`
- **Seleccionar environment**: "REM Development Environment"

### **2. Primer Test**
- **Health Checks** → **Gateway Health** → Send ✅
- **Authentication** → **Login Juan** → Send ✅
- **Persons** → **List All Persons** → Send ✅

**¡Si todo responde correctamente, ya tienes todo funcionando!** 🎉

---

## 🏗️ **Arquitectura en 2 minutos**

```
Cliente/Frontend → API Gateway (nginx:8081) → Microservicios
                                            ├── Auth (4002)
                                            ├── Person (4001)  
                                            └── Address (4000)
                                                    ↓
                                            PostgreSQL (BD compartida)
```

### **¿Por qué esta arquitectura?**
- 🎯 **Un solo punto de entrada** para el frontend
- 🔒 **Seguridad centralizada** en el gateway
- 🚀 **Desarrollo rápido** con BD compartida (por ahora)
- 📈 **Escalable** a BDs separadas en el futuro

---

## 🧩 **Servicios Principales**

| Servicio | Puerto | Responsabilidad | Ejemplo Endpoint |
|----------|--------|-----------------|------------------|
| **Auth** | 4002 | 🔐 Usuarios y login | `POST /login` |
| **Person** | 4001 | 👥 Personas y empresas | `GET /persons` |
| **Address** | 4000 | 🏠 Direcciones | `POST /addresses` |

### **💡 Tip**: Siempre usa el API Gateway
```bash
# ✅ CORRECTO - A través del gateway
curl http://localhost:8081/api/persons/

# ❌ EVITAR - Acceso directo (solo para debug)
curl http://localhost:4001/persons/
```

---

## 📚 **Código que Debes Conocer**

### **rem-common/** - Librería Compartida
```go
// Importar utilidades comunes
import (
    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/logger"
    "github.com/rem-gestion/rem-common/middleware"
)
```

### **Estructura de un Servicio**
```
person-svc/
├── cmd/api/main.go           # 🚀 Punto de entrada
├── cmd/migrate/main.go       # 🔄 Migraciones
├── src/
│   ├── router/router.go      # 🌐 Rutas API
│   ├── controllers/          # 🎮 Lógica de endpoints
│   ├── models/              # 🗃️ Estructuras de DB
│   ├── services/            # 💼 Lógica de negocio
│   └── repository/          # 🔍 Acceso a datos
└── migrations/pg/           # 📊 Scripts SQL
```

---

## 🎯 **Flujo de Desarrollo Típico**

### **Agregar un nuevo endpoint**
1. **Model** → Definir estructura en `models/`
2. **Repository** → Funciones de DB en `repository/`
3. **Service** → Lógica de negocio en `services/`
4. **Controller** → Handler HTTP en `controllers/`
5. **Router** → Registrar ruta en `router/`
6. **Test** → Probar con Postman

### **Comandos útiles**
```bash
# Hot-reload automático (ya incluido en dev.bat)
air

# Limpiar todo y reiniciar
.\scripts\clean.bat
.\scripts\dev.bat

# Solo poblar datos
.\scripts\seed-db.bat

# Test rápido de endpoints
.\scripts\test-gateway.bat
```

---

## 🔧 **Troubleshooting Rápido**

### **Error: "Puerto ya en uso"**
```bash
.\scripts\clean.bat  # Limpia todo
# Esperar 30 segundos
.\scripts\dev.bat    # Reiniciar
```

### **Error: "No se puede conectar a la DB"**
- Verificar que Docker esté corriendo
- `docker ps` debe mostrar contenedor de PostgreSQL

### **Error: "502 Bad Gateway"**
- Los servicios no están corriendo
- Verificar que los 3 terminales de Air estén activos

### **Error en migraciones**
- Limpiar DB: `.\scripts\clean.bat`
- Reiniciar: `.\scripts\dev.bat`

---

## 📖 **Documentación Completa**

- 📋 **README.md** - Documentación principal completa
- 📁 **docs/** - Documentación específica por tema
- 📖 **service-name/README.md** - Docs de cada servicio
- 📡 **Postman** - Collections con ejemplos

---

## 👥 **Datos de Prueba Listos**

Una vez que ejecutes `.\scripts\dev.bat` y aceptes poblar la DB:

### **Usuarios para login**
```
juan.perez@example.com / password
maria.gonzalez@example.com / password  
admin@acmecorp.com / password
```

### **Qué encontrarás**
- 👤 8 usuarios con diferentes estados
- 🏢 5 personas individuales + 3 empresas
- 📍 10 direcciones en Argentina
- 📞 Contactos variados para cada persona

---

## 🎉 **¡Ya estás listo!**

Con esto tienes todo lo necesario para:
- ✅ **Desarrollar** nuevas funcionalidades
- ✅ **Testing** de endpoints existentes
- ✅ **Debugging** de issues
- ✅ **Contribuir** al proyecto

### **¿Dudas?**
1. Revisar **README.md** principal
2. Buscar en **docs/** específicos
3. Preguntar al equipo

---

*¡Bienvenido al futuro de la gestión inmobiliaria!* 🏠✨
