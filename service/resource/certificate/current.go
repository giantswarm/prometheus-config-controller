package certificate

import (
	"context"
	"fmt"
	"path"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("debug", fmt.Sprintf("reading certificate directory: %s", r.certificateDirectory))

	fileInfos, err := afero.ReadDir(r.fs, r.certificateDirectory)
	if err != nil {
		return nil, microerror.Maskf(err, "could not read certificate directory")
	}

	certificateFiles := []certificateFile{}

	for _, fileInfo := range fileInfos {
		filePath := path.Join(r.certificateDirectory, fileInfo.Name())
		fileData, err := afero.ReadFile(r.fs, filePath)
		if err != nil {
			return nil, microerror.Maskf(err, "could not read certificate")
		}

		certificateFiles = append(certificateFiles, certificateFile{
			path: filePath,
			data: string(fileData),
		})
	}

	return certificateFiles, nil
}
