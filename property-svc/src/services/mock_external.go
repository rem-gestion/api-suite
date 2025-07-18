package services

import (
	"github.com/google/uuid"
)

// MockExternalServices simula las respuestas de otros microservicios
// Esto se reemplazará por comunicación gRPC real en el futuro
type MockExternalServices struct{}

// MockAddress representa una dirección simulada
type MockAddress struct {
	ID      uuid.UUID `json:"id"`
	Street  string    `json:"street"`
	Number  int       `json:"number"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	Country string    `json:"country"`
}

// MockPerson representa una persona simulada
type MockPerson struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"` // "individual" o "company"
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	LegalName *string   `json:"legal_name,omitempty"`
}

// MockOrganization representa una organización simulada
type MockOrganization struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// GetMockAddresses retorna direcciones dummy para testing
func (m *MockExternalServices) GetMockAddresses() map[uuid.UUID]MockAddress {
	return map[uuid.UUID]MockAddress{
		uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"): {
			ID:      uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Street:  "Av. Corrientes",
			Number:  1234,
			City:    "Buenos Aires",
			State:   "CABA",
			Country: "AR",
		},
		uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"): {
			ID:      uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			Street:  "Av. Santa Fe",
			Number:  5678,
			City:    "Buenos Aires",
			State:   "CABA",
			Country: "AR",
		},
		uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"): {
			ID:      uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			Street:  "Av. Rivadavia",
			Number:  9000,
			City:    "Buenos Aires",
			State:   "CABA",
			Country: "AR",
		},
		uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"): {
			ID:      uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
			Street:  "Puerto Madero",
			Number:  100,
			City:    "Buenos Aires",
			State:   "CABA",
			Country: "AR",
		},
		uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"): {
			ID:      uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
			Street:  "Panamericana KM 40",
			Number:  0,
			City:    "Tigre",
			State:   "Buenos Aires",
			Country: "AR",
		},
	}
}

// GetMockPersons retorna personas dummy para testing
func (m *MockExternalServices) GetMockPersons() map[uuid.UUID]MockPerson {
	firstName1 := "Juan"
	lastName1 := "Pérez"
	firstName2 := "María"
	lastName2 := "González"
	legalName1 := "ACME Corp SA"
	legalName2 := "Inmobiliaria Del Centro SRL"

	return map[uuid.UUID]MockPerson{
		uuid.MustParse("11111111-1111-1111-1111-111111111111"): {
			ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Type:      "individual",
			FirstName: &firstName1,
			LastName:  &lastName1,
		},
		uuid.MustParse("22222222-2222-2222-2222-222222222222"): {
			ID:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Type:      "individual",
			FirstName: &firstName2,
			LastName:  &lastName2,
		},
		uuid.MustParse("33333333-3333-3333-3333-333333333333"): {
			ID:        uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			Type:      "company",
			LegalName: &legalName1,
		},
		uuid.MustParse("44444444-4444-4444-4444-444444444444"): {
			ID:        uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			Type:      "company",
			LegalName: &legalName2,
		},
	}
}

// GetMockOrganizations retorna organizaciones dummy para testing
func (m *MockExternalServices) GetMockOrganizations() map[uuid.UUID]MockOrganization {
	return map[uuid.UUID]MockOrganization{
		uuid.MustParse("99999999-9999-9999-9999-999999999999"): {
			ID:   uuid.MustParse("99999999-9999-9999-9999-999999999999"),
			Name: "REM Inmobiliaria Principal",
		},
		uuid.MustParse("88888888-8888-8888-8888-888888888888"): {
			ID:   uuid.MustParse("88888888-8888-8888-8888-888888888888"),
			Name: "Gestión Comercial Especializada",
		},
	}
}

// ValidateAddressExists simula la validación de existencia de dirección
func (m *MockExternalServices) ValidateAddressExists(addressID uuid.UUID) bool {
	addresses := m.GetMockAddresses()
	_, exists := addresses[addressID]
	return exists
}

// ValidatePersonExists simula la validación de existencia de persona
func (m *MockExternalServices) ValidatePersonExists(personID uuid.UUID) bool {
	persons := m.GetMockPersons()
	_, exists := persons[personID]
	return exists
}

// ValidateOrganizationExists simula la validación de existencia de organización
func (m *MockExternalServices) ValidateOrganizationExists(orgID uuid.UUID) bool {
	orgs := m.GetMockOrganizations()
	_, exists := orgs[orgID]
	return exists
}

// GetAddress simula obtener una dirección por ID
func (m *MockExternalServices) GetAddress(addressID uuid.UUID) (*MockAddress, bool) {
	addresses := m.GetMockAddresses()
	addr, exists := addresses[addressID]
	if !exists {
		return nil, false
	}
	return &addr, true
}

// GetPerson simula obtener una persona por ID
func (m *MockExternalServices) GetPerson(personID uuid.UUID) (*MockPerson, bool) {
	persons := m.GetMockPersons()
	person, exists := persons[personID]
	if !exists {
		return nil, false
	}
	return &person, true
}

// GetOrganization simula obtener una organización por ID
func (m *MockExternalServices) GetOrganization(orgID uuid.UUID) (*MockOrganization, bool) {
	orgs := m.GetMockOrganizations()
	org, exists := orgs[orgID]
	if !exists {
		return nil, false
	}
	return &org, true
}
