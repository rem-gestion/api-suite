package services

import (
	"fmt"

	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/rem-common/validation"
)

// ValidatePersonDataHelper es una función helper para validar PersonData
func ValidatePersonDataHelper(validator *validation.PersonValidator, personData *dto.PersonData) error {
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
		if errors := validator.ValidateIndividualData(&IndividualDataAdapter{data: personData.Individual}, "person_data.individual"); errors.HasErrors() {
			return fmt.Errorf("individual validation failed: %v", errors)
		}
	case "company":
		if personData.Company == nil {
			return fmt.Errorf("company data is required for company type")
		}
		if errors := validator.ValidateCompanyData(&CompanyDataAdapter{data: personData.Company}, "person_data.company"); errors.HasErrors() {
			return fmt.Errorf("company validation failed: %v", errors)
		}
	}

	// Validar contactos si existen
	if len(personData.Contacts) > 0 {
		var contactAdapters []validation.ContactData
		for _, contact := range personData.Contacts {
			contactAdapters = append(contactAdapters, &ContactDataAdapter{data: &contact})
		}
		if errors := validator.ValidateContactData(contactAdapters, "person_data.contacts"); errors.HasErrors() {
			return fmt.Errorf("contacts validation failed: %v", errors)
		}
	}

	// Validar dirección si existe
	if personData.Address != nil {
		if errors := validator.ValidateAddressData(&AddressDataAdapter{data: personData.Address}, "person_data.address"); errors.HasErrors() {
			return fmt.Errorf("address validation failed: %v", errors)
		}
	}

	return nil
}

// Adapters para convertir DTOs a interfaces de validación

type IndividualDataAdapter struct {
	data *dto.IndividualData
}

func (i *IndividualDataAdapter) GetFirstName() string      { return i.data.FirstName }
func (i *IndividualDataAdapter) GetLastName() string       { return i.data.LastName }
func (i *IndividualDataAdapter) GetDocumentNumber() string { return i.data.DocumentNumber }

type CompanyDataAdapter struct {
	data *dto.CompanyData
}

func (c *CompanyDataAdapter) GetLegalName() string   { return c.data.LegalName }
func (c *CompanyDataAdapter) GetCUIT() string        { return c.data.CUIT }
func (c *CompanyDataAdapter) GetSocietyType() string { return c.data.SocietyType }

type ContactDataAdapter struct {
	data *dto.ContactData
}

func (c *ContactDataAdapter) GetType() string    { return c.data.Type }
func (c *ContactDataAdapter) GetValue() string   { return c.data.Value }
func (c *ContactDataAdapter) GetIsPrimary() bool { return c.data.IsPrimary }

type AddressDataAdapter struct {
	data *dto.AddressData
}

func (a *AddressDataAdapter) GetStreet() string  { return a.data.Street }
func (a *AddressDataAdapter) GetCity() string    { return a.data.City }
func (a *AddressDataAdapter) GetState() string   { return a.data.State }
func (a *AddressDataAdapter) GetCountry() string { return a.data.Country }
func (a *AddressDataAdapter) GetNumber() string  { return fmt.Sprintf("%d", a.data.Number) }
func (a *AddressDataAdapter) GetFloor() string   { return a.data.Floor }
func (a *AddressDataAdapter) GetUnit() string    { return a.data.Unit }
func (a *AddressDataAdapter) GetZip() string     { return a.data.Zip }
