package main

import (
	"context"
	"regexp"
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

func TestRunFunction(t *testing.T) {
	var ctx context.Context
	req := &fnv1beta1.RunFunctionRequest{}

	f := &Function{log: logging.NewNopLogger()}
	want := regexp.MustCompile("since XR created")

	rsp, err := f.RunFunction(ctx, req)
	if !want.MatchString(rsp.GetResults()[0].Message) || err != nil {
		t.Fatalf("Function must output time elapsed since XR creation")
	}
}
