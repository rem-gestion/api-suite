#!/bin/bash

set -e

echo "🏭 Preparando entorno de PRODUCCIÓN REM..."

# Verificar que estamos en el directorio correcto
if [ ! -d "auth-identity-svc" ] || [ ! -d "address-svc" ] || [ ! -d "person-svc" ]; then
    echo "❌ Error: No se encuentran los directorios de servicios"
    echo "   Ejecuta este script desde el directorio services/"
    exit 1
fi

# Setear variables de entorno para producción
export REM_ENVIRONMENT=production

echo "🔧 Configurando entorno de producción..."
echo "   🗃️  Bases de datos: Separadas por servicio"
echo "   🔧 Configuración: .env.production"

# Función para compilar un servicio
build_service() {
    local service=$1
    echo "🔨 Compilando $service..."
    cd "$service-svc"
    
    # Crear directorio bin si no existe
    mkdir -p bin
    
    # Compilar
    if [ -f "cmd/api/main.go" ]; then
        go build -ldflags="-w -s" -o "bin/$service-svc" ./cmd/api
    elif [ -f "cmd/main.go" ]; then
        go build -ldflags="-w -s" -o "bin/$service-svc" ./cmd
    else
        echo "❌ No se encontró main.go para $service"
        cd ..
        return 1
    fi
    
    echo "✅ $service compilado exitosamente"
    cd ..
}

# Función para ejecutar migraciones de producción
run_prod_migrations() {
    local service=$1
    echo "🔄 Ejecutando migraciones de producción para $service..."
    cd "$service-svc"
    
    if [ -f "bin/$service-svc" ]; then
        # Si el binario soporta comando migrate
        ./bin/$service-svc migrate 2>/dev/null || {
            # Fallback a go run
            if [ -f "cmd/migrate/main.go" ]; then
                go run ./cmd/migrate
            else
                echo "⚠️  No se encontraron migraciones para $service"
            fi
        }
    else
        echo "❌ Binario no encontrado para $service"
        cd ..
        return 1
    fi
    cd ..
}

# Verificar variables de entorno críticas
check_env_vars() {
    local missing=()
    
    # Variables críticas que deben estar definidas en producción
    local required_vars=(
        "AUTH_API_KEY"
        "ADDRESS_API_KEY" 
        "PERSON_API_KEY"
    )
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            missing+=("$var")
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        echo "❌ Variables de entorno faltantes:"
        printf '   %s\n' "${missing[@]}"
        echo ""
        echo "💡 Ejemplo de configuración:"
        echo "   export AUTH_API_KEY='tu-clave-auth-secreta'"
        echo "   export ADDRESS_API_KEY='tu-clave-address-secreta'"
        echo "   export PERSON_API_KEY='tu-clave-person-secreta'"
        return 1
    fi
}

# Verificar variables de entorno
echo "🔍 Verificando variables de entorno..."
if ! check_env_vars; then
    exit 1
fi

echo "✅ Variables de entorno verificadas"

# Compilar todos los servicios
echo "🔨 Compilando servicios..."
build_service "auth-identity"
build_service "address"
build_service "person"

echo "✅ Todos los servicios compilados"

# Verificar conectividad a bases de datos antes de migrar
echo "🔗 Verificando conectividad a bases de datos..."
# Aquí podrías agregar checks de conectividad específicos

# Ejecutar migraciones de producción
echo "🔄 Ejecutando migraciones de producción..."
run_prod_migrations "auth-identity"
run_prod_migrations "address"
run_prod_migrations "person"

echo "✅ Migraciones completadas"

# Función para iniciar un servicio en background
start_service() {
    local service=$1
    local port=$2
    echo "🚀 Iniciando $service en puerto $port..."
    cd "$service-svc"
    
    # Crear archivo de log
    mkdir -p logs
    
    # Iniciar servicio en background
    nohup ./bin/$service-svc > logs/$service.log 2>&1 &
    local pid=$!
    echo $pid > logs/$service.pid
    
    echo "   PID: $pid"
    echo "   Log: $(pwd)/logs/$service.log"
    
    cd ..
    
    # Verificar que el servicio inició correctamente
    sleep 2
    if kill -0 $pid 2>/dev/null; then
        echo "✅ $service iniciado correctamente"
    else
        echo "❌ Error iniciando $service"
        return 1
    fi
}

# Iniciar servicios
echo ""
echo "🚀 Iniciando servicios en modo producción..."

start_service "auth-identity" "8080"
start_service "address" "4000"
start_service "person" "4001"

echo ""
echo "🎉 Entorno de producción iniciado exitosamente!"
echo ""
echo "📋 Servicios activos:"
echo "   🔐 Auth Service:    http://localhost:8080 (PID: $(cat auth-identity-svc/logs/auth-identity.pid))"
echo "   🏠 Address Service: http://localhost:4000 (PID: $(cat address-svc/logs/address.pid))"
echo "   👤 Person Service:  http://localhost:4001 (PID: $(cat person-svc/logs/person.pid))"
echo ""
echo "📋 Gestión de servicios:"
echo "   Para ver logs:    tail -f {servicio}-svc/logs/{servicio}.log"
echo "   Para detener:     kill \$(cat {servicio}-svc/logs/{servicio}.pid)"
echo "   Para detener todo: ./scripts/stop.sh"
echo ""

# Crear script de parada
cat > scripts/stop.sh << 'EOF'
#!/bin/bash
echo "🛑 Deteniendo servicios REM..."

for service in auth-identity address person; do
    if [ -f "${service}-svc/logs/${service}.pid" ]; then
        pid=$(cat "${service}-svc/logs/${service}.pid")
        if kill -0 $pid 2>/dev/null; then
            echo "   Deteniendo $service (PID: $pid)..."
            kill $pid
            rm -f "${service}-svc/logs/${service}.pid"
        else
            echo "   $service ya estaba detenido"
            rm -f "${service}-svc/logs/${service}.pid"
        fi
    fi
done

echo "✅ Todos los servicios detenidos"
EOF

chmod +x scripts/stop.sh

echo "💾 Script de parada creado: ./scripts/stop.sh"
