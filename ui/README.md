# waypoint

This README outlines the details of collaborating on this Ember application.
A short introduction of this app could easily go here.

## Prerequisites

You will need the following things properly installed on your computer.

- [Node.js v14](https://nodejs.org/)
  - The current codebase has been tested to run well with `node` version 14 so it is _**strongly recommended**_ that you use this version
  - You can use node version managers to manage all of your node versions, for example [nvm](https://github.com/nvm-sh/nvm), [n](https://github.com/tj/n), etc.
- [Yarn](https://classic.yarnpkg.com/en/docs/install)
- [Ember CLI](https://ember-cli.com/)
- [Google Chrome](https://google.com/chrome/)

## Installation

- `cd waypoint/ui`
- `yarn install`

## Testing UI Changes

If you are trying out UI changes on a pull request and don't want to run Ember
locally, you can build the static assets and compile it into Waypoint server.

Run the following commands to build the UI and compile it into Waypoint Server:

```shell
cd ui
make
cd ..
make static-assets
make docker/server
```

Then when that finishes, you will be able to install the locally built server
into Docker or another platform to try out the UI changes.

## Running / Development

There are two modes of development.

### Running with Mocks

This returns data in-browser with Mirage.js (a mocking framework)
active. This means that the network requests will be intercepted
and return [mocked objects](https://github.com/hashicorp/waypoint/tree/master/ui/mirage/services)
that are static and are re-loaded on page refresh.

- `ember serve`
- The app will be available at [http://localhost:4200](http://localhost:4200).
- When prompted for a token, you can use any non-empty string (i.e. `my-cool-token`)

Troubleshooting:

- If you run into issues with `ember serve`, try deleting the directory `ui/node_modules/`, rerunning `yarn install`, and rerunning an `ember serve`.

### Running with a local Waypoint Server

This option assumes there is a Waypoint server running
at `https://localhost:9702`, which you can verify by visiting https://localhost:9702 in the browser. 

If you need to make any API changes to go along
with frontend changes, or just wish to run the server locally, you can follow
the instructions to run [Waypoint server locally](https://www.waypointproject.io/docs/server/run).

- Visit https://localhost:9702, and accept the invalid certificate warning.
- `ember serve local` 
- The app will be available at [http://localhost:4200](http://localhost:4200). Make sure that you are in the same browser session (e.g. a new tab) where you accepted the invalid certificate warning above.
- When prompted for a token, run `waypoint user token` in the command line, and enter the response. 

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
- Install [`protoc` 3.17.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3)
- Install `ts-protoc-gen`: `yarn global add ts-protoc-gen` or `npm i -g ts-protoc-gen`
- Install `protoc-gen-grpc-web`: `brew install protoc-gen-grpc-web`

#### Generate the API definitions

- `make docker/tools`
- `make docker/gen/server`
- `make gen/ts`

### Code Generators

Make use of the many generators for code, try `ember help generate` for more details

### Running Tests

- `ember test`
- `ember test --server`

(See “Percy” section for other ways of running the test suite)

### Linting

- `npm run lint:hbs`
- `npm run lint:js`
- `npm run lint:js -- --fix`

### Percy

We use [Percy](https://percy.io) for visual regression testing.

All the Percy snapshotting happens in [percy-test.ts](./tests/acceptance/percy-test.ts). The aim is to have a test for every significant state in the UI in this file. We keep it all in one file, rather than weaving Percy snapshotting through the rest of the test suite. We think this makes it more maintainable.

We are incrementally adding Percy tests, so it’s rather minimal at the moment. If you’d like to add a Percy test, please go ahead.

To run tests with Percy enabled (that is, Percy ready to receive snapshots), run the following:

```sh
yarn ember:test:percy
```

This is exactly the same command we run in CI. You will need to set the env var `PERCY_TOKEN` with a valid Percy token.

#### Percy troubleshooting

- If you need access to our Percy account (for approvals), please ask someone from @hashicorp/waypoint-frontend.
- Percy should only trigger visual diffs for changes to the UI. If you notice an unexpected Percy snapshot, that part of the UI may need the class `hide-in-percy` added to the Mirage test. It's safe to approve the Percy snapshot and tag @hashicorp/waypoint-frontend to fix the frontend Mirage tests in a separate PR.

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
