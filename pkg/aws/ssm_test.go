package aws

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/fr123k/aws-ssm-operator/api/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func TestLocalStackIntegration(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/")
		// Send response to be tested
		rw.Write([]byte(`{
			"Parameter": {
				"ARN": "arn:aws:ssm:us-east-2:111122223333:parameter/MyGitHubPassword",
				"DataType": "text",
				"LastModifiedDate": 1582657288.8,
				"Name": "MyGitHubPassword",
				"Type": "SecureString",
				"Value": "AYA39c3b3042cd2aEXAMPLE/AKIAIOSFODNN7EXAMPLE/fh983hg9awEXAMPLE==",
				"Version": 3
			}
		}`))
	}))
	// Close the server when test finishes
	defer server.Close()
	t.Setenv("LOCAL_STACK_ENDPOINT", server.URL)
	ssm := NewSSMClient(nil)
	result, err := ssm.SSMParameterValueToSecret(v1alpha1.ParameterStoreRef{Name: "MyGitHubPassword"})

	fmt.Printf("%+v", err)

	assert.Nil(t, err)
	assert.Equal(t, "AYA39c3b3042cd2aEXAMPLE/AKIAIOSFODNN7EXAMPLE/fh983hg9awEXAMPLE==", result["MyGitHubPassword"])
}

type SSMGetParameterImpl struct{}

func (dt SSMGetParameterImpl) GetParameter(ctx context.Context,
	params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {

	parameter := &types.Parameter{Name: aws.String("name"), Value: aws.String("aws-docs-example-parameter-value")}

	output := &ssm.GetParameterOutput{
		Parameter: parameter,
	}

	return output, nil
}

func TestFetchParameterStoreValues(t *testing.T) {
	ssm := NewSSMClient(nil)
	ssm.SSMGetParameterAPI = &SSMGetParameterImpl{}
	result, err := ssm.SSMParameterValueToSecret(v1alpha1.ParameterStoreRef{Name: "name"})

	assert.Nil(t, err)
	assert.Equal(t, "aws-docs-example-parameter-value", result["name"])
}

type SSMGetParametersByPathImpl struct{}

func (dt SSMGetParametersByPathImpl) GetParametersByPath(ctx context.Context,
	params *ssm.GetParametersByPathInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {

	parameter := types.Parameter{Name: aws.String("name"), Value: aws.String("aws-docs-example-parameter-value")}

	output := &ssm.GetParametersByPathOutput{
		Parameters: []types.Parameter{parameter},
	}

	return output, nil
}

func TestFetchParametersStoreValues(t *testing.T) {
	ssm := NewSSMClient(nil)
	ssm.SSMGetParameterAPI = &SSMGetParameterImpl{}
	result, anno, err := ssm.SSMParametersValueToSecret([]v1alpha1.ParametersStoreRef{
		v1alpha1.ParametersStoreRef{
			Name: "NAME",
			Key:  "name",
		}},
	)

	assert.Nil(t, err)
	assert.Len(t, anno, 0)
	assert.Equal(t, "aws-docs-example-parameter-value", result["NAME"])
}
