package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// FetchCountryData henter og formatterer informasjon om et land fra REST Countries API
func FetchCountryData(countryCode string) (map[string]interface{}, error) {
	// Kall REST Countries API
	url := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", countryCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("kunne ikke koble til API: %v", err)
	}
	defer resp.Body.Close()

	// Håndter statuskode
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returnerte statuskode %d", resp.StatusCode)
	}

	// Dekod JSON-respons
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("kunne ikke dekode JSON: %v", err)
	}

	// Formater dataene
	if len(result) > 0 {
		data := result[0]
		formatted := map[string]interface{}{
			"name":       data["name"].(map[string]interface{})["common"],
			"continents": data["continents"],
			"population": data["population"],
			"languages":  data["languages"],
			"borders":    data["borders"],
			"flag":       data["flags"].(map[string]interface{})["png"],
			"capital":    data["capital"].([]interface{})[0],
			"cities":     []string{}, // Placeholder for nå
		}
		return formatted, nil
	}
	return nil, fmt.Errorf("ingen data funnet for landkode %s", countryCode)
}

// FetchCitiesData henter byinformasjon fra CountriesNow API
func FetchCitiesData(country string) ([]string, error) {
	// Kall CountriesNow API
	url := "http://129.241.150.113:3500/api/v0.1/countries/cities"
	payload := map[string]string{"country": country}
	jsonPayload, _ := json.Marshal(payload)

	// Send POST-forespørsel
	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("kunne ikke koble til API: %v", err)
	}
	defer resp.Body.Close()

	// Håndter statuskode
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returnerte statuskode %d", resp.StatusCode)
	}

	// Dekod JSON-respons
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("kunne ikke dekode JSON: %v", err)
	}

	// Hent ut byer
	if cities, ok := result["data"].([]interface{}); ok {
		cityList := []string{}
		for _, city := range cities {
			cityList = append(cityList, city.(string))
		}
		return cityList, nil
	}

	return nil, fmt.Errorf("ingen byer funnet for land %s", country)
}

// PopulationResponse representerer API-responsen for befolkningsdata
type PopulationResponse struct {
	Error bool `json:"error"`
	Data  struct {
		PopulationCounts []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

// FetchPopulationData henter befolkningsdata for en gitt by
func FetchPopulationData(city string) (*PopulationResponse, error) {
	apiURL := "http://129.241.150.113:3500/api/v0.1/countries/population/cities"

	// Opprett JSON-body med bynavn
	requestBody, err := json.Marshal(map[string]string{
		"city": city,
	})
	if err != nil {
		return nil, fmt.Errorf("kunne ikke serialisere forespørsel: %v", err)
	}

	// Lag en POST-forespørsel
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("kunne ikke lage forespørsel: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send forespørselen
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kunne ikke sende forespørsel: %v", err)
	}
	defer resp.Body.Close()

	// Sjekk HTTP-statuskode
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feil fra server: %v", resp.Status)
	}

	// Les respons-body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("kunne ikke lese respons: %v", err)
	}

	// Parse JSON-responsen
	var result PopulationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("kunne ikke parse JSON: %v", err)
	}

	// Returner resultatet
	return &result, nil
}

// FetchCountryPopulationData henter befolkningsdata for et gitt land

func FetchCountryPopulationData(countryCode string) (*PopulationResponse, error) {
	apiURL := "http://129.241.150.113:3500/api/v0.1/countries/population"

	// Hent ISO3-kode for landet
	iso3Code, err := GetISO3Code(countryCode)
	if err != nil {
		log.Printf("Feil: Kunne ikke finne ISO3-kode for %s: %v\n", countryCode, err)
		return nil, fmt.Errorf("kunne ikke finne ISO3-kode for %s", countryCode)
	}

	requestBody, _ := json.Marshal(map[string]string{"iso3": iso3Code})
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Feil: API-kall til CountriesNow feilet for %s: %v\n", countryCode, err)
		return nil, fmt.Errorf("API-kall feilet")
	}
	defer resp.Body.Close()

	// Sjekk HTTP-statuskode
	if resp.StatusCode != http.StatusOK {
		log.Printf("Feil: API returnerte statuskode %d for %s\n", resp.StatusCode, countryCode)
		return nil, fmt.Errorf("API returnerte statuskode %d", resp.StatusCode)
	}

	// Les respons-body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Feil: Kunne ikke lese API-respons for %s\n", countryCode)
		return nil, fmt.Errorf("kunne ikke lese respons")
	}

	// Debugging: Print API-responsen
	log.Printf("DEBUG: API-respons for %s: %s\n", countryCode, string(body))

	// Parse JSON
	var result PopulationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Feil: Kunne ikke parse JSON for %s: %v\n", countryCode, err)
		return nil, fmt.Errorf("kunne ikke parse JSON")
	}

	// Sjekk om API-en returnerte en gyldig respons
	if result.Error || len(result.Data.PopulationCounts) == 0 {
		log.Printf("Feil: Ingen befolkningsdata funnet for landet %s\n", countryCode)
		return nil, fmt.Errorf("ingen befolkningsdata funnet")
	}

	return &result, nil
}
