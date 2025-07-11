package repository

import (
	model "github.com/rem-gestion/api-suite/address/src/models"
)

// AddressRepository define la interfaz común para todos los tipos de repositorios de direcciones
type AddressRepository interface {
	Create(a *model.Address) (*model.Address, error)
	Get(id string) (*model.Address, error)
	Update(a *model.Address) error
	Delete(id string) error
}

// Verificar que ambos tipos implementan la interfaz
var (
	_ AddressRepository = (*AddressRepo)(nil)
	_ AddressRepository = (*AdaptiveAddressRepo)(nil)
)
