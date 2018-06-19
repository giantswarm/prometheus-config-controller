package key

import (
	"fmt"
	"path"
)

func certPath(certificateDirectory, clusterID, suffix string) string {
	return path.Join(certificateDirectory, fmt.Sprintf("%s-%s.pem", clusterID, suffix))
}

func CAPath(certificateDirectory, clusterID string) string {
	return certPath(certificateDirectory, clusterID, "ca")
}

func CrtPath(certificateDirectory, clusterID string) string {
	return certPath(certificateDirectory, clusterID, "crt")
}

func KeyPath(certificateDirectory, clusterID string) string {
	return certPath(certificateDirectory, clusterID, "key")
}
