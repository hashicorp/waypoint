// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package cert has helpers for managing certificates in our services.
//
// The primary helper is the Cert struct, which manages a single TLS certificate.
// This has functions for getting a *tls.Config for creating a TLS listener,
// and automatically watches and reloads on any certificate file changes. It also
// provides functions for atomic replacement of certificates for zero downtime
// replacement.
package cert
