package key

import (
	"fmt"
	"path"
)

func certPath(certificateDirectory, groupName, suffix string) string {
	return path.Join(certificateDirectory, fmt.Sprintf("%s-%s.pem", groupName, suffix))
}

func CAPath(certificateDirectory, groupName string) string {
	return certPath(certificateDirectory, groupName, "ca")
}

func CrtPath(certificateDirectory, groupName string) string {
	return certPath(certificateDirectory, groupName, "crt")
}

func KeyPath(certificateDirectory, groupName string) string {
	return certPath(certificateDirectory, groupName, "key")
}
