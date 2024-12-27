package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

var (
	companyID    = "410"
	instanceType = "3"
	username     = "PoshAPICommandMW"
	password     = "RC2Thp9PtcQnGuUdvKJrHW"
	apiURL       = "https://devapi.teqtank.com/"
)

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized Method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	method := r.FormValue("method")
	data := r.FormValue("data")
	endpoint := r.FormValue("endpoint")

	if method == "" || endpoint == "" {
		http.Error(w, "Missing Parameters", http.StatusBadRequest)
		return
	}

	token := addToken(companyID, instanceType, username, password, apiURL)
	url := apiURL + endpoint

	switch method {
	case "login":
		customerUsername := r.FormValue("CustomerUsername")
		customerPassword := r.FormValue("CustomerPassword")
		headers := map[string]string{
			"Authorization":    "Bearer " + token,
			"Content-Type":     "application/json",
			"CustomerUsername": customerUsername,
			"CustomerPassword": customerPassword,
		}
		handleAPIRequest(w, "GET", url, headers, "")
	case "updatePassword":
		customerPassword := r.FormValue("CustomerPassword")
		headers := map[string]string{
			"Authorization":    "Bearer " + token,
			"Content-Type":     "application/json",
			"CustomerPassword": customerPassword,
		}
		handleAPIRequest(w, "GET", url, headers, "")
	case "resetPassword":
		customerUsername := r.FormValue("CustomerUsername")
		headers := map[string]string{
			"Authorization":    "Bearer " + token,
			"Content-Type":     "application/json",
			"CustomerUsername": customerUsername,
		}
		handleAPIRequest(w, "GET", url, headers, "")
	default:
		headers := map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}
		handleAPIRequest(w, method, url, headers, data)
	}
}

func addToken(companyID, instanceType, username, password, apiURL string) string {
	url := fmt.Sprintf("%s/Authorize/CompanyId/%s/%s", apiURL, companyID, instanceType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("MakoUsername", username)
	req.Header.Set("MakoPassword", password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching token: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading token response: %v", err)
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Fatalf("Error parsing token response: %v", err)
	}

	return tokenResponse.Data.Token
}

func handleAPIRequest(w http.ResponseWriter, method, url string, headers map[string]string, data string) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making API call", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
