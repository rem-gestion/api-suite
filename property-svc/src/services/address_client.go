package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AddressCreateRequest representa la estructura para crear una dirección
type AddressCreateRequest struct {
	Street  string `json:"street"`
	Number  int    `json:"number"`
	Floor   string `json:"floor,omitempty"`
	Unit    string `json:"unit,omitempty"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

// AddressResponse representa la respuesta del address service
type AddressResponse struct {
	ID      string `json:"id"`
	Street  string `json:"street"`
	Number  int    `json:"number"`
	Floor   string `json:"floor"`
	Unit    string `json:"unit"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

// AddressClient maneja la comunicación con el address service
type AddressClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewAddressClient crea una nueva instancia del cliente de address service
func NewAddressClient(baseURL, apiKey string, logger *zap.Logger) *AddressClient {
	return &AddressClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.Named("address-client"),
	}
}

// ValidateAddressExists verifica si una dirección existe
func (c *AddressClient) ValidateAddressExists(ctx context.Context, addressID uuid.UUID) bool {
	url := fmt.Sprintf("%s/api/addresses/%s", c.baseURL, addressID.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.Error("Error creating request", zap.Error(err))
		return false
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Error making request", zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	// Si el status es 200, la dirección existe
	return resp.StatusCode == http.StatusOK
}

// GetAddress obtiene una dirección por ID
func (c *AddressClient) GetAddress(ctx context.Context, addressID uuid.UUID) (*AddressResponse, error) {
	url := fmt.Sprintf("%s/api/addresses/%s", c.baseURL, addressID.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("address not found: %s", addressID.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var address AddressResponse
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &address, nil
}

// CreateAddress crea una nueva dirección si no existe
func (c *AddressClient) CreateAddress(ctx context.Context, address AddressCreateRequest) (*AddressResponse, error) {
	url := fmt.Sprintf("%s/api/addresses", c.baseURL)

	jsonData, err := json.Marshal(address)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error creating address: %s", string(body))
	}

	var newAddress AddressResponse
	if err := json.Unmarshal(body, &newAddress); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &newAddress, nil
}

// GetOrCreateAddress obtiene una dirección existente o la crea si no existe
func (c *AddressClient) GetOrCreateAddress(ctx context.Context, addressID uuid.UUID, fallbackAddress *AddressCreateRequest) (*AddressResponse, error) {
	// Primero intentar obtener la dirección existente
	if address, err := c.GetAddress(ctx, addressID); err == nil {
		c.logger.Info("Address found", zap.String("address_id", addressID.String()))
		return address, nil
	}

	// Si no existe y no tenemos datos para crear una nueva, retornar error
	if fallbackAddress == nil {
		return nil, fmt.Errorf("address %s not found and no fallback data provided", addressID.String())
	}

	// Crear la nueva dirección
	c.logger.Info("Address not found, creating new one", zap.String("address_id", addressID.String()))
	return c.CreateAddress(ctx, *fallbackAddress)
}
