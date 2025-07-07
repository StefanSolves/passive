package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Result struct to store the fetched data
type Result struct {
	QueryType   string
	Query       string
	FirstName   string
	LastName    string
	Address     string
	PhoneNumber string
	ISP         string
	City        string
	Lat         string
	Lon         string
	Socials     map[string]bool
}

func saveToFile(result Result) {
	var filename string

	// Determine the correct filename based on the query type
	switch result.QueryType {
	case "Full Name":
		filename = "result.txt"
	case "IP":
		filename = "result2.txt"
	case "Username":
		filename = "result3.txt"
	default:
		filename = "result.txt" // Default fallback
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the result to file
	output := fmt.Sprintf("QueryType: %s\nQuery: %s\n", result.QueryType, result.Query)
	if result.QueryType == "Full Name" {
		output += fmt.Sprintf("First Name: %s\nLast Name: %s\nAddress: %s\nPhone: %s\n", result.FirstName, result.LastName, result.Address, result.PhoneNumber)
	} else if result.QueryType == "IP" {
		output += fmt.Sprintf("ISP: %s\nCity: %s\nLat: %s\nLon: %s\n", result.ISP, result.City, result.Lat, result.Lon)
	} else if result.QueryType == "Username" {
		for social, exists := range result.Socials {
			status := "no"
			if exists {
				status = "yes"
			}
			output += fmt.Sprintf("%s: %s\n", social, status)
		}
	}

	output += "------------\n"
	file.WriteString(output)

	fmt.Println("Saved in", filename)
}

// searchFullName fetches user information for a full name
func searchFullName(firstName, lastName string) Result {
	apiUrl := fmt.Sprintf("https://randomuser.me/api/?name=%s+%s", firstName, lastName)

	// Make a GET request to the Random User API
	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return Result{}
	}
	defer resp.Body.Close()

	// Parse the JSON response
	body, _ := ioutil.ReadAll(resp.Body)
	var userData map[string]interface{}
	json.Unmarshal(body, &userData)

	// Extract data from the response
	results := userData["results"].([]interface{})
	if len(results) == 0 {
		fmt.Println("No results found for the given name.")
		return Result{}
	}

	user := results[0].(map[string]interface{})
	location := user["location"].(map[string]interface{})
	postcode := fmt.Sprintf("%v", location["postcode"]) // Safely convert postcode to string
	address := fmt.Sprintf("%s, %s %s",
		location["street"].(map[string]interface{})["name"],
		location["city"], postcode)

	phone := user["phone"].(string)

	fmt.Printf("First name: %s\nLast name: %s\nAddress: %s\nNumber: %s\n", firstName, lastName, address, phone)

	return Result{
		QueryType:   "Full Name",
		Query:       firstName + " " + lastName,
		FirstName:   firstName,
		LastName:    lastName,
		Address:     address,
		PhoneNumber: phone,
	}
}

// getLocalIP retrieves the local IP address of the machine
func getLocalIP() string {
	// Get a list of the system's network interfaces
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error fetching local IP address:", err)
		return "N/A"
	}

	for _, addr := range addrs {
		// Check if the address is not a loopback and is an IP address
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "N/A"
}

// getPublicIP retrieves the public IP of the machine by querying an external service
func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		fmt.Println("Error fetching public IP address:", err)
		return "N/A"
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading public IP response:", err)
		return "N/A"
	}

	return string(ip)
}

func isPrivateIP(ip string) bool {
	privateIPBlocks := []string{"127.", "10.", "172.", "192.168."}
	for _, block := range privateIPBlocks {
		if strings.HasPrefix(ip, block) {
			return true
		}
	}
	return false
}

