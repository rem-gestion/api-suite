REM filepath: scripts/prod.bat
@echo off
echo 🏭 Preparando entorno de PRODUCCIÓN...

REM Setear variables de entorno
set REM_ENVIRONMENT=production

REM Compilar servicios
echo 🔨 Compilando servicios...
cd auth-identity-svc
if not exist bin mkdir bin
go build -o bin/auth-svc.exe ./cmd/main.go
cd ..

cd address-svc
if not exist bin mkdir bin
go build -o bin/address-svc.exe ./cmd/api/main.go
cd ..

cd person-svc
if not exist bin mkdir bin
go build -o bin/person-svc.exe ./cmd/main.go
cd ..

REM Ejecutar migraciones de producción
echo 🔄 Ejecutando migraciones de producción...
cd auth-identity-svc
bin\auth-svc.exe migrate
cd ..

cd address-svc
bin\address-svc.exe migrate
cd ..

cd person-svc
bin\person-svc.exe migrate
cd ..

echo 🚀 Binarios compilados. Para iniciar en producción:
echo Terminal 1: cd auth-identity-svc ^&^& bin\auth-svc.exe
echo Terminal 2: cd address-svc ^&^& bin\address-svc.exe
echo Terminal 3: cd person-svc ^&^& bin\person-svc.exe
pause