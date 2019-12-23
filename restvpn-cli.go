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

type TunnelParams struct {
	customer   string
	remoteIp   string
	remotePort string
	desc       string
	mask       string
	gateway    string
}

func getTunnels(apiAddr string, apiKey string, customer string) string {
	client := http.Client{}
	request, error := http.NewRequest("GET", apiAddr+"/restvpn/tunnels/"+customer, nil)
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

func postTunnel(apiAddr string, apiKey string, params TunnelParams) string {
	requestBody, error := json.Marshal(map[string]string{
		"customer":    params.customer,
		"remote_ip":   params.remoteIp,
		"remote_port": params.remotePort,
		"description": params.desc,
		"mask":        params.mask,
		"gateway":     params.gateway,
	})
	if error != nil {
		log.Fatalln("Failed composing requestBody: ", error)
	}
	client := http.Client{}
	request, error := http.NewRequest("POST", apiAddr+"/restvpn/tunnels", bytes.NewBuffer(requestBody))
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

func updateTunnel(apiAddr string, apiKey string, params TunnelParams) string {
	requestBody, error := json.Marshal(map[string]string{
		"remote_port": params.remotePort,
		"description": params.desc,
		"mask":        params.mask,
		"gateway":     params.gateway,
	})
	if error != nil {
		log.Fatalln("Failed composing requestBody: ", error)
	}
	client := http.Client{}
	request, error := http.NewRequest("PUT", apiAddr+"/restvpn/tunnels/"+params.customer+"/"+params.remoteIp, bytes.NewBuffer(requestBody))
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

func deleteTunnel(apiAddr string, apiKey string, params TunnelParams) string {
	client := http.Client{}
	request, error := http.NewRequest("DELETE", apiAddr+"/restvpn/tunnels/"+params.customer+"/"+params.remoteIp, nil)
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
		if len(flag.Args()) > 1 {
			fmt.Println("WARNING: 'list' command does not accept any arguments")
		}
		fmt.Print(getTunnels(apiAddr, apiKey, ""))
	case "get":
		getSubcommand := flag.NewFlagSet("get", flag.ExitOnError)

		var customer string
		getSubcommand.StringVar(&customer, "customer", "", "Customer name (required)")

		getSubcommand.Parse(os.Args[2:])
		if getSubcommand.Parsed() {
			if customer == "" {
				getSubcommand.PrintDefaults()
				os.Exit(1)
			}
			fmt.Print(getTunnels(apiAddr, apiKey, customer))
		}
	case "add":
		addSubcommand := flag.NewFlagSet("add", flag.ExitOnError)

		var customer string
		var remoteIp string
		var remotePort string
		var desc string
		var mask string
		var gateway string
		addSubcommand.StringVar(&customer, "customer", "", "Customer name (required)")
		addSubcommand.StringVar(&remoteIp, "ip", "", "Tunnel ip (required)")
		addSubcommand.StringVar(&remotePort, "port", "", "Tunnel port (required)")
		addSubcommand.StringVar(&desc, "desc", "", "Brief description")
		addSubcommand.StringVar(&mask, "mask", "", "Route netmask")
		addSubcommand.StringVar(&gateway, "gw", "", "Route gateway")

		addSubcommand.Parse(os.Args[2:]) // flag.Parse() for addSubcommand
		if addSubcommand.Parsed() {
			if customer == "" || remoteIp == "" || remotePort == "" {
				addSubcommand.PrintDefaults()
				os.Exit(1)
			}
			params := TunnelParams{
				customer:   customer,
				remoteIp:   remoteIp,
				remotePort: remotePort,
			}
			fmt.Print(postTunnel(apiAddr, apiKey, params))
		}
	case "update":
		updateSubcommand := flag.NewFlagSet("update", flag.ExitOnError)

		var customer string
		var remoteIp string
		var remotePort string
		var desc string
		var mask string
		var gateway string
		updateSubcommand.StringVar(&customer, "customer", "", "Customer name (required)")
		updateSubcommand.StringVar(&remoteIp, "ip", "", "Tunnel ip (required)")
		updateSubcommand.StringVar(&remotePort, "port", "", "Tunnel port")
		updateSubcommand.StringVar(&desc, "desc", "", "Brief description")
		updateSubcommand.StringVar(&mask, "mask", "", "Route netmask")
		updateSubcommand.StringVar(&gateway, "gw", "", "Route gateway")

		updateSubcommand.Parse(os.Args[2:])
		if updateSubcommand.Parsed() {
			if customer == "" || remoteIp == "" {
				updateSubcommand.PrintDefaults()
				os.Exit(1)
			}
			params := TunnelParams{
				customer:   customer,
				remoteIp:   remoteIp,
				remotePort: remotePort,
				desc:       desc,
				mask:       mask,
				gateway:    gateway,
			}
			fmt.Print(updateTunnel(apiAddr, apiKey, params))
		}
	case "delete":
		deleteSubcommand := flag.NewFlagSet("delete", flag.ExitOnError)

		var customer string
		var remoteIp string
		deleteSubcommand.StringVar(&customer, "customer", "", "Customer name (required)")
		deleteSubcommand.StringVar(&remoteIp, "ip", "", "Tunnel ip (required)")

		deleteSubcommand.Parse(os.Args[2:])
		if deleteSubcommand.Parsed() {
			if customer == "" || remoteIp == "" {
				deleteSubcommand.PrintDefaults()
				os.Exit(1)
			}
			params := TunnelParams{
				customer: customer,
				remoteIp: remoteIp,
			}
			fmt.Print(deleteTunnel(apiAddr, apiKey, params))
		}
	default:
		fmt.Println("HELP: 'list', 'get', 'add', 'update' or 'delete' subcommand is required")
		os.Exit(1)
	}
}
