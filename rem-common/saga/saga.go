package saga

import (
	"context"
	"fmt"
	"log"
)

// Step representa un paso en la saga
type Step struct {
	Name     string
	Execute  func(ctx context.Context) (interface{}, error)
	Rollback func(ctx context.Context, result interface{}) error
}

// Saga representa una transacción distribuida
type Saga struct {
	steps      []Step
	results    []interface{}
	completed  int
	rolledBack bool
}

// NewSaga crea una nueva saga
func NewSaga() *Saga {
	return &Saga{
		steps:   make([]Step, 0),
		results: make([]interface{}, 0),
	}
}

// AddStep agrega un paso a la saga
func (s *Saga) AddStep(step Step) *Saga {
	s.steps = append(s.steps, step)
	return s
}

// Execute ejecuta todos los pasos de la saga
func (s *Saga) Execute(ctx context.Context) error {
	for i, step := range s.steps {
		log.Printf("Executing saga step %d: %s", i+1, step.Name)

		result, err := step.Execute(ctx)
		if err != nil {
			log.Printf("Saga step %d failed: %s - Error: %v", i+1, step.Name, err)

			// Rollback todos los pasos completados
			if rollbackErr := s.rollback(ctx); rollbackErr != nil {
				return fmt.Errorf("step %s failed: %w, rollback also failed: %v", step.Name, err, rollbackErr)
			}

			return fmt.Errorf("saga failed at step %s: %w", step.Name, err)
		}

		s.results = append(s.results, result)
		s.completed++
		log.Printf("Saga step %d completed: %s", i+1, step.Name)
	}

	log.Printf("Saga completed successfully with %d steps", s.completed)
	return nil
}

// rollback ejecuta el rollback de todos los pasos completados en orden inverso
func (s *Saga) rollback(ctx context.Context) error {
	if s.rolledBack {
		return nil // Ya se hizo rollback
	}

	log.Printf("Starting saga rollback for %d completed steps", s.completed)
	s.rolledBack = true

	var rollbackErrors []error

	// Rollback en orden inverso
	for i := s.completed - 1; i >= 0; i-- {
		step := s.steps[i]
		result := s.results[i]

		if step.Rollback != nil {
			log.Printf("Rolling back step %d: %s", i+1, step.Name)

			if err := step.Rollback(ctx, result); err != nil {
				log.Printf("Rollback failed for step %d (%s): %v", i+1, step.Name, err)
				rollbackErrors = append(rollbackErrors, fmt.Errorf("rollback step %s: %w", step.Name, err))
			} else {
				log.Printf("Rollback completed for step %d: %s", i+1, step.Name)
			}
		}
	}

	if len(rollbackErrors) > 0 {
		return fmt.Errorf("rollback errors: %v", rollbackErrors)
	}

	log.Printf("Saga rollback completed successfully")
	return nil
}

// GetResult obtiene el resultado de un paso específico
func (s *Saga) GetResult(stepIndex int) interface{} {
	if stepIndex < 0 || stepIndex >= len(s.results) {
		return nil
	}
	return s.results[stepIndex]
}

// IsCompleted retorna true si la saga se completó exitosamente
func (s *Saga) IsCompleted() bool {
	return s.completed == len(s.steps) && !s.rolledBack
}

// IsRolledBack retorna true si se hizo rollback
func (s *Saga) IsRolledBack() bool {
	return s.rolledBack
}

// GetCompletedSteps retorna el número de pasos completados
func (s *Saga) GetCompletedSteps() int {
	return s.completed
}
