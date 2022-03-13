package controllers

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/fr123k/aws-ssm-operator/api/v1alpha1"
	errs "github.com/pkg/errors"
)

// SSMClient preserves AWS client session and SSM client itself
type SSMClient struct {
	s   *session.Session
	ssm *ssm.SSM
}

func awsCfg() *aws.Config {
	if lsEp := os.Getenv("LOCAL_STACK_ENDPOINT"); lsEp != "" {
		log.Info("Setup localstack aws client")
		cfg := &aws.Config{
			Region:                 aws.String("us-east-1"),
			Credentials:            credentials.NewStaticCredentials("test", "test", ""),
			S3ForcePathStyle:       aws.Bool(true),
			Endpoint:               aws.String(lsEp),
			DisableParamValidation: aws.Bool(true),
		}
		return cfg
	}
	return nil
}

func newSSMClient(s *session.Session) *SSMClient {
	return &SSMClient{
		s: s,
	}
}

// FetchParameterStoreValue fetches decrypted values from SSM Parameter Store
func (c *SSMClient) FetchParameterStoreValues(ref v1alpha1.ParameterStoreRef) (map[string]string, error) {
	if c.s == nil {
		c.s = session.Must(session.NewSession(awsCfg()))
	}

	if c.ssm == nil {
		c.ssm = ssm.New(c.s)
	}

	if ref.Name != "" {
		log.Info("fetching values from SSM Parameter Store by name", "Name", ref.Name)
		got, err := c.ssm.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(ref.Name),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return nil, err
		}

		return map[string]string{"name": aws.StringValue(got.Parameter.Value)}, nil
	}

	log.Info("fetching values from SSM Parameter Store by path", "Path", ref.Path, "Recursive", ref.Recursive)
	got, err := c.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(ref.Path),
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(ref.Recursive),
		MaxResults:     aws.Int64(100),
	})
	if err != nil {
		return nil, err
	}

	log.Info("fetching values from SSM Parameter Store by path", "Params", fmt.Sprintf("%+v", got.Parameters))

	dict := make(map[string]string, len(got.Parameters))
	for _, p := range got.Parameters {
		ss := strings.Split(aws.StringValue(p.Name), "/")
		name := strings.ToUpper(ss[len(ss)-1])
		name = strings.ReplaceAll(name, "-", "_")
		dict[name] = aws.StringValue(p.Value)
	}

	return dict, nil
}

// SSMParameterValueToSecret shapes fetched value so as to store them into K8S Secret
func (c *SSMClient) SSMParameterValueToSecret(ref v1alpha1.ParameterStoreRef) (map[string]string, error) {
	if ref.Name != "" || ref.Path != "" {
		params, err := c.FetchParameterStoreValues(ref)
		if err != nil {
			return nil, err
		}
		if params == nil {
			return nil, errs.New("fetched value must not be nil")
		}

		return params, nil
	}
	return map[string]string{}, nil
}

func (c *SSMClient) FetchParametersStoreValues(refs []v1alpha1.ParametersStoreRef) (map[string]string, map[string]string, error) {
	if c.s == nil {
		c.s = session.Must(session.NewSession(awsCfg()))
	}

	if c.ssm == nil {
		c.ssm = ssm.New(c.s)
	}

	dict := make(map[string]string)
	anno := make(map[string]string)

	for _, ref := range refs {
		log.Info("fetching values from SSM Parameter Store", "Key", ref.Key, "Name", ref.Name)
		got, err := c.ssm.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(ref.Key),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Error(err, "error fetching values from SSM Parameter Store", "Key", ref.Key, "Name", ref.Name)
			anno[fmt.Sprintf("ssm.aws/%s_error", ref.Name)] = err.Error()
			continue
		}
		name := ref.Name
		if name == "" {
			ss := strings.Split(aws.StringValue(got.Parameter.Name), "/")
			name = strings.ToUpper(ss[len(ss)-1])
			name = strings.ReplaceAll(name, "-", "_")
		}
		dict[name] = aws.StringValue(got.Parameter.Value)
	}

	return dict, anno, nil
}

func (c *SSMClient) SSMParametersValueToSecret(ref []v1alpha1.ParametersStoreRef) (map[string]string, map[string]string, error) {
	params, anno, err := c.FetchParametersStoreValues(ref)
	if err != nil {
		return nil, nil, err
	}
	if params == nil {
		return nil, nil, errs.New("fetched value must not be nil")
	}

	return params, anno, nil
}
