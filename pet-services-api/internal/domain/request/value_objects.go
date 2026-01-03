package request

import (
	"fmt"
	"strings"
)

// PetType representa o tipo de animal de estimação
type PetType string

const (
	PetTypeDog    PetType = "dog"
	PetTypeCat    PetType = "cat"
	PetTypeBird   PetType = "bird"
	PetTypeRabbit PetType = "rabbit"
	PetTypeFish   PetType = "fish"
	PetTypeRodent PetType = "rodent"
	PetTypeReptile PetType = "reptile"
	PetTypeOther  PetType = "other"
)

// String retorna a representação em string do tipo de pet
func (pt PetType) String() string {
	return string(pt)
}

// DisplayName retorna o nome amigável do tipo de pet em português
func (pt PetType) DisplayName() string {
	names := map[PetType]string{
		PetTypeDog:     "Cachorro",
		PetTypeCat:     "Gato",
		PetTypeBird:    "Pássaro",
		PetTypeRabbit:  "Coelho",
		PetTypeFish:    "Peixe",
		PetTypeRodent:  "Roedor",
		PetTypeReptile: "Réptil",
		PetTypeOther:   "Outro",
	}
	if name, ok := names[pt]; ok {
		return name
	}
	return string(pt)
}

// IsValid verifica se o tipo de pet é válido
func (pt PetType) IsValid() bool {
	validTypes := []PetType{
		PetTypeDog, PetTypeCat, PetTypeBird, PetTypeRabbit,
		PetTypeFish, PetTypeRodent, PetTypeReptile, PetTypeOther,
	}
	for _, valid := range validTypes {
		if pt == valid {
			return true
		}
	}
	return false
}

// PetInfo value object - representa informações sobre um pet
type PetInfo struct {
	Name   string
	Type   PetType
	Breed  string
	Age    int
	Weight float64
	Notes  string
}

// NewPetInfo cria um novo PetInfo com validações
func NewPetInfo(name string, petType PetType, breed string, age int, weight float64, notes string) (PetInfo, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return PetInfo{}, fmt.Errorf("nome do pet é obrigatório")
	}

	if !petType.IsValid() {
		return PetInfo{}, fmt.Errorf("tipo de pet inválido: %s", petType)
	}

	if age < 0 {
		return PetInfo{}, fmt.Errorf("idade do pet não pode ser negativa")
	}

	if age > 50 {
		return PetInfo{}, fmt.Errorf("idade do pet parece inválida (máximo 50 anos)")
	}

	if weight < 0 {
		return PetInfo{}, fmt.Errorf("peso do pet não pode ser negativo")
	}

	if weight > 500 {
		return PetInfo{}, fmt.Errorf("peso do pet parece inválido (máximo 500kg)")
	}

	return PetInfo{
		Name:   name,
		Type:   petType,
		Breed:  strings.TrimSpace(breed),
		Age:    age,
		Weight: weight,
		Notes:  strings.TrimSpace(notes),
	}, nil
}

// DisplayInfo retorna informações formatadas do pet
func (p PetInfo) DisplayInfo() string {
	info := fmt.Sprintf("%s (%s)", p.Name, p.Type.DisplayName())
	
	if p.Breed != "" {
		info += fmt.Sprintf(", Raça: %s", p.Breed)
	}
	
	if p.Age > 0 {
		ageText := "ano"
		if p.Age > 1 {
			ageText = "anos"
		}
		info += fmt.Sprintf(", %d %s", p.Age, ageText)
	}
	
	if p.Weight > 0 {
		info += fmt.Sprintf(", %.1f kg", p.Weight)
	}
	
	return info
}

// ShortInfo retorna informação resumida do pet
func (p PetInfo) ShortInfo() string {
	if p.Breed != "" {
		return fmt.Sprintf("%s (%s)", p.Name, p.Breed)
	}
	return fmt.Sprintf("%s (%s)", p.Name, p.Type.DisplayName())
}

// HasSpecialNeeds verifica se o pet tem necessidades especiais mencionadas
func (p PetInfo) HasSpecialNeeds() bool {
	return p.Notes != ""
}

// IsValid valida as informações do pet
func (p PetInfo) IsValid() error {
	if p.Name == "" {
		return fmt.Errorf("nome do pet é obrigatório")
	}
	if !p.Type.IsValid() {
		return fmt.Errorf("tipo de pet inválido")
	}
	if p.Age < 0 || p.Age > 50 {
		return fmt.Errorf("idade do pet inválida")
	}
	if p.Weight < 0 || p.Weight > 500 {
		return fmt.Errorf("peso do pet inválido")
	}
	return nil
}