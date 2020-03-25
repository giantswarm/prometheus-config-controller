package certificate

import (
	"context"
	"fmt"
	"path"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("reading certificate directory: %s", r.certDirectory))

	fileInfos, err := afero.ReadDir(r.fs, r.certDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	certificateFiles := []certificateFile{}

	for _, fileInfo := range fileInfos {
		filePath := path.Join(r.certDirectory, fileInfo.Name())
		fileData, err := afero.ReadFile(r.fs, filePath)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		certificateFiles = append(certificateFiles, certificateFile{
			path: filePath,
			data: string(fileData),
		})
	}

	certificateCount.Set(float64(len(certificateFiles)))

	return certificateFiles, nil
}
