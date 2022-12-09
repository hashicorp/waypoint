---
layout: docs
page_title: 'Authenticator'
description: |-
  How to implement the Authenticator component for a Waypoint plugin
---

# Authenticator

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#Authenticator

The Authenticator component is executed when the waypoint init command is called, Authenticator is enabled
by implementing two interfaces, ValidateAuthFunc and AuthFunc. Typically you implement Authenticator along
with another Component, for example you have a Platform component which deploys a waypoint application to
Google Cloud Run. You could implement the Authenticator component to check that the GCP credentials are valid.

![Authenticator](/img/extending-waypoint/authenticator.png)

The interface definition which you implement is shown below.

```go
// Authenticator is responsible for authenticating different types of plugins.
type Authenticator interface {
  // AuthFunc should return the method for getting credentials for a
  // plugin. This should return AuthResult.
  AuthFunc() interface{}
  // ValidateAuthFunc should return the method for validating authentication
  // credentials for the plugin
  ValidateAuthFunc() interface{}
}
```

ValidateAuthFunc is called when you run waypoint init, this is where you would implement logic which checks that
the plugin has the correct requirements in order to perform its work.

The signature for the function you return from AuthFunc has a single output parameter which is an error. If a non-nil
error is returned then Waypoint calls the AuthFunc method. An example implementation of VaultAuthFunc can be found
in the example below.

```go
func (p *Deploy) ValidateAuthFunc() interface{} {
  return p.validateAuth
}

func (p *Deploy) validateAuth(
  ctx context.Context,
  log hclog.Logger,
  ui terminal.UI,
) error {
  s := ui.Status()
  defer s.Close()

  s.Update("Validate authentication")

  // checkLogin returns an error when user is not
  // authenticated
  err := checkLogin()

  // returning an error from ValidateAuthFunc causes Waypoint
  // to call AuthFunc
  return err
}
```

AuthFunc is only called when ValidateAuthFunc returns an error, this is where you would implement any prompts to the
user to authenticate or where you can attempt to authenticate.

The signature for an AuthFunc has two output parameters, \*component.AuthResult and an error. If authentication succeeds
you return an AuthResult message which has Authenticated set to true &component.AuthResult{Authenticated: true}, and for
failed authentication set this to false. If an error occurs during the authentication process you can return this as the
second output parameter.

A simple example of an AuthFunc implementation can be seen below.

```go
func (p *Deploy) AuthFunc() interface{} {
  return p.authenticate
}

func (p *Deploy) authenticate(
  ctx context.Context,
  log hclog.Logger,
  ui terminal.UI,
) (*component.AuthResult, error) {
  ui.Output("Describe the manual authentication steps here")
  return &component.AuthResult{false}, nil
}
```
