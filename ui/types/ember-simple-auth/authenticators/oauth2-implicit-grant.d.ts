interface parseResponseObject {
  authuser: string;
  prompt: string;
  scope: string;
  code: string;
  state: string;
}
declare module 'ember-simple-auth/authenticators/oauth2-implicit-grant' {
  export default class OAuth2ImplicitGrantAuthenticator {}
  export function parseResponse(url: string): parseResponseObject;
}
