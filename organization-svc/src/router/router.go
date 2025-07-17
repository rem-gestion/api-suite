package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/organization/src/controllers"
	"github.com/rem-gestion/rem-common/middleware"
)

// SetupRoutes configura todas las rutas del servicio organization-svc
// Organiza los endpoints en grupos lógicos según la funcionalidad
// jwtSecret: Secret compartido para validación local de JWT tokens
func SetupRoutes(r *gin.Engine, jwtSecret string) {
	// ============================================================================
	// MIDDLEWARES GLOBALES (en orden de ejecución)
	// ============================================================================

	// G-1: Request ID - Genera UUID único por petición para trazabilidad
	r.Use(middleware.RequestID())

	// G-2: Logger - Registra todas las peticiones HTTP con Gin logger
	r.Use(gin.Logger())

	// G-3: Recovery - Captura panics y devuelve 500 en lugar de crash
	r.Use(gin.Recovery())

	// G-6: Compression - Comprime respuestas grandes para ahorrar ancho de banda
	r.Use(middleware.GzipMiddleware())

	// G-7: Body Size Guard - Limita tamaño de payloads (1MB) para prevenir DoS
	r.Use(middleware.BodySizeLimit(1024 * 1024)) // 1MB

	// G-4: CORS - Configuración para peticiones cross-origin desde web UI
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ============================================================================
	// ENDPOINTS SIN AUTENTICACIÓN
	// ============================================================================

	// Health check endpoint (GET y HEAD para orquestadores)
	r.GET("/health", healthCheck)
	r.HEAD("/health", healthCheck)

	// Endpoints de métricas y debug (producción)
	// TODO: Aplicar admin_auth cuando esté disponible
	r.GET("/metrics", helloWorld)          // Prometheus metrics
	r.GET("/debug/pprof/*any", helloWorld) // Go pprof profiling

	// ============================================================================
	// API PRINCIPAL (con autenticación JWT obligatoria)
	// ============================================================================

	api := r.Group("")

	// A-1: JWT Authentication Local - Valida tokens usando secret compartido
	// IMPLEMENTADO: Validación local más eficiente que llamadas remotas
	api.Use(middleware.JWTAuthLocal(jwtSecret))

	// V-2: UUID Parameter Validation - Valida que :orgId sea UUID válido
	api.Use(middleware.UUIDParamValidator("orgId", "empId", "branchId", "roleId", "invId", "intId", "domId", "ownerId"))

	// V-3: Pagination Validation - Valida y normaliza parámetros de paginación
	api.Use(middleware.PaginationValidator())

	// ============================================================================
	// ENDPOINTS DE EJEMPLO (demuestran el JWT Auth Bridge funcionando)
	// ============================================================================

	// GET /test/auth - Endpoint de prueba para verificar que la autenticación funciona
	// Muestra la información del usuario autenticado desde el bridge
	api.GET("/test/auth", controllers.ExampleAuthenticatedController)

	// GET /test/organizations - Endpoint de ejemplo para listar organizaciones del usuario
	// Demuestra cómo usar la información del usuario autenticado en lógica de negocio
	api.GET("/test/organizations", controllers.ExampleOrganizationListController)

	// ============================================================================
	// ORGANIZATION CORE - Gestión básica de organizaciones
	// ============================================================================

	// POST /organizations - Crea una organización nueva
	// Recibe: display_name (requerido), logo_url (opcional), matricula (opcional), fiscal_address_id (opcional)
	// Respuesta: 201 Created + objeto Organization
	// Servicios externos: address-svc para validar/obtener dirección fiscal
	api.POST("/organizations", helloWorld)

	// GET /organizations - Lista organizaciones con filtros y paginación
	// Recibe: query strings para filtros (status, owner_id, search, page, per_page)
	// Respuesta: array de Organizations + metadata de paginación
	// Auth: usuario autenticado ve solo sus organizaciones
	api.GET("/organizations", helloWorld)

	// Grupo para endpoints específicos de una organización
	// Permite middleware de autorización por organización
	org := api.Group("/organizations/:orgId")
	{
		// GET /organizations/{orgId} - Obtiene una organización específica por ID
		// Recibe: orgId en path parameter
		// Respuesta: objeto Organization completo
		// Auth: usuario debe ser miembro de la organización
		org.GET("", helloWorld)

		// PUT /organizations/{orgId} - Actualiza datos generales de la organización
		// Recibe: display_name, logo_url, matricula, status (todos opcionales)
		// Respuesta: objeto Organization actualizado
		// Auth: requiere permisos de administrador en la organización
		org.PUT("", helloWorld)

		// PATCH /organizations/{orgId}/status - Cambia específicamente el estado
		// Recibe: status (active/suspended/deleted)
		// Respuesta: objeto Organization con nuevo status
		// Auth: solo administradores, validaciones especiales para 'deleted'
		org.PATCH("/status", helloWorld)

		// DELETE /organizations/{orgId} - Soft delete de la organización
		// Recibe: confirmación en headers o query
		// Respuesta: 204 No Content
		// Lógica: marca deleted_at, preserva datos para auditoría
		org.DELETE("", helloWorld) // ============================================================================
		// ORGANIZATION SETTINGS - Sistema clave-valor por organización
		// ============================================================================

		// GET /organizations/{orgId}/settings - Obtiene todas las configuraciones
		// Recibe: orgId en path
		// Respuesta: objeto JSON con todas las claves {currency: "USD", timezone: "UTC", ...}
		// Cache: configuraciones frecuentes se cachean en Redis
		org.GET("/settings", helloWorld)

		// GET /organizations/{orgId}/settings/{key} - Obtiene configuración específica
		// Recibe: key en path parameter
		// Respuesta: valor de la configuración específica
		// Cache: configuraciones frecuentes se cachean en Redis
		org.GET("/settings/:key", helloWorld)

		// PUT /organizations/{orgId}/settings - Actualiza/crea configuraciones en batch
		// Recibe: objeto JSON con pares clave-valor para actualizar
		// Respuesta: lista de configuraciones actualizadas
		// Validaciones: schemas específicos según la clave (ej: currency debe ser código ISO)
		org.PUT("/settings", helloWorld)

		// DELETE /organizations/{orgId}/settings/{key} - Elimina configuración específica
		// Recibe: key en path parameter
		// Respuesta: 204 No Content
		// Lógica: remueve la configuración, vuelve al valor por defecto si existe
		org.DELETE("/settings/:key", helloWorld) // ============================================================================
		// ORGANIZATION OWNERS - Gestión de propietarios
		// ============================================================================

		// GET /organizations/{orgId}/owners - Lista todos los propietarios
		// Recibe: filtros opcionales por owner_type, include_details
		// Respuesta: array de OrganizationOwner con porcentajes
		// Include: puede incluir detalles de person-svc si include_details=true
		// Útil para: UI de cap-table, reportes de ownership
		org.GET("/owners", helloWorld)

		// POST /organizations/{orgId}/owners - Añade un nuevo propietario
		// Recibe: owner_type (individual/company), owner_id, ownership_percentage
		// Respuesta: objeto OrganizationOwner creado
		// Servicios externos: person-svc para validar que owner_id existe
		// Validaciones: porcentajes no pueden exceder 100% total
		org.POST("/owners", helloWorld)

		// GET /organizations/{orgId}/owners/{ownerId} - Detalle de propietario específico
		// Recibe: ownerId en path parameter
		// Respuesta: objeto OrganizationOwner con detalles completos
		// Include: información completa del person-svc, historial de cambios
		org.GET("/owners/:ownerId", helloWorld)

		// DELETE /organizations/{orgId}/owners/{ownerId} - Remueve propietario
		// Recibe: ownerId en path
		// Respuesta: 204 No Content
		// Lógica: soft delete, recalcula porcentajes automáticamente
		org.DELETE("/owners/:ownerId", helloWorld)

		// ============================================================================
		// BRANCHES - Gestión de sucursales
		// ============================================================================

		// POST /organizations/{orgId}/branches - Crea nueva sucursal
		// Recibe: display_name, address_id (opcional), phone, email, is_main (opcional)
		// Respuesta: objeto OrganizationBranch creado
		// Servicios externos: address-svc para validar/obtener dirección
		// Lógica: si is_main=true, desmarca otras como principales
		org.POST("/branches", helloWorld)

		// GET /organizations/{orgId}/branches - Lista todas las sucursales
		// Recibe: filtros opcionales por status, include_address
		// Respuesta: array de OrganizationBranch
		// Include: puede incluir datos de address-svc si include_address=true
		org.GET("/branches", helloWorld)

		// GET /organizations/{orgId}/branches/{branchId} - Detalle de sucursal específica
		// Recibe: branchId en path, include query params opcionales
		// Respuesta: objeto OrganizationBranch completo
		// Include: dirección completa, estadísticas de empleados
		org.GET("/branches/:branchId", helloWorld)

		// PUT /organizations/{orgId}/branches/{branchId} - Actualiza sucursal
		// Recibe: campos a actualizar (display_name, phone, email, etc.)
		// Respuesta: objeto OrganizationBranch actualizado
		// Validaciones: no se puede cambiar is_main directamente (usar endpoint específico)
		org.PUT("/branches/:branchId", helloWorld)

		// PATCH /organizations/{orgId}/branches/{branchId}/make-main - Marca como sucursal principal
		// Recibe: confirmación opcional
		// Respuesta: objeto OrganizationBranch actualizado
		// Lógica: automáticamente desmarca la anterior sucursal principal
		org.PATCH("/branches/:branchId/make-main", helloWorld)

		// DELETE /organizations/{orgId}/branches/{branchId} - Elimina sucursal
		// Recibe: branchId en path
		// Respuesta: 204 No Content
		// Validaciones: no se puede eliminar si es la única sucursal o si tiene empleados activos
		org.DELETE("/branches/:branchId", helloWorld)

		// ============================================================================
		// ROLES - Gestión de roles internos
		// ============================================================================

		// POST /organizations/{orgId}/roles - Crea nuevo rol personalizado
		// Recibe: name, description (opcional), is_default (opcional)
		// Respuesta: objeto OrganizationRole creado
		// Validaciones: name debe ser único por organización
		org.POST("/roles", helloWorld)

		// GET /organizations/{orgId}/roles - Lista todos los roles
		// Recibe: filtros por is_default, include_usage_stats
		// Respuesta: array de OrganizationRole
		// Include: puede incluir estadísticas de uso (cuántos empleados tienen cada rol)
		org.GET("/roles", helloWorld)

		// GET /organizations/{orgId}/roles/{roleId} - Detalle de rol específico
		// Recibe: roleId en path
		// Respuesta: objeto OrganizationRole con empleados asignados
		// Include: lista de empleados que tienen este rol
		org.GET("/roles/:roleId", helloWorld)

		// PUT /organizations/{orgId}/roles/{roleId} - Actualiza rol
		// Recibe: name, description (campos editables)
		// Respuesta: objeto OrganizationRole actualizado
		// Restricciones: roles de sistema (is_default=true) tienen limitaciones de edición
		org.PUT("/roles/:roleId", helloWorld)

		// DELETE /organizations/{orgId}/roles/{roleId} - Elimina rol
		// Recibe: roleId en path, replacement_role_id opcional
		// Respuesta: 204 No Content
		// Lógica: reasigna empleados a replacement_role_id antes de eliminar
		org.DELETE("/roles/:roleId", helloWorld)

		// ============================================================================
		// INVITATIONS - Sistema de invitaciones
		// ============================================================================

		// POST /organizations/{orgId}/invites - Envía nueva invitación
		// Recibe: invitee_email, role_id, branch_id (opcional), expires_at (opcional)
		// Respuesta: objeto OrganizationInvite creado
		// Servicios externos: notification-svc para envío de email asíncrono
		// Lógica: genera token seguro, establece expiración automática
		org.POST("/invites", helloWorld)

		// GET /organizations/{orgId}/invites - Lista invitaciones con filtros
		// Recibe: filtros por status (pending/accepted/expired/cancelled), paginación
		// Respuesta: array de OrganizationInvite
		// Include: puede incluir estadísticas de aceptación
		org.GET("/invites", helloWorld)

		// GET /organizations/{orgId}/invites/{invId} - Detalle de invitación específica
		// Recibe: invId en path
		// Respuesta: objeto OrganizationInvite completo
		// Include: historial de reenvíos, intentos de aceptación
		org.GET("/invites/:invId", helloWorld)

		// POST /organizations/{orgId}/invites/{invId}/resend - Reenvía invitación
		// Recibe: invId en path, optional new_expires_at
		// Respuesta: 202 Accepted
		// Servicios externos: notification-svc para reenvío
		// Lógica: regenera token si la invitación estaba expirada
		org.POST("/invites/:invId/resend", helloWorld)

		// POST /organizations/{orgId}/invites/{invId}/cancel - Cancela invitación
		// Recibe: invId en path, cancellation_reason (opcional)
		// Respuesta: objeto OrganizationInvite actualizado
		// Lógica: marca como cancelled, invalida tokens
		org.POST("/invites/:invId/cancel", helloWorld)

		// ============================================================================
		// INTEGRATIONS - Hub de integraciones externas
		// ============================================================================

		// POST /organizations/{orgId}/integrations - Configura nueva integración
		// Recibe: integration_type_id, name, configuration, credentials (opcional)
		// Respuesta: objeto OrganizationIntegration creado
		// Servicios externos: APIs de terceros para validar credenciales según tipo
		// Lógica: encripta credenciales sensibles, valida configuración por tipo
		org.POST("/integrations", helloWorld)

		// GET /organizations/{orgId}/integrations - Lista integraciones activas
		// Recibe: filtros por status, integration_type, include_stats
		// Respuesta: array de OrganizationIntegration (sin credenciales sensibles)
		// Include: estadísticas de uso, último sync exitoso
		org.GET("/integrations", helloWorld)

		// GET /organizations/{orgId}/integrations/{intId} - Detalle de integración
		// Recibe: intId en path, include_config (solo admin)
		// Respuesta: objeto OrganizationIntegration completo
		// Security: credenciales solo visibles para administradores
		org.GET("/integrations/:intId", helloWorld)

		// GET /organizations/{orgId}/integrations/{intId}/config - Configuración de integración (solo admin)
		// Recibe: intId en path
		// Respuesta: configuración completa con credenciales ofuscadas
		// Security: solo administradores pueden ver credenciales (parcialmente ofuscadas)
		// Útil para: debugging, re-configuración sin resetear completamente
		org.GET("/integrations/:intId/config", helloWorld)

		// PUT /organizations/{orgId}/integrations/{intId} - Actualiza configuración
		// Recibe: configuration, credentials (campos editables)
		// Respuesta: objeto OrganizationIntegration actualizado
		// Servicios externos: re-validación con APIs de terceros
		// Lógica: re-encripta credenciales, invalida cache
		org.PUT("/integrations/:intId", helloWorld)

		// PATCH /organizations/{orgId}/integrations/{intId}/activate - Activa/desactiva integración
		// Recibe: status (active/inactive/error)
		// Respuesta: objeto OrganizationIntegration con nuevo status
		// Lógica: status 'error' requiere revisión manual, logs automáticos
		org.PATCH("/integrations/:intId/activate", helloWorld)

		// POST /organizations/{orgId}/integrations/{intId}/sync - Dispara sincronización manual
		// Recibe: sync_type (opcional: full/incremental)
		// Respuesta: 202 Accepted con job_id para tracking
		// Servicios externos: worker queue para procesamiento asíncrono
		// Lógica: rate limiting por organización, logs detallados
		org.POST("/integrations/:intId/sync", helloWorld)

		// GET /organizations/{orgId}/integrations/{intId}/events - Lista eventos de integración
		// Recibe: filtros por status, event_type, date_range
		// Respuesta: array de eventos con metadata
		// Include: detalles de payload, resultado, retries
		org.GET("/integrations/:intId/events", helloWorld)

		// GET /organizations/{orgId}/integrations/{intId}/logs - Historial de logs paginados
		// Recibe: filtros por level (info/warning/error), date_range, paginación
		// Respuesta: array de logs con timestamp y detalles
		// Performance: logs están particionados por fecha para queries eficientes
		org.GET("/integrations/:intId/logs", helloWorld)

		// DELETE /organizations/{orgId}/integrations/{intId} - Elimina integración
		// Recibe: intId en path, preserve_logs (opcional)
		// Respuesta: 204 No Content
		// Lógica: soft delete, desencripta y limpia credenciales, preserva logs si requested
		org.DELETE("/integrations/:intId", helloWorld)

		// ============================================================================
		// EMPLOYEES - Gestión de empleados
		// ============================================================================

		// POST /organizations/{orgId}/employees - Vincula usuario como empleado
		// Recibe: user_id, person_id (opcional), role_ids array, hired_at (opcional)
		// Respuesta: objeto Employee creado
		// Servicios externos: auth-identity-svc (validar user), person-svc (datos personales)
		// Lógica: crea automáticamente el vínculo persona-usuario si person_id es null
		org.POST("/employees", helloWorld)

		// GET /organizations/{orgId}/employees - Lista empleados con filtros
		// Recibe: filtros por status, role_id, branch_id, search, paginación
		// Respuesta: array de Employee con roles y datos básicos
		// Include: puede incluir datos de person-svc y estadísticas
		org.GET("/employees", helloWorld)

		// GET /organizations/{orgId}/employees/{empId} - Detalle completo de empleado
		// Recibe: empId en path, include query params
		// Respuesta: objeto Employee completo con roles, historial, branch
		// Include: datos personales completos, historial de roles, métricas
		org.GET("/employees/:empId", helloWorld)

		// GET /organizations/{orgId}/employees/{empId}/roles - Lista roles del empleado
		// Recibe: empId en path, include_historical (opcional)
		// Respuesta: array de EmployeeRole con detalles de cada rol
		// Include: puede incluir roles históricos si include_historical=true
		// Útil para: dashboards de empleados, auditoría de permisos
		org.GET("/employees/:empId/roles", helloWorld)

		// PUT /organizations/{orgId}/employees/{empId} - Actualiza datos de empleado
		// Recibe: hired_at, terminated_at, notes, metadata (campos editables)
		// Respuesta: objeto Employee actualizado
		// Validaciones: ciertas fechas no pueden ser futuras o inconsistentes
		org.PUT("/employees/:empId", helloWorld)

		// PATCH /organizations/{orgId}/employees/{empId}/status - Cambia estado del empleado
		// Recibe: status (active/inactive/terminated) y optional termination_reason
		// Respuesta: objeto Employee con nuevo status
		// Lógica: cambios de estado disparan workflows (notificaciones, cleanup de permisos)
		org.PATCH("/employees/:empId/status", helloWorld)

		// POST /organizations/{orgId}/employees/{empId}/roles/{roleId} - Asigna rol adicional
		// Recibe: is_primary (opcional, para designar rol principal)
		// Respuesta: objeto EmployeeRole creado
		// Lógica: si is_primary=true, desmarca otros roles como primarios
		org.POST("/employees/:empId/roles/:roleId", helloWorld)

		// PATCH /organizations/{orgId}/employees/{empId}/roles/{roleId}/primary - Marca rol como principal
		// Recibe: confirmación opcional
		// Respuesta: objeto EmployeeRole actualizado
		// Lógica: desmarca automáticamente otros roles como primarios, evita side-effects
		// Alternativa: más limpio que modificar is_primary en POST
		org.PATCH("/employees/:empId/roles/:roleId/primary", helloWorld)

		// DELETE /organizations/{orgId}/employees/{empId}/roles/{roleId} - Remueve rol
		// Recibe: empId y roleId en path
		// Respuesta: 204 No Content
		// Validaciones: empleado debe mantener al menos un rol activo
		org.DELETE("/employees/:empId/roles/:roleId", helloWorld)

		// DELETE /organizations/{orgId}/employees/{empId} - Termina empleado
		// Recibe: empId en path, termination_reason (opcional)
		// Respuesta: 204 No Content
		// Lógica: soft delete, actualiza status, preserva historial
		org.DELETE("/employees/:empId", helloWorld)

		// ============================================================================
		// CUSTOM DOMAINS - Gestión de dominios personalizados
		// ============================================================================

		// POST /organizations/{orgId}/domains - Registra nuevo dominio personalizado
		// Recibe: domain_name, subdomain (opcional), redirect_to_primary, is_primary (opcional)
		// Respuesta: objeto OrganizationDomain creado
		// Servicios externos: DNS provider para validaciones iniciales
		// Lógica: genera registros DNS requeridos, inicia verificación automática
		org.POST("/domains", helloWorld)

		// GET /organizations/{orgId}/domains - Lista dominios configurados
		// Recibe: filtros por status, include_dns_records
		// Respuesta: array de OrganizationDomain
		// Include: puede incluir records DNS y estado de verificación
		org.GET("/domains", helloWorld)

		// GET /organizations/{orgId}/domains/{domId} - Detalle completo de dominio
		// Recibe: domId en path, include_verification_details
		// Respuesta: objeto OrganizationDomain con DNS records y verificaciones
		// Include: historial de verificaciones, configuración SSL
		org.GET("/domains/:domId", helloWorld)

		// GET /organizations/{orgId}/domains/{domId}/verification-logs - Logs de verificación específicos
		// Recibe: domId en path, filtros por status, date_range
		// Respuesta: array de logs de verificación para este dominio específico
		// Include: detalles técnicos, records encontrados vs esperados
		// Útil para: debugging de problemas específicos de un dominio
		org.GET("/domains/:domId/verification-logs", helloWorld)

		// PUT /organizations/{orgId}/domains/{domId} - Actualiza configuración de dominio
		// Recibe: ssl_enabled, redirect_to_primary, auto_renew_ssl (campos editables)
		// Respuesta: objeto OrganizationDomain actualizado
		// Lógica: cambios de SSL disparan renovación automática
		org.PUT("/domains/:domId", helloWorld)

		// POST /organizations/{orgId}/domains/{domId}/verify - Inicia verificación DNS/SSL
		// Recibe: method (dns/http/email), force_recheck (opcional)
		// Respuesta: 202 Accepted con verification_job_id
		// Servicios externos: DNS checker job en background queue
		// Lógica: verifica records DNS, inicia proceso SSL si exitoso
		org.POST("/domains/:domId/verify", helloWorld)

		// POST /organizations/{orgId}/domains/{domId}/make-primary - Establece como dominio principal
		// Recibe: confirmación (require_confirmation=true)
		// Respuesta: objeto OrganizationDomain actualizado
		// Lógica: intercambia primary flags, configura redirects automáticos
		org.POST("/domains/:domId/make-primary", helloWorld)

		// DELETE /organizations/{orgId}/domains/{domId} - Elimina dominio personalizado
		// Recibe: domId en path, cleanup_dns (opcional)
		// Respuesta: 204 No Content
		// Lógica: soft delete, preserva logs, opcional cleanup de DNS records
		org.DELETE("/domains/:domId", helloWorld)

		// GET /organizations/{orgId}/domains/{domId}/dns-records - Obtiene registros DNS requeridos
		// Recibe: domId en path, record_type (opcional)
		// Respuesta: array de DomainDNS con registros necesarios
		// Include: instrucciones específicas por proveedor DNS
		org.GET("/domains/:domId/dns-records", helloWorld)

		// POST /organizations/{orgId}/domains/{domId}/dns-records/refresh - Re-verifica DNS
		// Recibe: domId en path
		// Respuesta: 202 Accepted con refresh_job_id
		// Servicios externos: DNS checker para re-validación
		// Lógica: útil cuando usuario ha actualizado DNS manualmente
		org.POST("/domains/:domId/dns-records/refresh", helloWorld)

		// ============================================================================
		// LOGS & MAINTENANCE - Auditoría específica por organización
		// ============================================================================

		// GET /organizations/{orgId}/invite-logs - Audit trail de invitaciones
		// Recibe: filtros por action_type, date_range, user_id, paginación
		// Respuesta: array de logs con acciones sobre invitaciones
		// Include: quien envió, aceptó, canceló invitaciones con timestamps
		org.GET("/invite-logs", helloWorld)

		// GET /organizations/{orgId}/domain-verification-logs - Historial de verificaciones DNS
		// Recibe: filtros por domain_id, verification_status, date_range
		// Respuesta: array de logs de verificaciones DNS/SSL
		// Include: detalles técnicos de fallos, records encontrados vs esperados
		org.GET("/domain-verification-logs", helloWorld)
	} // Cierre del grupo org

	// ============================================================================
	// INVITATIONS PUBLIC - Endpoints públicos para aceptar invitaciones
	// ============================================================================

	// POST /invites/accept - Acepta invitación usando token (endpoint público)
	// Recibe: token, person_data (opcional para crear perfil)
	// Respuesta: objeto Employee creado tras aceptación
	// Servicios externos: auth-identity-svc (crear/vincular usuario), person-svc (datos personales)
	// Lógica: valida token, crea empleado, invalida invitación
	api.POST("/invites/accept", helloWorld)

	// POST /invites/decline - Rechaza invitación usando token (endpoint público)
	// Recibe: token, decline_reason (opcional)
	// Respuesta: 204 No Content
	// Lógica: marca invitación como declined, opcional notificar a organización
	api.POST("/invites/decline", helloWorld)

	// GET /invites/validate - Valida token de invitación (endpoint público para landing page)
	// Recibe: token en query string
	// Respuesta: objeto InviteInfo con datos básicos para mostrar en UI
	// Include: nombre organización, rol ofrecido, tiempo restante
	api.GET("/invites/validate", helloWorld)

	// ============================================================================
	// MAINTENANCE - Administración del sistema
	// ============================================================================

	// GET /maintenance/run - Ejecuta mantenimiento programado (solo admin)
	// Recibe: maintenance_type (opcional): cleanup/optimize/partition
	// Respuesta: resultado de operaciones de mantenimiento
	// Lógica: ejecuta función run_organization_maintenance() de migración 0008
	// Security: endpoint protegido, solo accesible por administradores del sistema
	api.GET("/maintenance/run", helloWorld)
}

// healthCheck proporciona endpoint de health check básico
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "organization-svc",
		"version": "1.0.0",
	})
}

// helloWorld función temporal para endpoints sin implementar
// TODO: Reemplazar con implementaciones reales de controllers
func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World - Endpoint not implemented yet",
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	})
}
