// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package env

import (
	"os"
	"testing"
)

func TestGetBool(t *testing.T) {
	envVarTestKey := "WAYPOINT_GET_ENV_BOOL_TEST"

	tests := []struct {
		name       string
		defaultVal bool
		envVal     string
		want       bool
		wantErr    bool
	}{
		{
			name:       "Empty env var returns default 1",
			defaultVal: true,
			envVal:     "",
			want:       true,
			wantErr:    false,
		},
		{
			name:       "Empty env var returns default 2",
			defaultVal: false,
			envVal:     "",
			want:       false,
			wantErr:    false,
		},
		{
			name:       "Non-truthy env var returns err",
			defaultVal: false,
			envVal:     "unparseable",
			want:       false,
			wantErr:    true,
		},
		{
			name:       "'true' is true",
			defaultVal: false,
			envVal:     "true",
			want:       true,
			wantErr:    false,
		},
		{
			name:       "'false' is true",
			defaultVal: true,
			envVal:     "false",
			want:       false,
			wantErr:    false,
		},
		{
			name:       "1 is true",
			defaultVal: false,
			envVal:     "1",
			want:       true,
			wantErr:    false,
		},
		{
			name:       "0 is false",
			defaultVal: true,
			envVal:     "0",
			want:       false,
			wantErr:    false,
		},
		{
			name:       "Boolean parsing ignores capitalization",
			defaultVal: false,
			envVal:     "tRuE",
			want:       true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(envVarTestKey, tt.envVal)
			got, err := GetBool(envVarTestKey, tt.defaultVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}
