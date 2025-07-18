@echo off
echo 🧹 Limpiando entorno...

REM Detener contenedores (DB + API Gateway + RabbitMQ + Redis)
echo 📦 Deteniendo contenedores Docker...
docker-compose -f dev-full-compose.yml down 2>nul
docker-compose -f dev-postgres-compose.yml down 2>nul
docker-compose -f rabbitmq-compose.yml down 2>nul
docker-compose -f redis-compose.yml down 2>nul

REM Limpiar binarios (agregando organization-svc)
echo 🗑️ Limpiando binarios compilados...
if exist auth-identity-svc\bin rmdir /s /q auth-identity-svc\bin 2>nul
if exist address-svc\bin rmdir /s /q address-svc\bin 2>nul
if exist person-svc\bin rmdir /s /q person-svc\bin 2>nul
if exist organization-svc\bin rmdir /s /q organization-svc\bin 2>nul

REM Limpiar archivos temporales de Air (agregando organization-svc)
echo 🗑️ Limpiando archivos temporales...
if exist auth-identity-svc\tmp rmdir /s /q auth-identity-svc\tmp 2>nul
if exist address-svc\tmp rmdir /s /q address-svc\tmp 2>nul
if exist person-svc\tmp rmdir /s /q person-svc\tmp 2>nul
if exist organization-svc\tmp rmdir /s /q organization-svc\tmp 2>nul

REM Limpiar logs (agregando organization-svc)
echo 🗑️ Limpiando logs...
if exist auth-identity-svc\logs rmdir /s /q auth-identity-svc\logs 2>nul
if exist address-svc\logs rmdir /s /q address-svc\logs 2>nul
if exist person-svc\logs rmdir /s /q person-svc\logs 2>nul
if exist organization-svc\logs rmdir /s /q organization-svc\logs 2>nul

echo ✅ Limpieza completada
pause