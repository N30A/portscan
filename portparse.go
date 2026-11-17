package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	minPort = 1
	maxPort = 65535
)

func isValidPort(value string) bool {
	port, err := strconv.Atoi(value)
	if err != nil {
		return false
	}

	if port < minPort || port > maxPort {
		return false
	}

	return true
}

func hasDuplicates(s []string) bool {
	seen := make(map[string]bool)

	for _, v := range s {
		if seen[v] {
			return true
		}

		seen[v] = true
	}

	return false
}

func parsePorts(value string) ([]string, error) {
	if value == "--" {
		ports := make([]string, maxPort)
		for port := minPort; port <= maxPort; port++ {
			ports[port-1] = strconv.Itoa(port)
		}
		return ports, nil
	}

	if strings.Contains(value, ",") {
		return parsePortDelimiter(value)
	}

	if strings.Contains(value, "-") {
		return parsePortRange(value)
	}

	if !isValidPort(value) {
		return nil, fmt.Errorf("%s is not a valid port", value)
	}

	return []string{value}, nil
}

func parsePortDelimiter(value string) ([]string, error) {
	if value == "," {
		return nil, errors.New("no ports were specified between delimiter")
	}

	ports := strings.Split(value, ",")

	if hasDuplicates(ports) {
		return nil, fmt.Errorf("duplicate ports were specified")
	}

	for _, port := range ports {
		if !isValidPort(port) {
			return nil, fmt.Errorf("%s is not a valid port", port)
		}
	}

	return ports, nil
}

func parsePortRange(value string) ([]string, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return nil, errors.New("invalid port range format, ex: 1-1024")
	}

	first, err1 := strconv.Atoi(parts[0])
	last, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return nil, errors.New("invalid ports")
	}

	if first < minPort || last > maxPort {
		return nil, errors.New("invalid port range")
	}

	if first == last {
		return nil, errors.New("port range cannot be the same")
	}

	ports := make([]string, last)
	for port := first; port <= last; port++ {
		ports[port-1] = strconv.Itoa(port)
	}

	return ports, nil
}
