#!/bin/bash

echo "🧹 Limpiando entorno REM..."

# Detener servicios de producción si están corriendo
if [ -f "scripts/stop.sh" ]; then
    echo "🛑 Deteniendo servicios de producción..."
    ./scripts/stop.sh
fi

# Detener contenedores de desarrollo
echo "📦 Deteniendo contenedores de desarrollo..."
if [ -f "dev-postgres-compose.yml" ]; then
    docker-compose -f dev-postgres-compose.yml down
fi

# Detener contenedores de producción si existen
echo "📦 Deteniendo contenedores de producción..."
if [ -f "postgres-compose.yml" ]; then
    docker-compose -f postgres-compose.yml down
fi

# Matar sesiones tmux relacionadas
echo "🖥️  Limpiando sesiones tmux..."
if command -v tmux &> /dev/null; then
    tmux kill-session -t rem-dev 2>/dev/null || true
    tmux kill-session -t rem-prod 2>/dev/null || true
fi

# Limpiar binarios compilados
echo "🗑️  Eliminando binarios compilados..."
find . -name "bin" -type d -exec rm -rf {} + 2>/dev/null || true

# Limpiar archivos temporales de Air
echo "🗑️  Eliminando archivos temporales de Air..."
find . -name "tmp" -type d -exec rm -rf {} + 2>/dev/null || true

# Limpiar logs de producción
echo "🗑️  Eliminando logs de producción..."
find . -name "logs" -type d -exec rm -rf {} + 2>/dev/null || true

# Limpiar archivos PID
echo "🗑️  Eliminando archivos PID..."
find . -name "*.pid" -type f -delete 2>/dev/null || true

# Limpiar archivos de build de Go
echo "🗑️  Eliminando archivos de build de Go..."
find . -name "go_build_*" -type f -delete 2>/dev/null || true
find . -name "__debug_bin*" -type f -delete 2>/dev/null || true

# Limpiar archivos temporales del sistema
echo "🗑️  Eliminando archivos temporales..."
find . -name ".DS_Store" -type f -delete 2>/dev/null || true
find . -name "Thumbs.db" -type f -delete 2>/dev/null || true
find . -name "*.tmp" -type f -delete 2>/dev/null || true

# Limpiar módulos de Go (opcional)
if [ "$1" = "--deep" ]; then
    echo "🧹 Limpieza profunda: eliminando cache de módulos..."
    go clean -modcache 2>/dev/null || true
    
    echo "🧹 Limpieza profunda: eliminando volúmenes Docker..."
    docker volume prune -f 2>/dev/null || true
fi

echo ""
echo "✅ Limpieza completada"
echo ""

if [ "$1" = "--deep" ]; then
    echo "🔄 Después de una limpieza profunda, es recomendable ejecutar:"
    echo "   go mod download  # En cada servicio"
    echo "   docker system prune -f  # Para limpiar Docker completamente"
fi

echo "💡 Para iniciar de nuevo:"
echo "   Desarrollo:  ./scripts/dev.sh"
echo "   Producción:  ./scripts/prod.sh"
