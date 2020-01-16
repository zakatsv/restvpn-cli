package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RouteParams struct {
	commonName string
	remoteIp   string
	remotePort string
	desc       string
	netmask    string
}

func getRoutes(apiAddr string, apiKey string, commonName string) string {
	client := http.Client{}
	request, error := http.NewRequest("GET", apiAddr+"/restvpn/routes/"+commonName, nil)
	if error != nil {
		log.Fatalln(error)
	}
	if apiKey != "" {
		request.Header.Add("X-Api-Key", apiKey)
	}
	response, error := client.Do(request)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatalln("Failed reading response body: ", error)
	}
	return string(body)
}

func postRoute(apiAddr string, apiKey string, params RouteParams) string {
	requestBody, error := json.Marshal(map[string]string{
		"common_name": params.commonName,
		"remote_ip":   params.remoteIp,
		"remote_port": params.remotePort,
		"description": params.desc,
		"netmask":     params.netmask,
	})
	if error != nil {
		log.Fatalln("Failed composing requestBody: ", error)
	}
	client := http.Client{}
	request, error := http.NewRequest("POST", apiAddr+"/restvpn/routes", bytes.NewBuffer(requestBody))
	if error != nil {
		log.Fatalln(error)
	}
	request.Header.Add("Content-Type", "application/json")
	if apiKey != "" {
		request.Header.Add("X-Api-Key", apiKey)
	}
	response, error := client.Do(request)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatalln("Failed reading response body: ", error)
	}
	return string(body)
}

func updateRoute(apiAddr string, apiKey string, params RouteParams) string {
	requestBody, error := json.Marshal(map[string]string{
		"remote_port": params.remotePort,
		"description": params.desc,
		"netmask":     params.netmask,
	})
	if error != nil {
		log.Fatalln("Failed composing requestBody: ", error)
	}
	client := http.Client{}
	request, error := http.NewRequest("PUT", apiAddr+"/restvpn/routes/"+params.commonName+"/"+params.remoteIp, bytes.NewBuffer(requestBody))
	if error != nil {
		log.Fatalln(error)
	}
	request.Header.Add("Content-Type", "application/json")
	if apiKey != "" {
		request.Header.Add("X-Api-Key", apiKey)
	}
	response, error := client.Do(request)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatalln("Failed reading response body: ", error)
	}
	return string(body)
}

func deleteRoute(apiAddr string, apiKey string, params RouteParams) string {
	client := http.Client{}
	request, error := http.NewRequest("DELETE", apiAddr+"/restvpn/routes/"+params.commonName+"/"+params.remoteIp, nil)
	if error != nil {
		log.Fatalln(error)
	}
	if apiKey != "" {
		request.Header.Add("X-Api-Key", apiKey)
	}
	response, error := client.Do(request)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatalln("Failed reading response body: ", error)
	}
	return string(body)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("HELP: 'list', 'get', 'add', 'update' or 'delete' subcommand is required")
		os.Exit(1)
	}
	apiAddr := os.Getenv("RESTVPN_ADDR")
	apiKey := os.Getenv("RESTVPN_KEY")
	if apiAddr == "" {
		apiAddr = "http://localhost:5000"
	}

	switch os.Args[1] {
	case "-h", "--help":
		fmt.Println(`HELP:
		Warning: Make sure to set RESTVPN_ADDR and RESTVPN_KEY.
		CLI supports one of the following commands: list, get, add, update, delete`)
	case "list":
		if len(os.Args) > 2 {
			fmt.Println("WARNING: 'list' command does not accept any arguments")
		}
		fmt.Print(getRoutes(apiAddr, apiKey, ""))
	case "get":
		getSubcommand := flag.NewFlagSet("get", flag.ExitOnError)

		var commonName string
		getSubcommand.StringVar(&commonName, "cname", "", "Common name (required)")

		getSubcommand.Parse(os.Args[2:])
		if getSubcommand.Parsed() {
			if commonName == "" {
				getSubcommand.PrintDefaults()
				os.Exit(1)
			}
			fmt.Print(getRoutes(apiAddr, apiKey, commonName))
		}
	case "add":
		addSubcommand := flag.NewFlagSet("add", flag.ExitOnError)

		var commonName string
		var remoteIp string
		var remotePort string
		var desc string
		var netmask string
		addSubcommand.StringVar(&commonName, "cname", "", "Common name (required)")
		addSubcommand.StringVar(&remoteIp, "ip", "", "Remote ip (required)")
		addSubcommand.StringVar(&remotePort, "port", "", "Remote port (required)")
		addSubcommand.StringVar(&desc, "desc", "", "Brief description")
		addSubcommand.StringVar(&netmask, "mask", "", "Route netmask")

		addSubcommand.Parse(os.Args[2:]) // flag.Parse() for addSubcommand
		if addSubcommand.Parsed() {
			if commonName == "" || remoteIp == "" || remotePort == "" {
				addSubcommand.PrintDefaults()
				os.Exit(1)
			}
			params := RouteParams{
				commonName: commonName,
				remoteIp:   remoteIp,
				remotePort: remotePort,
			}
			fmt.Print(postRoute(apiAddr, apiKey, params))
		}
	case "update":
		updateSubcommand := flag.NewFlagSet("update", flag.ExitOnError)

		var commonName string
		var remoteIp string
		var remotePort string
		var desc string
		var netmask string
		updateSubcommand.StringVar(&commonName, "cname", "", "Common name (required)")
		updateSubcommand.StringVar(&remoteIp, "ip", "", "Remote ip (required)")
		updateSubcommand.StringVar(&remotePort, "port", "", "Remote port")
		updateSubcommand.StringVar(&desc, "desc", "", "Brief description")
		updateSubcommand.StringVar(&netmask, "mask", "", "Route netmask")

		updateSubcommand.Parse(os.Args[2:])
		if updateSubcommand.Parsed() {
			if commonName == "" || remoteIp == "" {
				updateSubcommand.PrintDefaults()
				os.Exit(1)
			}
			params := RouteParams{
				commonName: commonName,
				remoteIp:   remoteIp,
				remotePort: remotePort,
				desc:       desc,
				netmask:    netmask,
			}
			fmt.Print(updateRoute(apiAddr, apiKey, params))
		}
	case "delete":
		deleteSubcommand := flag.NewFlagSet("delete", flag.ExitOnError)

		var commonName string
		var remoteIp string
		deleteSubcommand.StringVar(&commonName, "cname", "", "Common name (required)")
		deleteSubcommand.StringVar(&remoteIp, "ip", "", "Remote ip (required)")

		deleteSubcommand.Parse(os.Args[2:])
		if deleteSubcommand.Parsed() {
			if commonName == "" || remoteIp == "" {
				deleteSubcommand.PrintDefaults()
				os.Exit(1)
			}
			params := RouteParams{
				commonName: commonName,
				remoteIp:   remoteIp,
			}
			fmt.Print(deleteRoute(apiAddr, apiKey, params))
		}
	default:
		fmt.Println("HELP: 'list', 'get', 'add', 'update' or 'delete' subcommand is required")
		os.Exit(1)
	}
}
