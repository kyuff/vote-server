package server
import (
	"testing"
)


func TestLocalHostLookup(t *testing.T) {
	for _, host := range ([]string{"127.0.0.1:0", "localhost:0", "localhost:8080", "127.0.0.1:8080"}) {
		if !IsLocalhost(host) {
			t.Error("Supposed to be localhost: " + host)
		}
	}


}


func TestLocalHostLookup_Errors(t *testing.T) {
	for _, host := range ([]string{"google.com:0", "pol.dk:0"}) {
		if IsLocalhost(host) {
			t.Error("Supposed to be localhost: " + host)
		}
	}
}


