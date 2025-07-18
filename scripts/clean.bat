@echo off
echo 🧹 Limpiando entorno...

REM Detener contenedores (DB + API Gateway)
echo 📦 Deteniendo contenedores Docker...
docker-compose -f dev-full-compose.yml down 2>nul
docker-compose -f dev-postgres-compose.yml down 2>nul

REM Limpiar binarios
echo 🗑️ Limpiando binarios compilados...
if exist auth-identity-svc\bin rmdir /s /q auth-identity-svc\bin 2>nul
if exist address-svc\bin rmdir /s /q address-svc\bin 2>nul
if exist person-svc\bin rmdir /s /q person-svc\bin 2>nul
if exist property-svc\bin rmdir /s /q property-svc\bin 2>nul

REM Limpiar archivos temporales de Air
echo 🗑️ Limpiando archivos temporales...
if exist auth-identity-svc\tmp rmdir /s /q auth-identity-svc\tmp 2>nul
if exist address-svc\tmp rmdir /s /q address-svc\tmp 2>nul
if exist person-svc\tmp rmdir /s /q person-svc\tmp 2>nul
if exist property-svc\tmp rmdir /s /q property-svc\tmp 2>nul

REM Limpiar logs
echo 🗑️ Limpiando logs...
if exist auth-identity-svc\logs rmdir /s /q auth-identity-svc\logs 2>nul
if exist address-svc\logs rmdir /s /q address-svc\logs 2>nul
if exist person-svc\logs rmdir /s /q person-svc\logs 2>nul
if exist property-svc\logs rmdir /s /q property-svc\logs 2>nul

echo ✅ Limpieza completada
pause