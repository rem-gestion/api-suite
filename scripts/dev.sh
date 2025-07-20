#!/bin/bash

set -e

echo "🚀 Iniciando entorno de DESARROLLO REM..."

# Verificar que estamos en el directorio correcto
if [ ! -f "dev-full-compose.yml" ]; then
    echo "❌ Error: No se encuentra dev-full-compose.yml"
    echo "   Ejecuta este script desde el directorio services/"
    exit 1
fi

# Setear variables de entorno para desarrollo
export REM_ENVIRONMENT=development

echo "📦 Configurando entorno de desarrollo..."
echo "   🗃️  Base de datos: rem_development (compartida)"
echo "   🐰 RabbitMQ: localhost:5672 (Management: localhost:15672)"
echo "   🔴 Redis: localhost:6379"
echo "   🔧 Configuración: .env.development"

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker no está corriendo. Por favor, inicia Docker."
    exit 1
fi

# Detener cualquier instancia previa
echo "🧹 Deteniendo contenedores previos..."
docker-compose -f dev-full-compose.yml down

# Levantar infraestructura completa (DB + RabbitMQ + Redis + API Gateway)
echo "📦 Levantando infraestructura completa..."
docker-compose -f dev-full-compose.yml up -d

# Esperar a que la infraestructura esté lista
echo "⏳ Esperando que la infraestructura esté lista..."
timeout=90
counter=0
while ! docker exec rem-postgres-dev pg_isready -U user -d rem_development > /dev/null 2>&1; do
    sleep 1
    counter=$((counter + 1))
    if [ $counter -gt $timeout ]; then
        echo "❌ Timeout esperando la base de datos"
        exit 1
    fi
done

echo "✅ PostgreSQL listo!"

# Verificar RabbitMQ
echo "⏳ Verificando RabbitMQ..."
counter=0
while ! docker exec rem-rabbitmq-dev rabbitmq-diagnostics -q ping > /dev/null 2>&1; do
    sleep 1
    counter=$((counter + 1))
    if [ $counter -gt $timeout ]; then
        echo "❌ Timeout esperando RabbitMQ"
        exit 1
    fi
done

echo "✅ RabbitMQ listo!"

# Verificar Redis
echo "⏳ Verificando Redis..."
counter=0
while ! docker exec rem-redis-dev redis-cli ping > /dev/null 2>&1; do
    sleep 1
    counter=$((counter + 1))
    if [ $counter -gt $timeout ]; then
        echo "❌ Timeout esperando Redis"
        exit 1
    fi
done

echo "✅ Redis listo!"
echo "✅ Toda la infraestructura está lista!"

# Función para ejecutar migraciones de un servicio
run_migrations() {
    local service=$1
    echo "🔄 Ejecutando migraciones de $service..."
    cd "$service-svc"
    if [ -f "cmd/migrate/main.go" ]; then
        go run ./cmd/migrate
    else
        echo "⚠️  No se encontraron migraciones para $service"
    fi
    cd ..
}

# Ejecutar migraciones en orden (agregando organization)
echo "🔄 Ejecutando migraciones..."
run_migrations "auth-identity"
run_migrations "address"
run_migrations "person"
run_migrations "property"
run_migrations "organization"

# Preguntar si poblar la base de datos
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "📊 POBLACION DE BASE DE DATOS"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "¿Desea poblar la base de datos con datos mock para pruebas?"
echo "  • Incluye usuarios, personas, direcciones, propiedades y contactos de ejemplo"
echo "  • Permite probar inmediatamente las APIs sin crear datos manualmente"
echo "  • Los datos pueden limpiarse posteriormente con clean.sh"
echo ""

echo "Opciones:"
echo "  [s] - SI, poblar la base de datos con datos de prueba"
echo "  [n] - NO, continuar sin poblar (por defecto)"
echo ""

