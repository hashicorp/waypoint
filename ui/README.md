# waypoint

This README outlines the details of collaborating on this Ember application.
A short introduction of this app could easily go here.

## Prerequisites

You will need the following things properly installed on your computer.

- [Git](https://git-scm.com/)
- [Node.js](https://nodejs.org/) (with npm)
- [Ember CLI](https://ember-cli.com/)
- [Google Chrome](https://google.com/chrome/)

## Installation

- `git clone <repository-url>` this repository
- `cd waypoint`
- `npm install`

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
