// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package appconfig provides the logic for watching and reading application
// configuration values. Application configuration values may be static
// or they may be dynamically loaded from external systems such as Vault,
// Kubernetes, AWS SSM, etc.
//
// Application configuration is primarily used by the entrypoint with
// Waypoint to load (and reload) configuration values. However, this package
// can be used standalone if application config wants to be pulled in anywhere.
//
// This package is also used for getting dynamic config values and rendering
// them into other stanzas in the waypoint.hcl.
package appconfig
