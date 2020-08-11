import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ExternalLink',
  component: 'ExternalLink',
};

export let ExternalLink = () => ({
  template: hbs`
    <ExternalLink
      href='https://www.hashicorp.com'
    >
      HashiCorp
    </ExternalLink>
  `,
});
