package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// GetISO3Code henter ISO3-koden for et land fra REST Countries API
func GetISO3Code(countryCode string) (string, error) {
	url := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", countryCode)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Feil: Kunne ikke koble til REST Countries API for %s: %v\n", countryCode, err)
		return "", fmt.Errorf("kunne ikke koble til API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Feil: REST Countries API returnerte statuskode %d for %s\n", resp.StatusCode, countryCode)
		return "", fmt.Errorf("API returnerte statuskode %d", resp.StatusCode)
	}

	// Dekod JSON-respons
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Feil: Kunne ikke dekode JSON for %s: %v\n", countryCode, err)
		return "", fmt.Errorf("kunne ikke dekode JSON")
	}

	// Hent ISO3-koden
	if len(result) > 0 {
		if codes, ok := result[0]["cca3"].(string); ok {
			log.Printf("DEBUG: ISO3-kode for %s er %s\n", countryCode, codes)
			return codes, nil
		}
	}

	log.Printf("Feil: Fant ikke ISO3-kode for %s\n", countryCode)
	return "", fmt.Errorf("ISO3-kode ikke funnet")
}
