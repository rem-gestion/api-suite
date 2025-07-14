package router

//iniciales

/* # Organizations
POST   /organizations                 # Crear org
GET    /organizations                 # Listar (paginado + search)
GET    /organizations/:id             # Detalle
PUT    /organizations/:id             # Update (logo, matricula…)
DELETE /organizations/:id             # Soft-delete

# Branches
POST   /organizations/:id/branches    # Alta sucursal
GET    /organizations/:id/branches
PUT    /branches/:branchId
DELETE /branches/:branchId

# Settings
PUT    /organizations/:id/settings    # Upsert JSON
GET    /organizations/:id/settings

# Employees
POST   /organizations/:id/employees   # Agregar user existente
GET    /organizations/:id/employees
DELETE /employees/:employeeId

# Invites
POST   /organizations/:id/invites     # Invitar por email
GET    /organizations/:id/invites
POST   /invites/:token/accept         # Aceptar invitación (sin auth)
POST   /invites/:token/decline

# Ownership
POST   /organizations/:id/owners      # Agregar owner (person/company)
DELETE /owners/:ownerId

# Health & admin
GET    /health
GET    /health/detailed
*/
