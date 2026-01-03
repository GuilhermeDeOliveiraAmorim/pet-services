package user

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// Email value object - representa um endereço de email válido
type Email struct {
	value string
}

// NewEmail cria e valida um novo email
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, ErrInvalidEmail
	}

	// RFC 5322 simplificado
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}

	// Valida tamanho máximo (RFC 5321)
	if len(email) > 254 {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: email}, nil
}

// String retorna o valor do email
func (e Email) String() string {
	return e.value
}

// Equals compara dois emails
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Phone value object - representa um número de telefone brasileiro
type Phone struct {
	value string
}

// NewPhone cria e valida um novo telefone
func NewPhone(phone string) (Phone, error) {
	// Remove todos caracteres não numéricos
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// Valida tamanho (10 ou 11 dígitos para números brasileiros)
	if len(phone) < 10 || len(phone) > 11 {
		return Phone{}, ErrInvalidPhone
	}

	// Valida DDD (códigos de 11 a 99)
	if len(phone) >= 2 {
		ddd := phone[0:2]
		if ddd < "11" || ddd > "99" {
			return Phone{}, ErrInvalidPhone
		}
	}

	return Phone{value: phone}, nil
}

// String retorna o valor bruto do telefone
func (p Phone) String() string {
	return p.value
}

// Formatted retorna o telefone formatado
func (p Phone) Formatted() string {
	if len(p.value) == 11 {
		// (XX) XXXXX-XXXX
		return fmt.Sprintf("(%s) %s-%s", p.value[0:2], p.value[2:7], p.value[7:11])
	}
	// (XX) XXXX-XXXX
	return fmt.Sprintf("(%s) %s-%s", p.value[0:2], p.value[2:6], p.value[6:10])
}

// Equals compara dois telefones
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}

// Location value object - representa uma localização geográfica
type Location struct {
	Latitude  float64
	Longitude float64
	Address   Address
}

// NewLocation cria uma nova localização com validação
func NewLocation(latitude, longitude float64, address Address) (Location, error) {
	// Valida range de latitude (-90 a 90)
	if latitude < -90 || latitude > 90 {
		return Location{}, fmt.Errorf("latitude deve estar entre -90 e 90")
	}

	// Valida range de longitude (-180 a 180)
	if longitude < -180 || longitude > 180 {
		return Location{}, fmt.Errorf("longitude deve estar entre -180 e 180")
	}

	return Location{
		Latitude:  latitude,
		Longitude: longitude,
		Address:   address,
	}, nil
}

// DistanceTo calcula a distância para outra localização em quilômetros
// Usa a fórmula de Haversine para calcular distância sobre a superfície de uma esfera
func (l Location) DistanceTo(other Location) float64 {
	const earthRadius = 6371.0 // raio da Terra em km

	// Converte graus para radianos
	lat1 := l.Latitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	deltaLat := (other.Latitude - l.Latitude) * math.Pi / 180
	deltaLon := (other.Longitude - l.Longitude) * math.Pi / 180

	// Fórmula de Haversine
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// IsWithinRadius verifica se outra localização está dentro de um raio em km
func (l Location) IsWithinRadius(other Location, radiusKM float64) bool {
	return l.DistanceTo(other) <= radiusKM
}

// Address value object - representa um endereço completo
type Address struct {
	Street     string
	Number     string
	Complement string
	District   string
	City       string
	State      string
	ZipCode    string
	Country    string
}

// NewAddress cria um novo endereço com validações básicas
func NewAddress(street, number, city, state, zipCode string) (Address, error) {
	if street == "" {
		return Address{}, fmt.Errorf("rua é obrigatória")
	}
	if city == "" {
		return Address{}, fmt.Errorf("cidade é obrigatória")
	}
	if state == "" {
		return Address{}, fmt.Errorf("estado é obrigatório")
	}

	// Valida CEP brasileiro (apenas dígitos, 8 caracteres)
	if zipCode != "" {
		cleanZip := regexp.MustCompile(`[^\d]`).ReplaceAllString(zipCode, "")
		if len(cleanZip) != 8 {
			return Address{}, fmt.Errorf("CEP inválido")
		}
		zipCode = cleanZip
	}

	return Address{
		Street:  street,
		Number:  number,
		City:    city,
		State:   state,
		ZipCode: zipCode,
		Country: "Brasil",
	}, nil
}

// FullAddress retorna o endereço completo formatado
func (a Address) FullAddress() string {
	parts := []string{}

	if a.Street != "" {
		streetPart := a.Street
		if a.Number != "" {
			streetPart += ", " + a.Number
		}
		if a.Complement != "" {
			streetPart += " - " + a.Complement
		}
		parts = append(parts, streetPart)
	}

	if a.District != "" {
		parts = append(parts, a.District)
	}

	if a.City != "" && a.State != "" {
		parts = append(parts, a.City+" - "+a.State)
	}

	if a.ZipCode != "" {
		parts = append(parts, "CEP: "+a.formatZipCode())
	}

	return strings.Join(parts, ", ")
}

// formatZipCode formata o CEP no padrão XXXXX-XXX
func (a Address) formatZipCode() string {
	if len(a.ZipCode) == 8 {
		return a.ZipCode[0:5] + "-" + a.ZipCode[5:8]
	}
	return a.ZipCode
}

// ShortAddress retorna versão resumida do endereço (rua e número)
func (a Address) ShortAddress() string {
	if a.Street != "" && a.Number != "" {
		return a.Street + ", " + a.Number
	}
	return a.Street
}

// Equals compara dois endereços
func (a Address) Equals(other Address) bool {
	return a.Street == other.Street &&
		a.Number == other.Number &&
		a.City == other.City &&
		a.State == other.State &&
		a.ZipCode == other.ZipCode
}