while true; do
    printf "Poblar base de datos? [s/N]: "
    read -r POPULATE_DB
    
    # Si está vacío, usar N como default
    if [ -z "$POPULATE_DB" ]; then
        POPULATE_DB="n"
    fi
    
    case $POPULATE_DB in
        [sS]|[sS][iI]|[yY]|[yY][eE][sS])
            echo ""
            echo "🌱 Poblando base de datos con datos de prueba..."
            if docker exec -i rem-postgres-dev psql -U user -d rem_development < scripts/seed-dev-db.sql; then
                echo "✅ Base de datos poblada exitosamente con datos de prueba"
                echo ""
                echo "👥 USUARIOS DE PRUEBA DISPONIBLES:"
                echo "  • juan.perez@example.com (password: password)"
                echo "  • maria.gonzalez@example.com (password: password)"
                echo "  • admin@acmecorp.com (password: password) - EMPRESA"
                echo "  • Ver más en scripts/seed-dev-db.sql"
            else
                echo "⚠️  Warning: Error poblando la base de datos. Los servicios seguirán funcionando pero sin datos de prueba."
            fi
            break
            ;;
        [nN]|[nN][oO])
            echo ""
            echo "⏭️  Saltando población de base de datos"
            echo "  💡 Tip: Puedes poblar después con: ./scripts/seed-db.sh"
            break
            ;;
        *)
            echo ""
            echo "⚠️  Por favor ingresa 's' para poblar o 'n' para saltar."
            echo ""
            ;;
    esac
done

echo ""
echo "🌟 Entorno de desarrollo configurado exitosamente!"
echo ""
echo "📋 Servicios disponibles:"
echo "   🔐 Auth Service:         http://localhost:4002"
echo "   🏠 Address Service:      http://localhost:4000"  
echo "   👤 Person Service:       http://localhost:4001"
echo "   🏢 Organization Service: http://localhost:4003"
echo "   🏢 Property Service:     http://localhost:4004"
echo ""
echo "🗄️  Infraestructura:"
echo "   📊 PostgreSQL: postgresql://user:supersecreta@localhost:5432/rem_development"
echo "   🐰 RabbitMQ:   amqp://user:supersecreta@localhost:5672 (Management: http://localhost:15672)"
echo "   🔴 Redis:      redis://localhost:6379"
echo ""

# Iniciar servicios con Air si está disponible
if command -v air &> /dev/null; then
    echo "🚀 Iniciando servicios en modo desarrollo con Air..."
    echo "   Para detener: Ctrl+C"
    echo ""
    
    # Usar tmux si está disponible
    if command -v tmux &> /dev/null; then
        # Crear sesión tmux con 4 servicios
        tmux new-session -d -s rem-dev -c "$PWD/auth-identity-svc"
        tmux send-keys -t rem-dev "air" Enter
        
        # Crear ventanas adicionales
        tmux new-window -t rem-dev -c "$PWD/address-svc"
        tmux send-keys -t rem-dev "air" Enter
        
        tmux new-window -t rem-dev -c "$PWD/person-svc"
        tmux send-keys -t rem-dev "air" Enter
        
        tmux new-window -t rem-dev -c "$PWD/property-svc"
        tmux new-window -t rem-dev -c "$PWD/organization-svc"
        tmux send-keys -t rem-dev "air" Enter
        
        # Volver a la primera ventana
        tmux select-window -t rem-dev:0
        
        echo "🖥️  Sesión tmux 'rem-dev' creada con 4 ventanas"
        echo "   Para conectar: tmux attach -t rem-dev"
        echo "   Para listar ventanas: Ctrl+B, w"
        echo "   Para cambiar ventana: Ctrl+B, número"
        echo ""
        echo "🔗 Conectando a la sesión tmux..."
        tmux attach -t rem-dev
    else
        echo "📝 Tmux no está disponible. Ejecuta manualmente en terminales separadas:"
        echo ""
        echo "   Terminal 1: cd auth-identity-svc && air"
        echo "   Terminal 2: cd address-svc && air"
        echo "   Terminal 3: cd person-svc && air"
        echo "   Terminal 4: cd property-svc && air"
        echo "   Terminal 5: cd organization-svc && air"
    fi
else
    echo "⚠️  Air no está instalado. Instálalo con: go install github.com/cosmtrek/air@latest"
    echo ""
    echo "📝 Para iniciar manualmente:"
    echo "   Terminal 1: cd auth-identity-svc && go run ./cmd/api"
    echo "   Terminal 2: cd address-svc && go run ./cmd/api"
    echo "   Terminal 3: cd person-svc && go run ./cmd/api"
    echo "   Terminal 4: cd property-svc && go run ./cmd/api"
    echo "   Terminal 5: cd organization-svc && go run ./cmd/api"
fi
