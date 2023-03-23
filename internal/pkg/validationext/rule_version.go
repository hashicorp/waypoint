// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validationext

import (
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	goversion "github.com/hashicorp/go-version"
)

var (
	// MeetsVersionConstraint implements validation.Rule to check if a
	// supplied version is a valid according to the constraints supplied.
	MeetsVersionConstraint validation.Rule = &ConstraintVersionRule{}
	// IsVersion implements validation.Rule to check if the supplied string
	// is considered a valid Semantic Versioning (semver) string.
	IsVersion validation.Rule = &ParseVersionRule{}
)

type ParseVersionRule struct{}

func (v *ParseVersionRule) Validate(value interface{}) error {
	valueVersion, err := extractVersionString(value)
	if err != nil {
		return err
	}

	_, err = goversion.NewVersion(valueVersion)
	return err
}

type ConstraintVersionRule struct {
	Constraint goversion.Constraints
}

func MeetsConstraints(constraints ...string) validation.Rule {
	constraintString := strings.Join(constraints, ",")
	constraint, err := goversion.NewConstraint(constraintString)
	if err != nil {
		panic(err)
	}
	return &ConstraintVersionRule{
		Constraint: constraint,
	}
}

func (v *ConstraintVersionRule) Validate(value interface{}) error {
	valueVersion, err := extractVersionString(value)
	if err != nil {
		return err
	}

	goVer, err := goversion.NewVersion(valueVersion)
	if err != nil {
		return err
	}

	if !v.Constraint.Check(goVer) {
		return fmt.Errorf("%s does not satisfy constraint: %s", goVer, v.Constraint)
	}

	return nil
}

func extractVersionString(value interface{}) (string, error) {
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	switch kind {
	case reflect.String:
		return value.(string), nil
	default:
		return "", fmt.Errorf("type not supported: %v, must be string", rv.Type())
	}

}
