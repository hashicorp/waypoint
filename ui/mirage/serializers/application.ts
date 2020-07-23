import { Serializer, Request } from 'ember-cli-mirage';
import * as jspb from 'google-protobuf'
import { Response } from "miragejs";

export default class ApplicationSerializer extends Serializer {
    serialize(model: jspb.Message, request: Request): Response {
        let resp = model
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
    }
}
