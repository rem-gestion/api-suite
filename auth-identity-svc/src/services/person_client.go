package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	personpb "github.com/rem-gestion/rem-common/protos/person/v1"
	"google.golang.org/grpc"
)

// PersonServiceClient maneja la comunicación gRPC con person-svc
type PersonServiceClient struct {
	conn   *grpc.ClientConn
	client personpb.PersonServiceClient
	addr   string
}

// NewPersonServiceClient crea una nueva instancia del cliente
func NewPersonServiceClient(addr string) (*PersonServiceClient, error) {
	conn, err := rcgrpc.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to person service at %s: %w", addr, err)
	}

	client := personpb.NewPersonServiceClient(conn)

	return &PersonServiceClient{
		conn:   conn,
		client: client,
		addr:   addr,
	}, nil
}

// NewPersonServiceClientWithConn crea una nueva instancia del cliente usando una conexión existente
func NewPersonServiceClientWithConn(conn *grpc.ClientConn) *PersonServiceClient {
	client := personpb.NewPersonServiceClient(conn)

	return &PersonServiceClient{
		conn:   conn,
		client: client,
		addr:   "existing_connection",
	}
}

// Close cierra la conexión gRPC
func (c *PersonServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreatePerson crea una persona via gRPC
func (p *PersonServiceClient) CreatePerson(personData *dto.PersonData) (*personpb.PersonResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &personpb.CreatePersonRequest{}

	// Mapear tipo
	switch personData.Type {
	case "individual":
		req.Type = personpb.PersonType_INDIVIDUAL
	case "company":
		req.Type = personpb.PersonType_COMPANY
	default:
		return nil, fmt.Errorf("invalid person type: %s", personData.Type)
	}

	// Mapear sexo si está presente
	if personData.Sexo != "" {
		switch personData.Sexo {
		case "masculino":
			req.Sexo = personpb.Sexo_MASCULINO
		case "femenino":
			req.Sexo = personpb.Sexo_FEMENINO
		default:
			req.Sexo = personpb.Sexo_SEXO_UNSPECIFIED
		}
	}

	// Mapear avatar_url
	if personData.AvatarURL != "" {
		req.AvatarUrl = personData.AvatarURL
	}

	// Mapear detalles de persona
	if personData.Individual != nil && req.Type == personpb.PersonType_INDIVIDUAL {
		req.Details = &personpb.CreatePersonRequest_Individual{
			Individual: &personpb.CreateIndividual{
				FirstName: personData.Individual.FirstName,
				LastName:  personData.Individual.LastName,
				Dni:       personData.Individual.DocumentNumber,
			},
		}
	}

	if personData.Company != nil && req.Type == personpb.PersonType_COMPANY {
		req.Details = &personpb.CreatePersonRequest_Company{
			Company: &personpb.CreateCompany{
				LegalName:   personData.Company.LegalName,
				Cuit:        personData.Company.CUIT,
				SocietyType: personData.Company.SocietyType,
			},
		}
	}

	// Mapear dirección si está presente
	if personData.Address != nil {
		number := int32(0)

		// Usar zip si está presente, sino usar postal_code
		zipCode := personData.Address.Zip
		if zipCode == "" {
			zipCode = personData.Address.PostalCode
		}

		req.Address = &personpb.CreatePersonRequest_AddressPayload{
			AddressPayload: &personpb.AddressPayload{
				Floor:   personData.Address.Floor,
				Unit:    personData.Address.Unit,
				Street:  personData.Address.Street,
				Number:  number,
				City:    personData.Address.City,
				State:   personData.Address.State,
				Zip:     zipCode,
				Country: personData.Address.Country,
			},
		}
	}

	// Mapear contactos
	if len(personData.Contacts) > 0 {
		req.Contacts = make([]*personpb.AddContactPayload, len(personData.Contacts))
		for i, contact := range personData.Contacts {
			contactType := personpb.ContactType_CONTACT_TYPE_UNSPECIFIED
			switch contact.Type {
			case "email":
				contactType = personpb.ContactType_EMAIL
			case "phone":
				contactType = personpb.ContactType_PHONE
			case "whatsapp":
				contactType = personpb.ContactType_WHATSAPP
			}

			req.Contacts[i] = &personpb.AddContactPayload{
				Tipo:      contactType,
				Dato:      contact.Value,
				IsPrimary: contact.IsPrimary,
			}
		}
	}

	return p.client.CreatePerson(ctx, req)
}

// GetPerson obtiene una persona por ID via gRPC
func (c *PersonServiceClient) GetPerson(personID string) (*dto.PersonProfileData, error) {
	if personID == "" {
		return nil, fmt.Errorf("person ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &personpb.GetPersonRequest{
		Id: personID,
	}

	resp, err := c.client.GetPerson(ctx, req)
	if err != nil {
		log.Printf("Error getting person via gRPC: %v", err)
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	// Mapear respuesta a DTO
	personData := &dto.PersonProfileData{
		ID:        resp.Id,
		AvatarURL: resp.AvatarUrl,
	}

	// Mapear tipo
	switch resp.Type {
	case personpb.PersonType_INDIVIDUAL:
		personData.Type = "individual"
		if resp.GetIndividual() != nil {
			personData.Individual = &dto.IndividualData{
				FirstName:      resp.GetIndividual().FirstName,
				LastName:       resp.GetIndividual().LastName,
				DocumentNumber: resp.GetIndividual().Dni,
			}
		}
	case personpb.PersonType_COMPANY:
		personData.Type = "company"
		if resp.GetCompany() != nil {
			personData.Company = &dto.CompanyData{
				LegalName:   resp.GetCompany().LegalName,
				CUIT:        resp.GetCompany().Cuit,
				SocietyType: resp.GetCompany().SocietyType,
			}
		}
	}

	// Mapear sexo
	switch resp.Sexo {
	case personpb.Sexo_MASCULINO:
		personData.Sexo = "masculino"
	case personpb.Sexo_FEMENINO:
		personData.Sexo = "femenino"
	}

	// Mapear contactos
	if len(resp.Contacts) > 0 {
		personData.Contacts = make([]dto.ContactData, len(resp.Contacts))
		for i, contact := range resp.Contacts {
			contactType := ""
			switch contact.Tipo {
			case personpb.ContactType_EMAIL:
				contactType = "email"
			case personpb.ContactType_PHONE:
				contactType = "phone"
			case personpb.ContactType_WHATSAPP:
				contactType = "whatsapp"
			}

			personData.Contacts[i] = dto.ContactData{
				Type:      contactType,
				Value:     contact.Dato,
				IsPrimary: contact.IsPrimary,
			}
		}
	}

	return personData, nil
}

// HealthCheck verifica la conectividad con person-svc
func (c *PersonServiceClient) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Intentar una llamada simple para verificar conectividad
	_, err := c.client.ListPersons(ctx, &personpb.ListPersonsRequest{
		Page:    1,
		PerPage: 1,
	})

	if err != nil {
		return fmt.Errorf("person service health check failed: %w", err)
	}

	return nil
}
