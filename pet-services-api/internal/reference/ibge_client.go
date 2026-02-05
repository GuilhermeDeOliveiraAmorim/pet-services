package reference

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const ibgeBaseURL = "https://servicodados.ibge.gov.br/api/v1/localidades"

type IBGEClient struct {
	httpClient *http.Client
	ttl        time.Duration
	mu         sync.RWMutex
	states     []State
	statesAt   time.Time
	cities     map[int][]City
	citiesAt   map[int]time.Time
}

func NewIBGEClient(httpClient *http.Client, ttl time.Duration) *IBGEClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &IBGEClient{
		httpClient: httpClient,
		ttl:        ttl,
		cities:     map[int][]City{},
		citiesAt:   map[int]time.Time{},
	}
}

type ibgeState struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}

type ibgeCity struct {
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	Microrregiao struct {
		Mesorregiao struct {
			UF ibgeState `json:"UF"`
		} `json:"mesorregiao"`
	} `json:"microrregiao"`
}

func (c *IBGEClient) ListStates(ctx context.Context) ([]State, error) {
	if states, ok := c.getCachedStates(); ok {
		return states, nil
	}

	url := fmt.Sprintf("%s/estados", ibgeBaseURL)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("IBGE states request failed: status %d: %s", response.StatusCode, string(body))
	}

	var payload []ibgeState
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	states := make([]State, 0, len(payload))
	for _, item := range payload {
		states = append(states, State{
			ID:   item.ID,
			Code: item.Sigla,
			Name: item.Nome,
		})
	}

	c.setCachedStates(states)
	return states, nil
}

func (c *IBGEClient) ListCities(ctx context.Context, stateID *int) ([]City, error) {
	cacheKey := 0
	if stateID != nil {
		cacheKey = *stateID
	}

	if cities, ok := c.getCachedCities(cacheKey); ok {
		return cities, nil
	}

	url := fmt.Sprintf("%s/municipios", ibgeBaseURL)
	if stateID != nil {
		url = fmt.Sprintf("%s/estados/%d/municipios", ibgeBaseURL, *stateID)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("IBGE cities request failed: status %d: %s", response.StatusCode, string(body))
	}

	var payload []ibgeCity
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	cities := make([]City, 0, len(payload))
	for _, item := range payload {
		currentStateID := cacheKey
		if stateID == nil {
			currentStateID = item.Microrregiao.Mesorregiao.UF.ID
		}
		cities = append(cities, City{
			ID:      item.ID,
			Name:    item.Nome,
			StateID: currentStateID,
		})
	}

	c.setCachedCities(cacheKey, cities)
	return cities, nil
}

func (c *IBGEClient) getCachedStates() ([]State, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.states) == 0 {
		return nil, false
	}
	if time.Since(c.statesAt) > c.ttl {
		return nil, false
	}
	return c.states, true
}

func (c *IBGEClient) setCachedStates(states []State) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states = states
	c.statesAt = time.Now()
}

func (c *IBGEClient) getCachedCities(cacheKey int) ([]City, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cached, ok := c.cities[cacheKey]
	if !ok || len(cached) == 0 {
		return nil, false
	}
	if time.Since(c.citiesAt[cacheKey]) > c.ttl {
		return nil, false
	}
	return cached, true
}

func (c *IBGEClient) setCachedCities(cacheKey int, cities []City) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cities[cacheKey] = cities
	c.citiesAt[cacheKey] = time.Now()
}
