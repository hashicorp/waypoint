import { Serializer } from 'ember-cli-mirage';
import * as jspb from 'google-protobuf';
import { Response } from 'miragejs';
import { encode } from '../helpers/protobufs';

export default class ApplicationSerializer extends Serializer {
  serialize(msg: jspb.Message): Response {
    return new Response(200, {}, encode(msg));
  }
}
