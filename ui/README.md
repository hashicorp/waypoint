# waypoint

This README outlines the details of collaborating on this Ember application.
A short introduction of this app could easily go here.

## Prerequisites

You will need the following things properly installed on your computer.

- [Node.js v12](https://nodejs.org/)
  - The current codebase has been tested to run well with `node` version 12 so it is _**strongly recommended**_ that you use this version
  - You can use node version managers to manage all of your node versions, for example [nvm](https://github.com/nvm-sh/nvm), [n](https://github.com/tj/n), etc.
- [Yarn](https://classic.yarnpkg.com/en/docs/install)
- [Ember CLI](https://ember-cli.com/)
- [Google Chrome](https://google.com/chrome/)

## Installation

- `cd waypoint/ui`
- `yarn install`

## Running / Development

There are two modes of development.

### Running with Mocks

This returns data in-browser with Mirage.js (a mocking framework)
active. This means that the network requests will be intercepted
and return [mocked objects](https://github.com/hashicorp/waypoint/tree/master/ui/mirage/services)
that are static and are re-loaded on page refresh.

- `ember serve`
- The app will be available at [http://localhost:4200](http://localhost:4200).

### Running with a local Waypoint Server

This option assumes there is a Waypoint server running
at `https://localhost:9702`. If you need to make any API changes to go along
with frontend changes, or just wish to run the server locally, you can follow
the instructions to run [Waypoint server locally](https://www.waypointproject.io/docs/server/run).

Note: You'll need to visit the above address in the same browser session to
accept the invalid certificate warning in your browser for this to work.

- `ember serve local`
- The app will be available at [http://localhost:4200](http://localhost:4200).

If you need to build the server and run it locally, you'll want to stop the existing instance, build and reinstall it in docker:

- `docker stop waypoint-server; docker rm waypoint-server; docker volume prune -f`
- `make docker/server`
- `waypoint install -platform=docker -docker-server-image="waypoint:dev" -accept-tos`

Then run the authentication steps above again.

### Generating Type Definitions after making api changes

if you've made API changes in `/internal/server` and want to use those on the frontend, you'll need to generate the type definitions again: 

#### Required dependencies for build step

- MacOS only: `brew install gnu-sed` then follow the instructions to replace the default `sed`
- Download [the 1.1.2 release of `mockery`](https://github.com/vektra/mockery/releases/tag/v1.1.2) and install in your `/go/bin` directory
- install `ts-protoc-gen`: `yarn global add ts-protoc-gen` or `npm i -g ts-protoc-gen`

#### Generate the API definitions

- `go generate ./internal/server`
- `make gen/ts`

### Code Generators

Make use of the many generators for code, try `ember help generate` for more details

### Running Tests

- `ember test`
- `ember test --server`

### Linting

- `npm run lint:hbs`
- `npm run lint:js`
- `npm run lint:js -- --fix`

### Building

- `ember build` (development)
- `ember build --environment production` (production)

### Deploying

Specify what it takes to deploy your app.

## Further Reading / Useful Links

- [ember.js](https://emberjs.com/)
- [ember-cli](https://ember-cli.com/)
- Development Browser Extensions
  - [ember inspector for chrome](https://chrome.google.com/webstore/detail/ember-inspector/bmdblncegkenkacieihfhpjfppoconhi)
  - [ember inspector for firefox](https://addons.mozilla.org/en-US/firefox/addon/ember-inspector/)
