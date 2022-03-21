package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMGetParameterAPI defines the interface for the GetParameter function.
// We use this interface to test the function using a mocked service.
type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// SSMGetParametersByPathAPI defines the interface for the GetParameterByPath function.
// We use this interface to test the function using a mocked service.
type SSMGetParametersByPathAPI interface {
	GetParametersByPath(ctx context.Context,
		params *ssm.GetParametersByPathInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
}
