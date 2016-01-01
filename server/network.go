package server
import (
	"net"
)

func IsLocalhost(hostname string) bool {
	host, _, err := net.SplitHostPort(hostname)
	if err != nil {
		panic(err)
	}

	names, err := net.LookupIP(host)
	if err != nil {
		panic("Could not resolve " + host + ": " + err.Error())
	}

	for _, name := range names {
		if name.IsLoopback() {
			return true
		}
	}
	return false
}
