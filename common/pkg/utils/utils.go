package utils

import (
	"fmt"
	"log"
	"net"
)

func GetServiceAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Failed to get service address: %v", err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	log.Fatalf("Failed to find a suitable service address")
	return ""
}

func FormatServiceAddress(addr string) string {
	return addr + ":8080"
}

func FormatMemorySize[T int | []byte](input T) string {
	var size int

	switch v := any(input).(type) {
	case int:
		size = v
	case []byte:
		size = len(v)
	}

	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}
