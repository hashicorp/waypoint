import OAuth2PasswordGrant from 'ember-simple-auth/authenticators/oauth2-password-grant';
import classic from 'ember-classic-decorator';

@classic
export default class OIDCAuthenticator extends OAuth2PasswordGrant {

};
