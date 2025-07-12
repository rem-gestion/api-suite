#!/bin/bash

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "SCRIPT DE POBLACION DE BASE DE DATOS - DESARROLLO"
echo "═══════════════════════════════════════════════════════════════"
echo

# Verificar que el contenedor de PostgreSQL esté corriendo
if ! docker ps --format "table {{.Names}}\t{{.Status}}" | grep -q rem-postgres-dev; then
    echo "❌ El contenedor rem-postgres-dev no está corriendo."
    echo "   Ejecuta primero: docker-compose -f dev-full-compose.yml up -d"
    echo "   O usa: ./scripts/dev.sh"
    exit 1
fi

echo "✅ Contenedor PostgreSQL detectado"
echo

# Verificar conectividad a la base de datos
echo "🔍 Verificando conectividad a la base de datos..."
if ! docker exec rem-postgres-dev pg_isready -U user -d rem_development > /dev/null 2>&1; then
    echo "❌ No se puede conectar a la base de datos rem_development"
    echo "   Verifica que el contenedor esté funcionando correctamente"
    exit 1
fi

echo "✅ Conectividad verificada"
echo

# Mostrar opción de limpiar datos existentes
echo "⚠️  ATENCION: Este script insertará datos de prueba en la base de datos."
echo "   Si ya existen datos, podrían generarse conflictos de claves duplicadas."
echo
read -p "¿Quieres limpiar los datos existentes antes de poblar? (s/N): " CLEAN_CHOICE

if [[ "$CLEAN_CHOICE" =~ ^[Ss]$ ]]; then
    echo
    echo "🧹 Limpiando datos existentes..."
    
    # Script para limpiar datos existentes respetando las FK
    docker exec -i rem-postgres-dev psql -U user -d rem_development << 'EOF'
BEGIN;
-- Limpiar en orden inverso a las dependencias
DELETE FROM person_contacto;
DELETE FROM auth_users;
DELETE FROM person_individual;
DELETE FROM person_company;
DELETE FROM person_person;
DELETE FROM auth_accounts;
DELETE FROM address_addresses;
COMMIT;
EOF
    
    echo "✅ Datos existentes limpiados"
    echo
fi

# Ejecutar script de población
echo "🌱 Poblando base de datos con datos de prueba..."
if docker exec -i rem-postgres-dev psql -U user -d rem_development < scripts/seed-dev-db.sql; then
    echo
    echo "✅ ¡Base de datos poblada exitosamente!"
    echo
    echo "═══════════════════════════════════════════════════════════════"
    echo "USUARIOS DE PRUEBA DISPONIBLES:"
    echo "═══════════════════════════════════════════════════════════════"
    echo
    echo "👤 INDIVIDUALES:"
    echo "  • juan.perez@example.com (password: password) - ACTIVO"
    echo "  • maria.gonzalez@example.com (password: password) - ACTIVO"
    echo "  • carlos.rodriguez@example.com (password: password) - PENDIENTE"
    echo "  • ana.martinez@example.com (password: password) - ACTIVO"
    echo "  • laura.fernandez@gmail.com (Google OAuth) - ACTIVO"
    echo
    echo "🏢 EMPRESAS:"
    echo "  • admin@acmecorp.com (password: password) - ACTIVO"
    echo "  • contacto@techsolutions.com (password: password) - ACTIVO"
    echo "  • ventas@innovatech.com.ar (password: password) - ACTIVO"
    echo
    echo "═══════════════════════════════════════════════════════════════"
    echo "DATOS INSERTADOS:"
    echo "  📍 10 direcciones en Argentina"
    echo "  🔐 8 cuentas de autenticación"
    echo "  👥 8 personas (5 individuales + 3 empresas)"
    echo "  📞 20+ contactos (emails, teléfonos, WhatsApp)"
    echo
    echo "🌐 Prueba los endpoints a través del API Gateway:"
    echo "  http://localhost:8081/api/persons/"
    echo "  http://localhost:8081/api/addresses/"
    echo "  http://localhost:8081/api/auth/"
    echo
else
    echo
    echo "❌ Error poblando la base de datos"
    echo "   Revisa los logs arriba para más detalles"
    echo
    exit 1
fi
