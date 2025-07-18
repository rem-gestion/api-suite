@echo off
echo.
echo 🧪 TESTING API GATEWAY - REM Platform
echo ═══════════════════════════════════════════════════════════════
echo.

echo 📡 Testing API Gateway Health:
curl -s http://localhost:8081/health
echo.
echo.

echo 🔐 Testing Auth Service via Gateway:
curl -s http://localhost:8081/api/health/auth
echo.
echo.

echo 👤 Testing Person Service via Gateway:  
curl -s http://localhost:8081/api/health/person
echo.
echo.

echo 🏠 Testing Address Service via Gateway:
curl -s -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/health/address
echo.
echo.

echo 🏢 Testing Property Service via Gateway:
curl -s http://localhost:8081/api/properties/health
echo.
echo.

echo ═══════════════════════════════════════════════════════════════
echo 🧪 Testing Real API Endpoints...
echo ═══════════════════════════════════════════════════════════════
echo.

echo 📋 Testing List Persons (should show persons if populated):
curl -s -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/persons/
echo.
echo.

echo 📋 Testing List Addresses (should show addresses if populated):  
curl -s -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/addresses/
echo.
echo.

echo 📋 Testing List Properties (should show properties if populated):
curl -s -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/properties/
echo.
echo.

echo 📋 Testing List Amenities (should show amenities if populated):
curl -s -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/amenities/
echo.
echo.

echo ✅ INTERPRETACION DE RESULTADOS:
echo ════════════════════════════════════════════════════════════════
echo   📡 Gateway Health:    Debe devolver JSON con status "ok"
echo   🔐 Auth Health:       Debe devolver "ok" (texto plano)  
echo   👤 Person Health:     Debe devolver "ok" (texto plano)
echo   🏠 Address Health:    Debe devolver "ok" (texto plano)
echo   🏢 Property Health:   Debe devolver "ok" (texto plano)
echo   📋 Persons List:      Debe devolver JSON con array (si poblado)
echo   📋 Addresses List:    Debe devolver JSON con array (si poblado)
echo   📋 Properties List:   Debe devolver JSON con array (si poblado)
echo   📋 Amenities List:    Debe devolver JSON con array (si poblado)
echo.
echo 🚨 ERRORES COMUNES:
echo   • "502 Bad Gateway" = Servicio no está corriendo
echo   • "api key inválida" = Falta X-Api-Key header
echo   • Sin respuesta = Puerto incorrecto o servicio down
echo.
echo 📋 URLs disponibles:
echo    API Gateway:   http://localhost:8081
echo    Auth API:      http://localhost:8081/api/auth/
echo    Users API:     http://localhost:8081/api/users/
echo    Person API:    http://localhost:8081/api/persons/
echo    Address API:   http://localhost:8081/api/addresses/
echo.
echo 💡 TIPS:
echo   • Si algo falla, verifica que los 3 terminales estén corriendo
echo   • Para datos de prueba: .\scripts\seed-db.bat
echo   • Para limpiar: .\scripts\clean.bat
echo.
pause
