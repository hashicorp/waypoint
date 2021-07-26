# Architecture

This document describes the high-level architecture of Waypoint. The goal
of this document is to give anyone interested in contributing to Waypoint
a high-level overview of how Waypoint works. This tries to stay high level
to guide to the right area of Waypoint, but stops short of explaining in detail
how something might work since its easy for that to get out of sync with the
codebase. The recommended usage for this doc is to use it to hone in on a
specific part of Waypoint, then study the code there for specifics.

## Bird's Eye View

<p align="center"><img src=".github/images/birds-eye.png" width="70%"></p>

At the highest level, Waypoint is made up of _four_ distinct Waypoint-specific components (colored in green above):

1. **Waypoint Client (CLI, UI, etc.)** - This is an API client focused on taking
user input and calling the correct APIs to interact with Waypoint such as
viewing a list of deployments, triggering a new deployment, etc. This communicates
only to the Waypoint server.

2. **Waypoint Server** - The main API server. This responds to all the API requests
and stores data. _This communicates directly to nothing._ The Waypoint clients, runners,
entrypoints all open an outbound connection to the server, _not_ the reverse. Therefore,
the server queues information that it then communicates to other components when they
connect.

3. **Waypoint Runners** - One or more runners which are responsible for executing 
logic for builds, deploys, etc. Importantly, these are the only components in the 
entire architecture that need access to the target platform.

4. **Waypoint Entrypoint** - An optional component a deployment may have that 
connects back to the Waypoint server for features such as application configuration,
logs, and more.