func searchIP(ip string) Result {
	if isPrivateIP(ip) {
		fmt.Printf("The IP %s is a private or reserved IP address.\n", ip)
		return Result{
			QueryType: "IP",
			Query:     ip,
			ISP:       "N/A",
			City:      "N/A",
			Lat:       "N/A",
			Lon:       "N/A",
		}
	}

	apiUrl := fmt.Sprintf("http://ipinfo.io/%s/json", ip)
	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error fetching IP data:", err)
		return Result{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var ipInfo map[string]interface{}
	json.Unmarshal(body, &ipInfo)

	isp, ok := ipInfo["org"].(string)
	if !ok {
		isp = "N/A"
	}
	city, ok := ipInfo["city"].(string)
	if !ok {
		city = "N/A"
	}
	loc := fmt.Sprintf("%v", ipInfo["loc"])
	if ipInfo["loc"] == nil {
		loc = "N/A"
	}

	latLon := strings.Split(loc, ",")
	lat := "N/A"
	lon := "N/A"
	if len(latLon) == 2 {
		lat = latLon[0]
		lon = latLon[1]
	}

	fmt.Printf("ISP: %s\nCity: %s\nLat/Lon: (%s) / (%s)\n", isp, city, lat, lon)

	return Result{
		QueryType: "IP",
		Query:     ip,
		ISP:       isp,
		City:      city,
		Lat:       lat,
		Lon:       lon,
	}
}

func searchUsername(username string) Result {
	// Remove "@" if present
	username = strings.TrimPrefix(username, "@")

	socialPlatforms := map[string]string{
		"Facebook":  "https://www.facebook.com/%s",
		"Twitter":   "https://www.twitter.com/%s",
		"LinkedIn":  "https://www.linkedin.com/in/%s",
		"Instagram": "https://www.instagram.com/%s",
		"Skype":     "https://join.skype.com/invite/%s",
	}

	socials := make(map[string]bool)

	for platform, urlFormat := range socialPlatforms {
		url := fmt.Sprintf(urlFormat, username)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error checking %s: %v\n", platform, err)
			socials[platform] = false
			continue
		}
		defer resp.Body.Close()

		// If status code is 200, the username exists
		if resp.StatusCode == 200 {
			socials[platform] = true
		} else {
			socials[platform] = false
		}
	}

	// Print results
	for platform, exists := range socials {
		status := "no"
		if exists {
			status = "yes"
		}
		fmt.Printf("%s: %s\n", platform, status)
	}

	// Create the result object
	result := Result{
		QueryType: "Username",
		Query:     username,
		Socials:   socials,
	}

	// Save the result to a file
	saveToFile(result)

	return result
}

// validateInput validates the user input for full name, IP address, or username
func validateInput(input, option string) bool {
	switch option {
	case "-fn":
		// Ensure full name has both first and last name
		names := strings.Fields(input)
		if len(names) != 2 {
			fmt.Println("Invalid full name format. Please enter as 'First Last'.")
			return false
		}
	case "-ip":
		// Validate IP address format
		re := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
		if !re.MatchString(input) {
			fmt.Println("Invalid IP address format.")
			return false
		}
	case "-u":
		// Validate username format (ensure it's non-empty)
		if len(input) == 0 {
			fmt.Println("Invalid username format.")
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 3 || os.Args[1] == "--help" {
		fmt.Println("\nWelcome to passive v1.0.0\n")
		fmt.Println("OPTIONS:")
		fmt.Println("    -fn         Search with full-name")
		fmt.Println("    -ip         Search with IP address")
		fmt.Println("    -u          Search with username")
		return
	}

	option := os.Args[1]
	input := strings.Join(os.Args[2:], " ")

	if !validateInput(input, option) {
		return
	}

	var result Result

	switch option {
	case "-fn":
		names := strings.Fields(input)
		result = searchFullName(names[0], names[1])
	case "-ip":
		result = searchIP(input)
	case "-u":
		result = searchUsername(input)
	default:
		fmt.Println("Invalid option. Use --help for more information.")
		return
	}

	// Only fetch and display local and public IP information if it's not a username or full name query
	if option != "-u" && option != "-fn" {
		fmt.Println("Fetching local IP and ISP information...")

		// Get local IP address
		localIP := getLocalIP()
		fmt.Println("Local IP Address:", localIP)

		// Get public IP address
		publicIP := getPublicIP()
		fmt.Println("Public IP Address:", publicIP)

		// Fetch ISP and location info using the public IP
		result = searchIP(publicIP)

		// Save the result to a file (IP result)
		saveToFile(result)
	} else {
		// Save the full name or username result directly to file
		saveToFile(result)
	}
}
