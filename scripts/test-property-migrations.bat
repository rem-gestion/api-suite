@echo off
echo Ejecutando migraciones y pruebas del Property Service...

REM Cambiar al directorio del property service
cd /d "%~dp0..\property-svc"

echo.
echo ═══════════════════════════════════════════════════════════════
echo EJECUTANDO MIGRACIONES DEL PROPERTY SERVICE
echo ═══════════════════════════════════════════════════════════════

REM Verificar que el contenedor de PostgreSQL esté ejecutándose
docker ps | findstr rem-postgres-dev >nul
if errorlevel 1 (
    echo.
    echo ❌ PostgreSQL no está ejecutándose.
    echo Ejecuta primero: .\scripts\dev.bat
    echo O inicia solo la base de datos: docker-compose -f dev-full-compose.yml up -d postgres
    echo.
    pause
    exit /b 1
)

echo.
echo 🔄 Ejecutando migraciones del Property Service...
go run ./cmd/migrate

if errorlevel 1 (
    echo.
    echo ❌ Error ejecutando las migraciones.
    echo Verifica que la base de datos esté disponible y las migraciones sean válidas.
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ Migraciones ejecutadas exitosamente!

echo.
echo ═══════════════════════════════════════════════════════════════
echo VERIFICANDO NUEVAS TABLAS
echo ═══════════════════════════════════════════════════════════════

echo.
echo 🔍 Verificando que las nuevas tablas fueron creadas...

REM Verificar las nuevas tablas
docker exec rem-postgres-dev psql -U user -d rem_development -c "\dt+ property_listings property_media property_valuations"

if errorlevel 1 (
    echo.
    echo ⚠️  No se pudieron verificar las tablas. Las migraciones pueden haber fallado.
) else (
    echo.
    echo ✅ Nuevas tablas verificadas exitosamente!
)

echo.
echo ═══════════════════════════════════════════════════════════════
echo COMPILANDO Y PROBANDO EL SERVICIO
echo ═══════════════════════════════════════════════════════════════

echo.
echo 🔨 Compilando el Property Service...
go build -o property-api.exe ./cmd/api

if errorlevel 1 (
    echo.
    echo ❌ Error compilando el servicio.
    echo Revisa los errores de compilación arriba.
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ Servicio compilado exitosamente!

echo.
echo 🧪 Ejecutando pruebas rápidas...
go test ./src/... -short

echo.
echo ═══════════════════════════════════════════════════════════════
echo RESUMEN DE NUEVAS FUNCIONALIDADES
echo ═══════════════════════════════════════════════════════════════
echo.
echo ✅ NUEVAS TABLAS CREADAS:
echo   • property_listings - Listings de propiedades (venta/alquiler)
echo   • property_media - Fotos, videos, planos de propiedades  
echo   • property_valuations - Valuaciones y tasaciones
echo.
echo ✅ NUEVOS ENDPOINTS DISPONIBLES:
echo   • GET/POST /api/property-listings
echo   • GET/POST /api/property-media
echo   • GET/POST /api/property-valuations
echo.
echo ✅ MIGRACIONES COMPLETADAS:
echo   • 0002_add_property_listings.up.sql
echo   • 0003_add_property_media.up.sql  
echo   • 0004_add_property_valuations.up.sql
echo.
echo 🚀 PRÓXIMOS PASOS:
echo   • Ejecutar .\scripts\seed-db.bat para datos de prueba
echo   • Probar endpoints con Postman
echo   • Iniciar servicio: .\scripts\dev.bat
echo.
echo ═══════════════════════════════════════════════════════════════

echo.
echo Presiona cualquier tecla para continuar...
pause >nul
