package entities

type Breed struct {
	Base
	Name      string `json:"name"`
	SpeciesID string `json:"species_id"`
}

type BreedRepository interface {
	ListBySpecies(speciesID string) ([]*Breed, error)
}
