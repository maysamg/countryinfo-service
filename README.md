# Testing av endepunkter

Dette dokumentet beskriver hvordan du kan teste API-et både **lokalt** og på **Render (Cloud Agent)**.


##  Før testing

1. Ved lokal testing, start serveren ved å kjøre følgende kommando i terminalen:
   ```
   go run main.go
   ```

2. Sørg for at API-et kjører på:
   ```
   http://localhost:8080
   ```
##  Viktig informasjon
⚠ **Render bruker en gratis instans som kan gå i dvale etter 15 min.**  
Hvis API-et ikke svarer med én gang, prøv igjen etter noen sekunder:)

## Testing av API-endepunkter lokalt og på Render

### 1. Hent landinformasjon
- Endpoint: `/countryinfo/v1/info/{landkode}`
- Eksempel-URL:
  ```
  http://localhost:8080/countryinfo/v1/info/no
  ```
    - Eksempel-URL-Render:
  ```
  https://countryinfo-service.onrender.com/countryinfo/v1/info/no
  ```
### 2. Hent befolkningsdata
- **Endpoint:** `/countryinfo/v1/population/{landkode}?limit={startår}-{sluttår}`
- **Eksempel-URL:**
  ```
  http://localhost:8080/countryinfo/v1/population/no?limit=2010-2015
  ```
    - **Eksempel-URL-Render:**
  ```
  https://countryinfo-service.onrender.com/countryinfo/v1/population/no?limit=2010-2015
    ```

### 3. Hent API-status
- **Endpoint:** `/countryinfo/v1/status/`
- **Eksempel-URL:**
  ```
  http://localhost:8080/countryinfo/v1/status/
  ```
    - **Eksempel-URL-Render:**
  ```
  https://countryinfo-service.onrender.com/countryinfo/v1/status/
    ```
- **Beskrivelse:** Gir en statusrapport for API-ene som brukes, inkludert uptime.


##  Testverktøy
Du kan teste API-et ved å bruke:
- **Postman** 

---



