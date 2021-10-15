export const ALERT_BANNER_ACTIVE = true

// https://github.com/hashicorp/web-components/tree/master/packages/alert-banner
export default {
  tag: 'new',
  url: 'https://www.hashicorp.com/blog/announcing-hashicorp-waypoint-0-6 ',
  text:
    'Waypoint 0.6 is available to download! Check out the announcement blog for more details.',
  linkText: 'View Post',
  // Set the `expirationDate prop with a datetime string (e.g. `2020-01-31T12:00:00-07:00`)
  // if you'd like the component to stop showing at or after a certain date
  expirationDate: '',
}
