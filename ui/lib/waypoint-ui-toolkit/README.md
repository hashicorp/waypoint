# Waypoint UI Toolkit

Waypoint UI Toolkit is a set of components, helpers, and utilities for the
Waypoint Web UI.

It is not general purpose, and is only intended for use on HashiCorp projects.

## Installing

Install the toolkit just like any other NPM package.

```
$ yarn add -D @hashicorp/waypoint-ui-toolkit
$ npm install -D @hashicorp/waypoint-ui-toolkit
```

## Updating

Update the toolkit just like any other NPM package, using `latest` or a specific
version reference.

```
$ yarn add -D @hashicorp/waypoint-ui-toolkit@latest
$ npm install -D @hashicorp/waypoint-ui-toolkit@latest
```

## Usage

Components in the toolkit are scoped under the `Waypoint` namespace. They are
available automatically in the consuming project:

```hbs
<Waypoint::Timeline @model={{...}} />
```

Helpers are similarly scoped, and also automatically available:

```hbs
{{waypoint/icon-for-component @operation.component.name}}
```

Utilities may be imported from the package:

```ts
import { imageRef } from '@hashicorp/waypoint-ui-toolkit';
```

TODO: Where do I find documentation on individual elements of the kit?
TODO: API.md?
TODO: Storybook?

## Reporting Issues

Please report issues on
[github.com/hashicorp/waypoint](https://github.com/hashicorp/waypoint/issues).

## Contributing

See [CONTRIBUTING.md].
