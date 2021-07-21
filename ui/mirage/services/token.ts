import { Token, NewTokenResponse } from 'waypoint-pb';
import { Response } from 'miragejs';

function createToken(): Token {
  let token = new Token();
  token.setTokenId(
    'bM152PWkXxfoy4vA51JFhR7LodiDkeSXVYEFiP2ShC1phS8BEjfNjiwqD1yJ17Pwz2DmkDDg2xJS8tTtGSZ5PrPLqaG5Fo2vHKdev'
  );
  return token;
}

export function create(): Response {
  let resp = new NewTokenResponse();
  resp.setToken(createToken().getTokenId_asB64());
  return this.serialize(resp, 'application');
}
