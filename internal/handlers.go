package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// **
// RootHandler håndterer forespørsler til "/"
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<h1>Welcome to the Country Info API!</h1>")
	fmt.Fprintln(w, "<p>Use the following endpoints to get data:</p>")
	fmt.Fprintln(w, "<ul>")
	fmt.Fprintln(w, "<li><a href='/countryinfo/v1/info/no'>/countryinfo/v1/info/{country_code}</a> - Get country info</li>")
	fmt.Fprintln(w, "<li><a href='/countryinfo/v1/population/no'>/countryinfo/v1/population/{country_code}</a> - Get population data</li>")
	fmt.Fprintln(w, "<li><a href='/countryinfo/v1/status'>/countryinfo/v1/status</a> - Check API status</li>")
	fmt.Fprintln(w, "</ul>")
	fmt.Fprintln(w, "<p>Enjoy using the API!</p>")
}

// InfoHandler håndterer forespørsler til /countryinfo/v1/info/
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	// Ekstraher landkoden fra URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Landkode er påkrevd. Eksempel: /countryinfo/v1/info/no", http.StatusBadRequest)
		return
	}
	countryCode := parts[4]

	// Hent 'limit' parameter (standard er 10)
	queryParams := r.URL.Query()
	limit := 10 // Standardverdi
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Hent landdata fra REST Countries API
	data, err := FetchCountryData(countryCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Kunne ikke hente data for landkode %s: %v", countryCode, err), http.StatusInternalServerError)
		return
	}

	// Hent bydata fra CountriesNow API
	countryName, ok := data["name"].(string)
	if !ok {
		http.Error(w, "Feil ved henting av landnavn", http.StatusInternalServerError)
		return
	}
	cities, err := FetchCitiesData(countryName)
	if err != nil {
		fmt.Printf("Advarsel: Kunne ikke hente byer for %s: %v\n", countryName, err)
	} else {
		// Sorter byene alfabetisk og begrens antall
		sort.Strings(cities)
		if len(cities) > limit {
			cities = cities[:limit]
		}
		data["cities"] = cities
	}

	// Returner JSON-respons
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		http.Error(w, "Kunne ikke formatere JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)

}

// PopulationHandler håndterer forespørsler til /countryinfo/v1/population/
func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Landkode er påkrevd. Eksempel: /countryinfo/v1/population/no", http.StatusBadRequest)
		return
	}
	countryCode := parts[4]

	// Håndter limit (startYear-endYear)
	queryParams := r.URL.Query()
	startYear, endYear := 0, 0

	if limit := queryParams.Get("limit"); limit != "" {
		yearRange := strings.Split(limit, "-")
		if len(yearRange) == 2 {
			startYear, _ = strconv.Atoi(yearRange[0])
			endYear, _ = strconv.Atoi(yearRange[1])
		}
	}

	// Debugging: Sjekk om parsing er riktig
	fmt.Println("Received limit:", queryParams.Get("limit"))
	fmt.Println("Parsed startYear:", startYear, "endYear:", endYear)

	// Hent befolkningsdata
	populationData, err := FetchCountryPopulationData(countryCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ingen befolkningsdata tilgjengelig for %s", countryCode), http.StatusNotFound)
		return
	}

	// Filtrer på årstall
	var filteredData []map[string]string
	totalPop, count := 0, 0

	for _, record := range populationData.Data.PopulationCounts {
		if startYear != 0 && endYear != 0 {
			if record.Year < startYear || record.Year > endYear {
				continue // Hopper over verdier utenfor intervallet
			}
		}

		filteredData = append(filteredData, map[string]string{
			"year":  strconv.Itoa(record.Year),
			"value": strconv.Itoa(record.Value),
		})

		totalPop += record.Value
		count++
	}

	// Beregn gjennomsnitt
	mean := 0
	if count > 0 {
		mean = totalPop / count
	}

	// Send JSON-respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mean":   mean,
		"values": filteredData,
	})
}

// StatusResponse representerer JSON-responsen for API-status
type StatusResponse struct {
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	APIStatus    string `json:"api_status"`
	CitiesAPI    string `json:"cities_api"`
	RequestCount int    `json:"request_count"`
}

// StartTime holder på tidspunktet for oppstart av serveren
var StartTime time.Time

// Variabel for å telle antall forespørsler
var requestCount int

// checkAPIStatus tester tilgjengeligheten til et eksternt API
func checkAPIStatus(url string, method string, body []byte) string {
	client := &http.Client{Timeout: 5 * time.Second}

	// Opprett forespørselen basert på HTTP-metoden
	var req *http.Request
	var err error
	if method == http.MethodPost {
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	}

	if err != nil {
		fmt.Printf("Feil ved opprettelse av forespørsel til %s: %v\n", url, err)
		return "ned"
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Feil ved sjekk av API-status (%s): %v\n", url, err)
		return "ned"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "opp"
	}

	fmt.Printf("API %s returnerte statuskode: %d\n", url, resp.StatusCode)
	return "ned"
}

// StatusHandler håndterer statusforespørsler
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Øk antall forespørsler
	requestCount++

	// Beregn uptime i sekunder
	uptime := int(time.Since(StartTime).Seconds())

	// API-endepunkter
	restCountriesAPI := "https://restcountries.com/v3.1/all"
	citiesAPI := "https://countriesnow.space/api/v0.1/countries/cities"

	// JSON-body for POST-forespørselen
	requestBody, _ := json.Marshal(map[string]string{"country": "Norway"})

	// Sjekk status for eksterne API-er
	restStatus := checkAPIStatus(restCountriesAPI, http.MethodGet, nil)
	citiesStatus := checkAPIStatus(citiesAPI, http.MethodPost, requestBody)

	// Lag JSON-respons
	status := map[string]interface{}{
		"countriesnowapi":  citiesStatus, // HTTP-statuskode for CountriesNow API
		"restcountriesapi": restStatus,   // HTTP-statuskode for REST Countries API
		"version":          "v1",         // API-versjon
		"uptime":           uptime,       // Tid i sekunder siden oppstart
	}

	// Send respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
