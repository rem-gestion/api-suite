package validation

import "fmt"

// PersonValidator valida datos de persona según esquemas de person-svc
type PersonValidator struct {
	BaseValidator
}

// NewPersonValidator crea una nueva instancia del validador de persona
func NewPersonValidator() *PersonValidator {
	return &PersonValidator{}
}

// ValidatePersonData valida estructura PersonData completa
func (pv *PersonValidator) ValidatePersonData(data interface{}) ValidationErrors {
	var errors ValidationErrors

	// Aquí necesitarías convertir el interface{} al tipo específico
	// Por ahora, usaremos una interfaz genérica que cada servicio puede implementar

	return errors
}

// PersonData representa la interfaz que debe implementar cualquier PersonData
type PersonData interface {
	GetType() string
	GetIndividual() IndividualData
	GetCompany() CompanyData
	GetContacts() []ContactData
	GetAddress() AddressData
}

type IndividualData interface {
	GetFirstName() string
	GetLastName() string
	GetDocumentNumber() string
}

type CompanyData interface {
	GetLegalName() string
	GetCUIT() string
	GetSocietyType() string
}

type ContactData interface {
	GetType() string
	GetValue() string
	GetIsPrimary() bool
}

type AddressData interface {
	GetStreet() string
	GetCity() string
	GetState() string
	GetCountry() string
	GetNumber() string
	GetFloor() string
	GetUnit() string
	GetZip() string
}

// ValidatePersonType valida el tipo de persona
func (pv *PersonValidator) ValidatePersonType(personType, fieldName string) *ValidationError {
	allowedTypes := []string{"individual", "company"}
	return pv.ValidateEnum(personType, fieldName, allowedTypes)
}

// ValidateIndividualData valida datos de persona individual
func (pv *PersonValidator) ValidateIndividualData(data IndividualData, fieldPrefix string) ValidationErrors {
	var errors ValidationErrors

	if err := pv.ValidateRequired(data.GetFirstName(), fieldPrefix+".first_name"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateMaxLength(data.GetFirstName(), fieldPrefix+".first_name", 64); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateRequired(data.GetLastName(), fieldPrefix+".last_name"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateMaxLength(data.GetLastName(), fieldPrefix+".last_name", 64); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateDNI(data.GetDocumentNumber(), fieldPrefix+".document_number"); err != nil {
		errors = append(errors, *err)
	}

	return errors
}

// ValidateCompanyData valida datos de empresa
func (pv *PersonValidator) ValidateCompanyData(data CompanyData, fieldPrefix string) ValidationErrors {
	var errors ValidationErrors

	if err := pv.ValidateRequired(data.GetLegalName(), fieldPrefix+".legal_name"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateMaxLength(data.GetLegalName(), fieldPrefix+".legal_name", 120); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateCUIT(data.GetCUIT(), fieldPrefix+".cuit"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateMaxLength(data.GetSocietyType(), fieldPrefix+".society_type", 12); err != nil {
		errors = append(errors, *err)
	}

	return errors
}

// ValidateContactData valida datos de contacto
func (pv *PersonValidator) ValidateContactData(contacts []ContactData, fieldPrefix string) ValidationErrors {
	var errors ValidationErrors

	for i, contact := range contacts {
		prefix := fmt.Sprintf("%s.contacts[%d]", fieldPrefix, i)

		// Validar tipo de contacto
		allowedTypes := []string{"email", "phone", "whatsapp"}
		if err := pv.ValidateEnum(contact.GetType(), prefix+".type", allowedTypes); err != nil {
			errors = append(errors, *err)
		}

		// Validar valor según tipo
		switch contact.GetType() {
		case "email":
			if err := pv.ValidateEmail(contact.GetValue(), prefix+".value"); err != nil {
				errors = append(errors, *err)
			}
		case "phone", "whatsapp":
			if err := pv.ValidatePhone(contact.GetValue(), prefix+".value"); err != nil {
				errors = append(errors, *err)
			}
		}

		// Validar que el valor no esté vacío
		if err := pv.ValidateRequired(contact.GetValue(), prefix+".value"); err != nil {
			errors = append(errors, *err)
		}

		// Validar longitud máxima del dato (según schema de person-svc)
		if err := pv.ValidateMaxLength(contact.GetValue(), prefix+".value", 128); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// ValidateAddressData valida datos de dirección
func (pv *PersonValidator) ValidateAddressData(data AddressData, fieldPrefix string) ValidationErrors {
	var errors ValidationErrors

	if err := pv.ValidateRequired(data.GetStreet(), fieldPrefix+".street"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateRequired(data.GetCity(), fieldPrefix+".city"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateRequired(data.GetState(), fieldPrefix+".state"); err != nil {
		errors = append(errors, *err)
	}

	if err := pv.ValidateRequired(data.GetCountry(), fieldPrefix+".country"); err != nil {
		errors = append(errors, *err)
	}

	return errors
}
