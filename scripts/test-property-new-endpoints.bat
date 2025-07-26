@echo off
echo Probando nuevos endpoints del Property Service con Postman...

REM Verificar que newman (Postman CLI) esté instalado
where newman >nul 2>&1
if errorlevel 1 (
    echo.
    echo ❌ Newman (Postman CLI) no está instalado.
    echo.
    echo Para instalarlo:
    echo   npm install -g newman
    echo.
    echo Alternativamente, abre Postman GUI y carga la colección manualmente.
    pause
    exit /b 1
)

echo.
echo ═══════════════════════════════════════════════════════════════
echo PROBANDO NUEVOS ENDPOINTS DEL PROPERTY SERVICE
echo ═══════════════════════════════════════════════════════════════

REM Verificar que el Property Service esté corriendo
curl -s http://localhost:4004/health >nul 2>&1
if errorlevel 1 (
    echo.
    echo ❌ Property Service no está corriendo en http://localhost:4004
    echo.
    echo Para iniciarlo:
    echo   cd property-svc
    echo   air
    echo.
    echo O ejecuta el entorno completo:
    echo   .\scripts\dev.bat
    echo.
    pause
    exit /b 1
)

echo ✅ Property Service está ejecutándose

echo.
echo 🧪 Ejecutando tests de los nuevos endpoints...

REM Crear directorio para reportes si no existe
if not exist "tmp" mkdir "tmp"

REM Ejecutar tests por sección individual para mejor seguimiento
echo.
echo 📋 Probando Property Listings...
newman run "REM-API-Collection.postman_collection.json" ^
    --environment "REM-Development.postman_environment.json" ^
    --folder "📋 Property Listings" ^
    --reporters cli,html ^
    --reporter-html-export "tmp/property-listings-test-report.html"

echo.
echo 📷 Probando Property Media...
newman run "REM-API-Collection.postman_collection.json" ^
    --environment "REM-Development.postman_environment.json" ^
    --folder "📷 Property Media" ^
    --reporters cli,html ^
    --reporter-html-export "tmp/property-media-test-report.html"

echo.
echo 💰 Probando Property Valuations...
newman run "REM-API-Collection.postman_collection.json" ^
    --environment "REM-Development.postman_environment.json" ^
    --folder "💰 Property Valuations" ^
    --reporters cli,html ^
    --reporter-html-export "tmp/property-valuations-test-report.html"

echo.
echo 🏢 Probando suite completa de Properties...
newman run "REM-API-Collection.postman_collection.json" ^
    --environment "REM-Development.postman_environment.json" ^
    --folder "🏢 Properties & Real Estate" ^
    --reporters cli,html ^
    --reporter-html-export "tmp/complete-property-test-report.html"

if errorlevel 1 (
    echo.
    echo ⚠️  Algunos tests fallaron. Revisa los reportes HTML generados en tmp/
) else (
    echo.
    echo ✅ Todos los tests pasaron exitosamente!
)

echo.
echo ═══════════════════════════════════════════════════════════════
echo TESTS EJECUTADOS PARA:
echo ═══════════════════════════════════════════════════════════════
echo.
echo 📋 PROPERTY LISTINGS:
echo   • GET /api/properties/listings - Listar todos los listings
echo   • POST /api/properties/listings - Crear nuevo listing
echo   • GET /api/properties/listings/{id} - Obtener listing específico
echo   • PUT /api/properties/listings/{id} - Actualizar listing
echo   • DELETE /api/properties/listings/{id} - Eliminar listing
echo   • GET /api/properties/listings/search - Buscar listings con filtros
echo   • GET /api/properties/{id}/listings - Listings por propiedad
echo.
echo 📸 PROPERTY MEDIA:
echo   • GET /api/properties/media - Listar medios de propiedad
echo   • POST /api/properties/media - Subir nuevo medio
echo   • GET /api/properties/media/{id} - Obtener medio específico
echo   • PUT /api/properties/media/{id} - Actualizar medio
echo   • DELETE /api/properties/media/{id} - Eliminar medio
echo   • PATCH /api/properties/media/{id}/main - Marcar como principal
echo   • GET /api/properties/{id}/media - Media por propiedad
echo.
echo 💰 PROPERTY VALUATIONS:
echo   • GET /api/properties/valuations - Listar valuaciones
echo   • POST /api/properties/valuations - Crear nueva valuación
echo   • GET /api/properties/valuations/{id} - Obtener valuación específica
echo   • PUT /api/properties/valuations/{id} - Actualizar valuación
echo   • DELETE /api/properties/valuations/{id} - Eliminar valuación
echo   • GET /api/properties/{id}/valuations - Valuaciones por propiedad
echo   • GET /api/properties/{id}/valuations/latest - Última valuación
echo.
echo 📊 REPORTES GENERADOS:
echo   • tmp/property-listings-test-report.html
echo   • tmp/property-media-test-report.html
echo   • tmp/property-valuations-test-report.html
echo   • tmp/complete-property-test-report.html
echo.
echo ═══════════════════════════════════════════════════════════════

echo.
echo Presiona cualquier tecla para abrir los reportes...
pause >nul

REM Abrir reportes principales en el navegador
echo Abriendo reportes...
start tmp/complete-property-test-report.html
timeout /t 2 >nul
start tmp/property-listings-test-report.html
timeout /t 1 >nul  
start tmp/property-media-test-report.html
timeout /t 1 >nul
start tmp/property-valuations-test-report.html
