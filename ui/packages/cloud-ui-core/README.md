cloud-ui-core
==============================================================================

:wave: Hello! This is an Ember addon that's shared between the host application for 
HCP (cloud-ui) and all of the engines `cloud-ui` uses. Things that are shared 
at the platform level - components, styles, helpers, services, utilities - 
should go in this addon.


Installation
------------------------------------------------------------------------------

Installation outside of the `cloud-ui` monorepo is not currently supported.

To add this addon to an engine or an app _inside_ the monorepo, please make 
sure `cloud-ui-core` is in your engine's `package.json` `dependencies` block and
run `yarn install` either from the root of this project or in the addon's directory.


Usage
------------------------------------------------------------------------------

All of the components are exported into the host app. To see their Storybook 
documentation, run the app or addon's Storybook server, and they should be
included automatically.


Contributing
------------------------------------------------------------------------------

See the [Contributing](CONTRIBUTING.md) guide for details.

