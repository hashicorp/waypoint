import { Schema } from "ember-cli-mirage";
import { Build, ListBuildsResponse, Component, Status, Ref} from 'waypoint-pb';
import { fakeId, fakeComponentForKind } from '../utils';;
import { Timestamp } from "google-protobuf/google/protobuf/timestamp_pb";
import { subMinutes } from 'date-fns'

function createBuild(): Build {
    let build = new Build()
    build.setId(fakeId())

    // todo(pearkes): create util
    let workspace = new Ref.Workspace()
    workspace.setWorkspace("default")

    let component = new Component()
    component.setType(Component.Type.BUILDER)
    component.setName(fakeComponentForKind(Component.Type.BUILDER))

    // todo(pearkes): random state
    let status = new Status()
    status.setState(Status.State.SUCCESS)

    // todo(pearkes): helpers
    let timestamp = new Timestamp()
    let result = Math.floor(subMinutes(new Date(), 30).getTime() / 1000)
    timestamp.setSeconds(result)

    // Same thing for now
    status.setCompleteTime(timestamp)
    status.setStartTime(timestamp)

    build.setComponent(component)
    build.setStatus(status)
    build.setWorkspace(workspace)

    return build
}

export function list(schema: Schema, { params, requestHeaders }) {
    let resp = new ListBuildsResponse()
    let builds = new Array(
        createBuild(),
        createBuild(),
        createBuild(),
        createBuild()
    )
    resp.setBuildsList(builds);
    return this.serialize(resp, "application")
}
