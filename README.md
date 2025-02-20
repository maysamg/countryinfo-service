# Testing av endepunkter

Dette dokumentet beskriver hvordan du kan teste API-et og de ulike endepunktene.

##  Før testing

1. Start serveren ved å kjøre følgende kommando i terminalen:
   ```
   go run main.go
   ```

2. Sørg for at API-et kjører på:
   ```
   http://localhost:8080
   ```

## Testing av API-endepunkter

### 1. Hent landinformasjon
- Endpoint: `/countryinfo/v1/info/{landkode}`
- Eksempel-URL:
  ```
  http://localhost:8080/countryinfo/v1/info/no
  ```
### 2. Hent befolkningsdata
- **Endpoint:** `/countryinfo/v1/population/{landkode}?limit={startår}-{sluttår}`
- **Eksempel-URL:**
  ```
  http://localhost:8080/countryinfo/v1/population/no?limit=2010-2015
  ```

### 3. Hent API-status
- **Endpoint:** `/countryinfo/v1/status/`
- **Eksempel-URL:**
  ```
  http://localhost:8080/countryinfo/v1/status/
  ```
- **Beskrivelse:** Gir en statusrapport for API-ene som brukes, inkludert uptime.


##  Testverktøy
Du kan teste API-et ved å bruke:
- **Postman** 

---



