package key

import (
	"fmt"
)

func CAPath(groupName string) string {
	return fmt.Sprintf("/certs/%s/ca.pem", groupName)
}

func CrtPath(groupName string) string {
	return fmt.Sprintf("/certs/%s/crt.pem", groupName)
}

func KeyPath(groupName string) string {
	return fmt.Sprintf("/certs/%s/key.pem", groupName)
}
