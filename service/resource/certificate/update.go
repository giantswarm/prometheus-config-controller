package certificate

import (
	"context"
	"fmt"
	"path"
	"reflect"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

func (r *Resource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	currentCertificateFiles, err := toCertificateFiles(currentState)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	desiredCertificateFiles, err := toCertificateFiles(desiredState)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	if currentCertificateFiles == nil || desiredCertificateFiles == nil {
		return nil, nil, nil, nil
	}

	if !reflect.DeepEqual(currentCertificateFiles, desiredCertificateFiles) {
		r.logger.Log("debug", "current certificates do not match desired certificates, need to update")
		return nil, nil, desiredCertificateFiles, nil
	}

	r.logger.Log("debug", "current certificates matches desired certificates, no update needed")

	return nil, nil, nil, nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	updateCertificateFiles, err := toCertificateFiles(updateState)
	if err != nil {
		return microerror.Mask(err)
	}

	// Write the update state certificate files.
	for _, fileToWrite := range updateCertificateFiles {
		r.logger.Log("debug", fmt.Sprintf("writing certificate: %s", fileToWrite.path))
		if err := afero.WriteFile(r.fs, fileToWrite.path, []byte(fileToWrite.data), r.certificatePermission); err != nil {
			return microerror.Mask(err)
		}
	}

	// Remove any unwanted certificate files.
	fileInfos, err := afero.ReadDir(r.fs, r.certificateDirectory)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, fileInfo := range fileInfos {
		fileDesired := false
		filePath := path.Join(r.certificateDirectory, fileInfo.Name())

		for _, updateCertificateFile := range updateCertificateFiles {
			if filePath == updateCertificateFile.path {
				fileDesired = true
			}
		}

		if !fileDesired {
			r.logger.Log("debug", fmt.Sprintf("removing certificate: %s", filePath))
			if err := r.fs.Remove(filePath); err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func toCertificateFiles(v interface{}) ([]certificateFile, error) {
	if v == nil {
		return nil, nil
	}

	certificateFiles, ok := v.([]certificateFile)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", []certificateFile{}, v)
	}

	return certificateFiles, nil
}
