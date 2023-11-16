package main

import (
	"context"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
)

const (
	// AnnotationKeyTTL is the key in the annotations map of a resource that indicates
	// the minimum duration composed resources should be considered valid.
	AnnotationKeyTTL = "fn.crossplane.io/ttl"
)

// error codes
const (
	errGetObserved = "cannot get observed composite resource"
	errGetTTL      = "cannot get ttl from composite resource"
	errExpired     = "composed resources have expired"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	rsp := response.To(req, response.DefaultTTL)

	// Read the observed XR from the request.
	xr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		// If the function can't read the XR, the request is malformed. This
		// should never happen. The function returns a fatal result. This tells
		// Crossplane to stop running functions and return an error.
		f.log.Debug(errGetObserved, req)
		response.Fatal(rsp, errors.Wrap(err, errGetObserved))
		return rsp, nil
	}

	// Check if the observed XR has a TTL annotation
	ttl, ok := xr.Resource.GetAnnotations()[AnnotationKeyTTL]
	if !ok {
		// no TTL, just exit
		return rsp, nil
	}

	timeout, err := time.ParseDuration(ttl)
	if err != nil {
		// if we couldn't parse the timeout from the TTL annotation, bail out with an error.
		f.log.Debug(errGetTTL, req)
		response.Fatal(rsp, errors.Wrap(err, errGetTTL))
		return rsp, nil
	}

	// Calculate the elapsed time since the Composite Resource was created
	elapsed := time.Time.Sub(time.Now(), xr.Resource.GetCreationTimestamp().Time)

	// if the elapsed time is longer than the timeout return an empty list for the desired composed resources and set a warning
	if elapsed > timeout {
		// TODO: Should we log the expired resources somewhere?
		rsp.Desired.Resources = map[string]*fnv1beta1.Resource{}
		response.Warning(rsp, errors.New(errExpired))
	}

	return rsp, nil
}
