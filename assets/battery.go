package assets

import (
	"errors"
	"fmt"
)

// Charger is an interface onto a battery or any other type of asset which is able to charge or discharge.
type Charger interface {
	SetCharging(value float64) error
}

// Battery is an asset that implements the Charger interface
type Battery struct {
	ID string
}

// NewBattery creates a new Battery
func NewBattery(ID string) *Battery {
	return &Battery{ID: ID}
}

func (b Battery) SetCharging(value float64) error {
	if value == 1.0 {
		fmt.Printf("Battery %s is set to draw energy from the grid at its max rate.\n", b.ID)
		return nil
	}

	if value == -1.0 {
		fmt.Printf("Battery %s is set to dump energy into the grid at its max rate.\n", b.ID)
		return nil
	}

	return errors.New(fmt.Sprintf("unknown value: %f", value))
}
