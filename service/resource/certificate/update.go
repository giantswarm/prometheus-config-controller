package certificate

import (
	"context"
	"reflect"

	"k8s.io/client-go/pkg/api/v1"

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
		return nil, nil, &desiredCertificateFiles, nil
	}

	r.logger.Log("debug", "current certificates matches desired certificates, no update needed")

	return nil, nil, nil, nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	return nil
}

func toCertificateFiles(v interface{}) ([]certificateFile, error) {
	if v == nil {
		return nil, nil
	}

	certificateFiles, ok := v.([]certificateFile)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", v1.ConfigMap{}, v)
	}

	return certificateFiles, nil
}
