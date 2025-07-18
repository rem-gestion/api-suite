@echo off
echo Iniciando entorno de DESARROLLO...

REM Setear variables de entorno
set REM_ENVIRONMENT=development
set REM_POSTGRES_DB=rem_development

REM Verificar Docker
docker info >nul 2>&1
if errorlevel 1 (
    echo Docker no esta corriendo. Por favor, inicia Docker.
    pause
    exit /b 1
)

REM Levantar infraestructura (DB + API Gateway + RabbitMQ + Redis)
echo Levantando base de datos, API Gateway, RabbitMQ y Redis...
docker-compose -f dev-full-compose.yml up -d

REM Esperar a que la infraestructura este lista
echo Esperando que la infraestructura este lista...
timeout /t 15 /nobreak >nul

REM Ejecutar migraciones (agregando organization-svc)
echo Ejecutando migraciones...
cd auth-identity-svc
go run ./cmd/migrate
cd ..

cd address-svc
go run ./cmd/migrate
cd ..

cd person-svc
go run ./cmd/migrate
cd ..

cd organization-svc
go run ./cmd/migrate
cd ..

REM Preguntar si poblar la base de datos
echo.
echo ═══════════════════════════════════════════════════════════════
echo POBLACION DE BASE DE DATOS
echo ═══════════════════════════════════════════════════════════════
echo.
echo ¿Desea poblar la base de datos con datos mock para pruebas?
echo   • Incluye usuarios, personas, direcciones y contactos de ejemplo
echo   • Permite probar inmediatamente las APIs sin crear datos manualmente
echo   • Los datos pueden limpiarse posteriormente con clean.bat
echo.
echo   [S] - SI, poblar la base de datos con datos de prueba
echo   [N] - NO, continuar sin poblar (por defecto)
echo.

choice /c SN /n /m "Poblar base de datos? [S/N]: "
if errorlevel 2 goto skip_populate
if errorlevel 1 goto do_populate
REM Si choice falla, usar default
goto skip_populate

:do_populate
echo.
echo Poblando base de datos con datos de prueba...
docker exec -i rem-postgres-dev psql -U user -d rem_development < scripts\seed-dev-db.sql
if errorlevel 1 (
    echo ⚠️  Warning: Error poblando la base de datos. Los servicios seguirán funcionando pero sin datos de prueba.
) else (
    echo ✅ Base de datos poblada exitosamente con datos de prueba
    echo.
    echo 👥 USUARIOS DE PRUEBA DISPONIBLES:
    echo   • juan.perez@example.com (password: password)
    echo   • maria.gonzalez@example.com (password: password)
    echo   • admin@acmecorp.com (password: password) - EMPRESA
    echo   • Ver más en scripts/seed-dev-db.sql
)
goto continue_setup

:skip_populate
echo.
echo ⏭️  Saltando población de base de datos
echo   💡 Tip: Puedes poblar después con: .\scripts\seed-db.bat
goto continue_setup

:continue_setup

REM Verificar que Air este instalado
where air >nul 2>&1
if errorlevel 1 (
    echo Air no esta instalado. Instalando...
    go install github.com/cosmtrek/air@latest
    if errorlevel 1 (
        echo Error instalando Air. Por favor instalalo manualmente: go install github.com/cosmtrek/air@latest
        pause
        exit /b 1
    )
)

REM Iniciar servicios automaticamente
echo Iniciando servicios automaticamente...
echo.
echo 🌐 API Gateway: http://localhost:8081
echo   ├─ Auth API:         http://localhost:8081/api/auth/
echo   ├─ Users API:        http://localhost:8081/api/users/
echo   ├─ Person API:       http://localhost:8081/api/persons/
echo   ├─ Address API:      http://localhost:8081/api/addresses/
echo   ├─ Organization API: http://localhost:8081/api/organizations/
echo   └─ Health:           http://localhost:8081/health
echo.
echo 🔧 Servicios individuales:
echo   ├─ Auth Service:         http://localhost:4002
echo   ├─ Address Service:      http://localhost:4000  
echo   ├─ Person Service:       http://localhost:4001
echo   └─ Organization Service: http://localhost:4003
echo.
echo 🗄️ Infraestructura:
echo   ├─ PostgreSQL:    localhost:5432 (rem_development)
echo   ├─ RabbitMQ:      localhost:5672 (Management: http://localhost:15672)
echo   └─ Redis:         localhost:6379

REM Abrir terminales con Air para cada servicio (agregando organization-svc)
echo Abriendo terminales para cada servicio...

REM Terminal 1: Auth Identity Service
start "REM Auth Service" cmd /k "cd /d %cd%\auth-identity-svc && echo [AUTH] Iniciando Auth Identity Service... && air"

REM Terminal 2: Address Service  
start "REM Address Service" cmd /k "cd /d %cd%\address-svc && echo [ADDRESS] Iniciando Address Service... && air"

REM Terminal 3: Person Service
start "REM Person Service" cmd /k "cd /d %cd%\person-svc && echo [PERSON] Iniciando Person Service... && air"

REM Terminal 4: Organization Service
start "REM Organization Service" cmd /k "cd /d %cd%\organization-svc && echo [ORGANIZATION] Iniciando Organization Service... && air"

echo.
echo ✅ Entorno de desarrollo iniciado exitosamente!
echo.
echo 📡 USAR API GATEWAY para todas las peticiones:
echo    URL Base: http://localhost:8081
echo.
echo 📋 Ejemplos de uso:
echo    curl http://localhost:8081/api/persons/
echo    curl http://localhost:8081/api/addresses/
echo    curl http://localhost:8081/api/organizations/
echo    curl http://localhost:8081/api/auth/validate
echo.
echo 🛠️ Se han abierto 4 terminales con los servicios ejecutandose
echo Para detener todo, cierra las terminales o usa Ctrl+C en cada una
echo Para limpiar el entorno, ejecuta: .\scripts\clean.bat
echo.
echo 💡 Credenciales por defecto:
echo    RabbitMQ Management: http://localhost:15672 (user/supersecreta)
echo    PostgreSQL: user/supersecreta (rem_development)
echo.
echo Presiona cualquier tecla para salir...
pause >nul
