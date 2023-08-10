// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package virtualceb is used to provide "virtual" CEB functionality. A
// virtual CEB acts like an entrypoint but doesn't represent a real physical
// instance of a deployment, hence the "virtual" labeling.
//
// This functionality is used in situations where a real entrypoint either
// can't run or is impractical to run. Most commonly this is used for serverless
// type environments such as Lambda. In those scenarios, it is impractical to
// run the entrypoint for tasks such as exec and logs since the instances are
// so ephemeral.
package virtualceb
