import { Token, NewTokenResponse } from 'waypoint-pb';
import { fakeId, fakeComponentForKind } from '../utils';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { subMinutes } from 'date-fns';

function createToken(): Token {
  let token = new Token();
  token.setAccessorId(
    'xjLoe9b2j2jTYLpM5vAv1Z6JrUW438HabDQ7fvyBzWCLozD6L2oBWRE8G2zk62V1UzcrcGf1LvmwbuQYyAFwRq63n3996WHsqrbyb8XLXfAqDbCePNX96Fkt'
  );
  return token;
}

export function create(schema: any, { params, requestHeaders }) {
  let resp = new NewTokenResponse();
  resp.setToken(createToken().getAccessorId_asB64());
  return this.serialize(resp, 'application');
}
