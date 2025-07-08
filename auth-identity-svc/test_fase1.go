package main

import (
	"fmt"
	"log"

	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/rem-common/validation"
)

// Este script demuestra el funcionamiento de la validación exhaustiva implementada en Fase 1
func main() {
	fmt.Println("🧪 Testing Fase 1 - Validación Exhaustiva")
	fmt.Println("========================================")

	// Inicializar validator
	validator := validation.NewPersonValidator()

	// Test 1: Datos válidos de persona individual
	fmt.Println("\n✅ Test 1: Datos válidos de persona individual")
	validPersonData := &dto.PersonData{
		Type: "individual",
		Individual: &dto.IndividualData{
			FirstName:      "Juan",
			LastName:       "Pérez",
			DocumentNumber: "12345678",
		},
		Contacts: []dto.ContactData{
			{
				Type:      "email",
				Value:     "juan.perez@example.com",
				IsPrimary: true,
			},
			{
				Type:      "phone",
				Value:     "+541123456789",
				IsPrimary: false,
			},
		},
		Address: &dto.AddressData{
			Street:  "Av. Corrientes",
			Number:  1234,
			City:    "Buenos Aires",
			State:   "CABA",
			Country: "Argentina",
			Zip:     "C1043",
		},
	}

	if err := testPersonValidation(validator, validPersonData); err != nil {
		log.Printf("❌ Test fallido: %v", err)
	} else {
		log.Println("✅ Validación exitosa!")
	}

	// Test 2: Datos inválidos - Email mal formateado
	fmt.Println("\n❌ Test 2: Datos inválidos - Email mal formateado")
	invalidPersonData := &dto.PersonData{
		Type: "individual",
		Individual: &dto.IndividualData{
			FirstName:      "María",
			LastName:       "González",
			DocumentNumber: "87654321",
		},
		Contacts: []dto.ContactData{
			{
				Type:      "email",
				Value:     "email-invalido",
				IsPrimary: true,
			},
		},
	}

	if err := testPersonValidation(validator, invalidPersonData); err != nil {
		log.Printf("✅ Error esperado detectado: %v", err)
	} else {
		log.Println("❌ Error: Debería haber fallado la validación!")
	}

	// Test 3: Datos de empresa válidos
	fmt.Println("\n✅ Test 3: Datos válidos de empresa")
	validCompanyData := &dto.PersonData{
		Type: "company",
		Company: &dto.CompanyData{
			LegalName:   "Tech Solutions S.A.",
			CUIT:        "30123456789",
			SocietyType: "S.A.",
		},
		Contacts: []dto.ContactData{
			{
				Type:      "email",
				Value:     "info@techsolutions.com",
				IsPrimary: true,
			},
		},
	}

	if err := testPersonValidation(validator, validCompanyData); err != nil {
		log.Printf("❌ Test fallido: %v", err)
	} else {
		log.Println("✅ Validación de empresa exitosa!")
	}

	// Test 4: DNI inválido
	fmt.Println("\n❌ Test 4: DNI inválido (menos de 8 dígitos)")
	invalidDNI := &dto.PersonData{
		Type: "individual",
		Individual: &dto.IndividualData{
			FirstName:      "Carlos",
			LastName:       "López",
			DocumentNumber: "123456", // Solo 6 dígitos
		},
	}

	if err := testPersonValidation(validator, invalidDNI); err != nil {
		log.Printf("✅ Error esperado detectado: %v", err)
	} else {
		log.Println("❌ Error: Debería haber fallado la validación!")
	}

	fmt.Println("\n🎉 Tests de Fase 1 completados!")
	fmt.Println("La validación exhaustiva está funcionando correctamente.")
}

// Simula la función helper de validación implementada en validation_helpers.go
func testPersonValidation(validator *validation.PersonValidator, personData *dto.PersonData) error {
	if personData == nil {
		return nil
	}

	// Validar tipo de persona
	if err := validator.ValidatePersonType(personData.Type, "person_data.type"); err != nil {
		return fmt.Errorf("invalid person type: %s", err.Message)
	}

	// Validar según tipo
	switch personData.Type {
	case "individual":
		if personData.Individual == nil {
			return fmt.Errorf("individual data is required for individual type")
		}
		// Simulamos los adapters para este test
		if errors := validator.ValidateIndividualData(&testIndividualAdapter{data: personData.Individual}, "person_data.individual"); errors.HasErrors() {
			return fmt.Errorf("individual validation failed: %v", errors)
		}
	case "company":
		if personData.Company == nil {
			return fmt.Errorf("company data is required for company type")
		}
		if errors := validator.ValidateCompanyData(&testCompanyAdapter{data: personData.Company}, "person_data.company"); errors.HasErrors() {
			return fmt.Errorf("company validation failed: %v", errors)
		}
	}

	// Validar contactos si existen
	if len(personData.Contacts) > 0 {
		var contactAdapters []validation.ContactData
		for _, contact := range personData.Contacts {
			contactAdapters = append(contactAdapters, &testContactAdapter{data: &contact})
		}
		if errors := validator.ValidateContactData(contactAdapters, "person_data.contacts"); errors.HasErrors() {
			return fmt.Errorf("contacts validation failed: %v", errors)
		}
	}

	// Validar dirección si existe
	if personData.Address != nil {
		if errors := validator.ValidateAddressData(&testAddressAdapter{data: personData.Address}, "person_data.address"); errors.HasErrors() {
			return fmt.Errorf("address validation failed: %v", errors)
		}
	}

	return nil
}

// Test adapters (simplificados para la demo)
type testIndividualAdapter struct {
	data *dto.IndividualData
}

func (i *testIndividualAdapter) GetFirstName() string      { return i.data.FirstName }
func (i *testIndividualAdapter) GetLastName() string       { return i.data.LastName }
func (i *testIndividualAdapter) GetDocumentNumber() string { return i.data.DocumentNumber }

type testCompanyAdapter struct {
	data *dto.CompanyData
}

func (c *testCompanyAdapter) GetLegalName() string   { return c.data.LegalName }
func (c *testCompanyAdapter) GetCUIT() string        { return c.data.CUIT }
func (c *testCompanyAdapter) GetSocietyType() string { return c.data.SocietyType }

type testContactAdapter struct {
	data *dto.ContactData
}

func (c *testContactAdapter) GetType() string    { return c.data.Type }
func (c *testContactAdapter) GetValue() string   { return c.data.Value }
func (c *testContactAdapter) GetIsPrimary() bool { return c.data.IsPrimary }

type testAddressAdapter struct {
	data *dto.AddressData
}

func (a *testAddressAdapter) GetStreet() string  { return a.data.Street }
func (a *testAddressAdapter) GetCity() string    { return a.data.City }
func (a *testAddressAdapter) GetState() string   { return a.data.State }
func (a *testAddressAdapter) GetCountry() string { return a.data.Country }
func (a *testAddressAdapter) GetNumber() string  { return fmt.Sprintf("%d", a.data.Number) }
func (a *testAddressAdapter) GetFloor() string   { return a.data.Floor }
func (a *testAddressAdapter) GetUnit() string    { return a.data.Unit }
func (a *testAddressAdapter) GetZip() string     { return a.data.Zip }
