@echo off
echo.
echo 🧪 TESTING ENHANCED POSTMAN FEATURES - REM Platform
echo ═══════════════════════════════════════════════════════════════
echo.
echo Este script demuestra las funcionalidades mejoradas de Postman:
echo   ✅ Auto-inyección de API Keys y JWT tokens
echo   ✅ Auto-guardado de IDs de respuestas
echo   ✅ Auto-completado de variables en request bodies
echo   ✅ Environment poblado con datos mock del seed
echo   ✅ Scripts de pre-request y test inteligentes
echo.
echo ═══════════════════════════════════════════════════════════════
echo 📋 DATOS MOCK DISPONIBLES EN EL ENVIRONMENT:
echo ═══════════════════════════════════════════════════════════════

echo.
echo 👤 PERSONAS INDIVIDUALES:
echo   • Juan Pérez:     person_id = 11111111-1111-1111-1111-111111111111
echo   • María González: maria_person_id = 22222222-2222-2222-2222-222222222222
echo   • Carlos Rodríguez: carlos_person_id = 33333333-3333-3333-3333-333333333333 [PENDING]
echo   • Ana Martínez:   ana_person_id = 44444444-4444-4444-4444-444444444444
echo   • Laura Fernández: laura_person_id = 55555555-5555-5555-5555-555555555555 [Google OAuth]

echo.
echo 🏢 EMPRESAS:
echo   • ACME Corporation: acme_person_id = 66666666-6666-6666-6666-666666666666
echo   • Tech Solutions:   tech_person_id = 77777777-7777-7777-7777-777777777777
echo   • Innova Tech:      innovatech_person_id = 88888888-8888-8888-8888-888888888888

echo.
echo 🏠 DIRECCIONES:
echo   • Av. Corrientes 1234, Buenos Aires:     address_id = 11111111-1111-1111-1111-111111111111
echo   • Av. Santa Fe 2567, Piso 5 A, CABA:    maria_address_id = 22222222-2222-2222-2222-222222222222
echo   • Belgrano 890, Mar del Plata:          carlos_address_id = 33333333-3333-3333-3333-333333333333
echo   • Av. Pueyrredón 1456, Piso 12 B, CABA: ana_address_id = 44444444-4444-4444-4444-444444444444
echo   • San Martín 345, Rosario:              laura_address_id = 55555555-5555-5555-5555-555555555555

echo.
echo 📞 CONTACTOS:
echo   • Juan Pérez Email:    contact_id = 11111111-1111-1111-1111-111111111111
echo   • Juan Pérez Phone:    contact_phone_id = 11111111-2222-1111-1111-111111111111
echo   • Juan Pérez WhatsApp: contact_whatsapp_id = 11111111-3333-1111-1111-111111111111
echo   • María González Email: maria_contact_id = 22222222-1111-2222-2222-222222222222

echo.
echo 🔐 CREDENCIALES DE PRUEBA:
echo   • test_email = juan.perez@example.com
echo   • test_password = password (todos los usuarios usan "password")
echo   • api_key = supersecret-api-key-for-dev

echo.
echo ═══════════════════════════════════════════════════════════════
echo 🚀 CÓMO USAR EN POSTMAN:
echo ═══════════════════════════════════════════════════════════════

echo.
echo 1. 📥 Importa los archivos en Postman:
echo      • REM-API-Collection.postman_collection.json
echo      • REM-Development.postman_environment.json

echo.
echo 2. 🎯 Selecciona el environment "REM Development Environment"

echo.
echo 3. 🔐 Para probar autenticación:
echo      • Ejecuta "Auth ^> Login" con {{test_email}} y {{test_password}}
echo      • El JWT se guardará automáticamente en {{access_token}}

echo.
echo 4. 👤 Para probar personas:
echo      • Ejecuta "Persons ^> List All Persons" 
echo      • Los IDs se guardarán automáticamente
echo      • Usa {{person_id}} en otros requests

echo.
echo 5. 🏠 Para probar direcciones:
echo      • Ejecuta "Addresses ^> List All Addresses"
echo      • Los IDs se guardarán automáticamente
echo      • Usa {{address_id}} en otros requests

echo.
echo 6. 📞 Para probar contactos:
echo      • Ejecuta "Contacts ^> List All Contacts" 
echo      • Los IDs se guardarán automáticamente
echo      • Usa {{contact_id}} en otros requests

echo.
echo ═══════════════════════════════════════════════════════════════
echo 🎯 FUNCIONALIDADES AUTOMÁTICAS:
echo ═══════════════════════════════════════════════════════════════

echo.
echo ✅ Headers automáticos:
echo      • X-Api-Key se agrega automáticamente para endpoints que lo necesitan
echo      • Authorization Bearer se agrega automáticamente para endpoints autenticados

echo.
echo ✅ Auto-guardado de IDs:
echo      • person_id se guarda al crear/obtener personas
echo      • address_id se guarda al crear/obtener direcciones  
echo      • contact_id se guarda al crear/obtener contactos
echo      • access_token se guarda al hacer login

echo.
echo ✅ Auto-completado de variables:
echo      • Los request bodies usan variables como {{person_id}}, {{test_email}}, etc.
echo      • Se reemplazan automáticamente con valores del environment

echo.
echo ✅ Logging inteligente:
echo      • Cada request muestra logs en la consola de Postman
echo      • Los errores se muestran con detalles
echo      • Los IDs guardados se registran automáticamente

echo.
echo ═══════════════════════════════════════════════════════════════
echo 📖 DOCUMENTACIÓN ADICIONAL:
echo ═══════════════════════════════════════════════════════════════

echo.
echo 📄 README-Postman.md:     Guía completa de uso de las colecciones
echo 📄 Postman-Examples.md:   Ejemplos de payloads y casos de uso
echo 📄 README.md:             Documentación principal del proyecto

echo.
echo ✅ ¡Environment y colección listos para pruebas inmediatas!
echo    Todos los datos mock están precargados y los scripts automatizan
echo    la inyección de headers, guardado de IDs y reemplazo de variables.

pause
