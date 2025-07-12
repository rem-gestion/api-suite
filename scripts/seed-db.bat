@echo off
echo ═══════════════════════════════════════════════════════════════
echo SCRIPT DE POBLACION DE BASE DE DATOS - DESARROLLO
echo ═══════════════════════════════════════════════════════════════
echo.

REM Verificar que el contenedor de PostgreSQL esté corriendo
docker ps --format "table {{.Names}}\t{{.Status}}" | findstr rem-postgres-dev >nul
if errorlevel 1 (
    echo ❌ El contenedor rem-postgres-dev no está corriendo.
    echo    Ejecuta primero: docker-compose -f dev-full-compose.yml up -d
    echo    O usa: .\scripts\dev.bat
    pause
    exit /b 1
)

echo ✅ Contenedor PostgreSQL detectado
echo.

REM Verificar conectividad a la base de datos
echo 🔍 Verificando conectividad a la base de datos...
docker exec rem-postgres-dev pg_isready -U user -d rem_development >nul 2>&1
if errorlevel 1 (
    echo ❌ No se puede conectar a la base de datos rem_development
    echo    Verifica que el contenedor esté funcionando correctamente
    pause
    exit /b 1
)

echo ✅ Conectividad verificada
echo.

REM Mostrar opción de limpiar datos existentes
echo ⚠️  ATENCION: Este script insertará datos de prueba en la base de datos.
echo    Si ya existen datos, podrían generarse conflictos de claves duplicadas.
echo.
set /p CLEAN_CHOICE="¿Quieres limpiar los datos existentes antes de poblar? (s/N): "

if /i "%CLEAN_CHOICE%"=="s" (
    echo.
    echo 🧹 Limpiando datos existentes...
    
    REM Script para limpiar datos existentes respetando las FK
    echo BEGIN; > temp_clean.sql
    echo -- Limpiar en orden inverso a las dependencias >> temp_clean.sql
    echo DELETE FROM person_contacto; >> temp_clean.sql
    echo DELETE FROM auth_users; >> temp_clean.sql
    echo DELETE FROM person_individual; >> temp_clean.sql
    echo DELETE FROM person_company; >> temp_clean.sql
    echo DELETE FROM person_person; >> temp_clean.sql
    echo DELETE FROM auth_accounts; >> temp_clean.sql
    echo DELETE FROM address_addresses; >> temp_clean.sql
    echo COMMIT; >> temp_clean.sql
    
    docker exec -i rem-postgres-dev psql -U user -d rem_development < temp_clean.sql
    del temp_clean.sql
    
    if errorlevel 1 (
        echo ❌ Error limpiando datos existentes
        pause
        exit /b 1
    )
    
    echo ✅ Datos existentes limpiados
    echo.
)

REM Ejecutar script de población
echo 🌱 Poblando base de datos con datos de prueba...
type "%~dp0seed-dev-db.sql" | docker exec -i rem-postgres-dev psql -U user -d rem_development

if errorlevel 1 (
    echo.
    echo ❌ Error poblando la base de datos
    echo    Revisa los logs arriba para más detalles
    echo.
    pause
    exit /b 1
) else (
    echo.
    echo ✅ ¡Base de datos poblada exitosamente!
    echo.
    echo ═══════════════════════════════════════════════════════════════
    echo USUARIOS DE PRUEBA DISPONIBLES:
    echo ═══════════════════════════════════════════════════════════════
    echo.
    echo 👤 INDIVIDUALES:
    echo   • juan.perez@example.com (password: password) - ACTIVO
    echo   • maria.gonzalez@example.com (password: password) - ACTIVO
    echo   • carlos.rodriguez@example.com (password: password) - PENDIENTE
    echo   • ana.martinez@example.com (password: password) - ACTIVO
    echo   • laura.fernandez@gmail.com (Google OAuth) - ACTIVO
    echo.
    echo 🏢 EMPRESAS:
    echo   • admin@acmecorp.com (password: password) - ACTIVO
    echo   • contacto@techsolutions.com (password: password) - ACTIVO
    echo   • ventas@innovatech.com.ar (password: password) - ACTIVO
    echo.
    echo ═══════════════════════════════════════════════════════════════
    echo DATOS INSERTADOS:
    echo   📍 10 direcciones en Argentina
    echo   🔐 8 cuentas de autenticación
    echo   👥 8 personas (5 individuales + 3 empresas)
    echo   📞 20+ contactos (emails, teléfonos, WhatsApp)
    echo.
    echo 🌐 Prueba los endpoints a través del API Gateway:
    echo   http://localhost:8081/api/persons/
    echo   http://localhost:8081/api/addresses/
    echo   http://localhost:8081/api/auth/
    echo.
)

echo Presiona cualquier tecla para continuar...
pause >nul
