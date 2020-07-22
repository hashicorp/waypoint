import { Response } from "miragejs";
import { Build, ListBuildsResponse } from "waypoint-pb";

export default function() {
  this.namespace = "hashicorp.waypoint.Waypoint";
  this.urlPrefix = "http://localhost:1235";
  this.timing = 0;
  this.logging = true;

  this.pretender.prepareHeaders = (headers) =>{
    headers['Content-Type'] = 'application/grpc-web-text';
    headers['X-Grpc-Web'] = '1';
    return headers;
  };

  this.post("/ListBuilds", function(schema, request) {
    let build = new Build()
    build.setId("foobar")
    let resp = new ListBuildsResponse()
    let builds = new Array(build);
    resp.setBuildsList(builds);
    
    let serialized = resp.serializeBinary()
    var len = serialized.length;
    var bytesArray = [0, 0, 0, 0];
    var payload = new Uint8Array(5 + len);
    for (var i = 3; i >= 0; i--) {
      bytesArray[i] = (len % 256);
      len = len >>> 8;
    }
    payload.set(new Uint8Array(bytesArray), 1);
    payload.set(serialized, 5);

    return new Response(
      200,
      {},
      btoa(String.fromCharCode(...payload))
    );
  });

  this.passthrough()
}
