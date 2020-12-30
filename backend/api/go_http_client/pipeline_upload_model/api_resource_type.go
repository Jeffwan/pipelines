// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by go-swagger; DO NOT EDIT.

package pipeline_upload_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// APIResourceType api resource type
// swagger:model apiResourceType
type APIResourceType string

const (

	// APIResourceTypeUNKNOWNRESOURCETYPE captures enum value "UNKNOWN_RESOURCE_TYPE"
	APIResourceTypeUNKNOWNRESOURCETYPE APIResourceType = "UNKNOWN_RESOURCE_TYPE"

	// APIResourceTypeEXPERIMENT captures enum value "EXPERIMENT"
	APIResourceTypeEXPERIMENT APIResourceType = "EXPERIMENT"

	// APIResourceTypeJOB captures enum value "JOB"
	APIResourceTypeJOB APIResourceType = "JOB"

	// APIResourceTypePIPELINE captures enum value "PIPELINE"
	APIResourceTypePIPELINE APIResourceType = "PIPELINE"

	// APIResourceTypePIPELINEVERSION captures enum value "PIPELINE_VERSION"
	APIResourceTypePIPELINEVERSION APIResourceType = "PIPELINE_VERSION"

	// APIResourceTypeNAMESPACE captures enum value "NAMESPACE"
	APIResourceTypeNAMESPACE APIResourceType = "NAMESPACE"

	// APIResourceTypeUSER captures enum value "USER"
	APIResourceTypeUSER APIResourceType = "USER"
)

// for schema
var apiResourceTypeEnum []interface{}

func init() {
	var res []APIResourceType
	if err := json.Unmarshal([]byte(`["UNKNOWN_RESOURCE_TYPE","EXPERIMENT","JOB","PIPELINE","PIPELINE_VERSION","NAMESPACE","USER"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		apiResourceTypeEnum = append(apiResourceTypeEnum, v)
	}
}

func (m APIResourceType) validateAPIResourceTypeEnum(path, location string, value APIResourceType) error {
	if err := validate.Enum(path, location, value, apiResourceTypeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this api resource type
func (m APIResourceType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateAPIResourceTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
