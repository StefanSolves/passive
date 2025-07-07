Passive OSINT Tool
Overview
Passive OSINT v1.0.0 is a command-line tool designed to perform passive Open-Source Intelligence (OSINT) investigations. This tool allows users to:

Search for information about a person using their full name.
Retrieve details about an IP address, including ISP and geographic information.
Check for the presence of a username on popular social media platforms.
What is OSINT?
OSINT (Open-Source Intelligence) refers to the collection and analysis of publicly available data to gather useful information. It leverages freely accessible resources like websites, social media platforms, APIs, and online databases. OSINT investigations are conducted for purposes such as threat analysis, security assessment, law enforcement, or personal curiosity.

This tool strictly performs passive reconnaissance, meaning it does not interact aggressively with the targets, ensuring minimal risk of detection or intrusion.

Investigative Methods Used
Full Name Lookup

The program simulates searching for details about an individual by name. It use Random User API Returns realistic-looking user data, including names, addresses, and phone numbers.
IP Address Investigation

The tool queries the ipinfo.io API to retrieve details about the Internet Service Provider (ISP), city, and geographic coordinates (latitude and longitude) of the specified IP address.
Username Investigation

The tool constructs URLs for popular social media platforms (e.g., Facebook, Twitter, Instagram) using the given username.
It makes HTTP requests to check whether the username exists on these platforms based on the HTTP response status codes.
Local IP and Public IP Retrieval

The tool identifies the local IP address of the machine using system network interfaces.
It queries the api.ipify.org service to fetch the public IP address.
File Saving

Results are saved in result.txt, result2.txt, or result3.txt based on availability to ensure no existing results are overwritten.
How the Program Works
Installation
Ensure that you have:

Go Programming Language installed on your machine.
Internet connectivity for making API requests.
Clone this repository and run the program:

$ git clone https://github.com/your-repo/passive-osint.git
$ cd passive-osint
Usage
Run the tool using the following syntax:

$ go run main.go [OPTION] [INPUT]
Options
--help: Displays usage information.
$ go run main.go --help
Output:

Welcome to passive v1.0.0

OPTIONS:
    -fn         Search with full-name
    -ip         Search with IP address
    -u          Search with username
-fn: Perform a search by full name.
$ go run main.go -fn "Jean Dupont"
Output:

First name: Jean
Last name: Dupont
Address: 7 rue du Progrès, 75016 Paris
Number: +33601010101
Saved in result.txt
-ip: Investigate an IP address.
$ go run main.go -ip 127.0.0.1
Output:

ISP: FSociety, S.A.
City Lat/Lon: (13.731) / (-1.1373)
Saved in result2.txt
-u: Check for a username on social media platforms.
$ go run main.go -u "@user01"
Output:

Facebook : yes
Twitter : yes
Linkedin : yes
Instagram : no
Skype : yes
Saved in result3.txt
Steps to Build and Run as passive
To simplify usage, you can build the program into an executable and run it directly using the passive command:

Build the Executable:

$ go build -o passive main.go
This creates an executable named passive in the current directory.

Move the Executable to a Directory in Your PATH:

$ mv passive /usr/local/bin/
This ensures the executable can be run from anywhere using the passive command.

Run the Program: Now you can execute the program using:

$ passive [OPTION] [INPUT]
Example:
Display Help:

$ passive --help
Search by Full Name:

$ passive -fn "Jean Dupont"
Investigate IP Address:

$ passive -ip 127.0.0.1
Check for Username:

$ passive -u "@user01"
Features
Intuitive Interface: The tool offers a simple and clear interface with minimal configuration.
Automatic File Handling: Results are automatically saved to a new file to preserve historical data.
Accurate Data Retrieval: Utilizes APIs like ipinfo.io and api.ipify.org for reliable results.
Limitations
The social media checks rely on HTTP status codes, which may not always guarantee accurate results if platforms block automated requests.
Random User API for searchFullName its not providing real search capabilities for arbitrary names.
Contribution
Feel free to contribute to this project by submitting pull requests or reporting issues. Ensure all contributions comply with ethical guidelines.