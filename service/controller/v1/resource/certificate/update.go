package certificate

import (
	"context"
	"fmt"
	"path"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource/crud"
	"github.com/spf13/afero"
)

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := crud.NewPatch()
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentCertificateFiles, err := toCertificateFiles(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desiredCertificateFiles, err := toCertificateFiles(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if currentCertificateFiles == nil || desiredCertificateFiles == nil {
		return nil, nil
	}

	if !reflect.DeepEqual(currentCertificateFiles, desiredCertificateFiles) {
		r.logger.LogCtx(ctx, "debug", "current certificates do not match desired certificates, need to update")
		return desiredCertificateFiles, nil
	}

	r.logger.LogCtx(ctx, "debug", "current certificates match desired certificates, no update needed")

	return nil, nil
}

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	updateCertificateFiles, err := toCertificateFiles(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	// In case the update state is nil, don't process at all.
	if updateCertificateFiles == nil {
		return nil
	}

	// Write the update state certificate files.
	for _, fileToWrite := range updateCertificateFiles {
		r.logger.LogCtx(ctx, "debug", fmt.Sprintf("writing certificate: %s", fileToWrite.path))
		if err := afero.WriteFile(r.fs, fileToWrite.path, []byte(fileToWrite.data), r.certPermission); err != nil {
			return microerror.Mask(err)
		}
	}

	// Remove any unwanted certificate files.
	fileInfos, err := afero.ReadDir(r.fs, r.certDirectory)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, fileInfo := range fileInfos {
		fileDesired := false
		filePath := path.Join(r.certDirectory, fileInfo.Name())

		for _, updateCertificateFile := range updateCertificateFiles {
			if filePath == updateCertificateFile.path {
				fileDesired = true
			}
		}

		if !fileDesired {
			r.logger.LogCtx(ctx, "debug", fmt.Sprintf("removing certificate: %s", filePath))
			if err := r.fs.Remove(filePath); err != nil {
				return microerror.Mask(err)
			}
		}
	}

	r.logger.LogCtx(ctx, "debug", "certificates have been updated, requesting reload")
	r.prometheusReloader.RequestReload(ctx)

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
