// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStore":       schema_pkg_apis_ssm_v1alpha1_ParameterStore(ref),
		"github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStoreSpec":   schema_pkg_apis_ssm_v1alpha1_ParameterStoreSpec(ref),
		"github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStoreStatus": schema_pkg_apis_ssm_v1alpha1_ParameterStoreStatus(ref),
	}
}

func schema_pkg_apis_ssm_v1alpha1_ParameterStore(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ParameterStore is the Schema for the parameterstores API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStoreSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStoreStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStoreSpec", "github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1.ParameterStoreStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_ssm_v1alpha1_ParameterStoreSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ParameterStoreSpec defines the desired state of ParameterStore",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_ssm_v1alpha1_ParameterStoreStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ParameterStoreStatus defines the observed state of ParameterStore",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
