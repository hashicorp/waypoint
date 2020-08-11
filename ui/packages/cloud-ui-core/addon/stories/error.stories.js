import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Error',
  component: 'Error',
};

export let Error = () => ({
  template: hbs`
    <Error>
      <:title>Not Found</:title>
      <:subtitle>Error 404</:subtitle>
      <:content>Some content message</:content>
      <:footer>
        <LinkTo @route="cloud">
          <Icon @type='chevron-left' @size='sm' aria-hidden='true' />
          Go back
        </LinkTo>
      </:footer>
    </Error>
  `,
  context: {}
});

export let ErrorWithIconType = () => ({
  template: hbs`
    <Error @iconType="help-circle-outline">
      <:title>Not Found</:title>
      <:subtitle>Error 404</:subtitle>
      <:content>Some content message</:content>
      <:footer>
        <LinkTo @route="cloud">
          <Icon @type='chevron-left' @size='sm' aria-hidden='true' />
          Go back
        </LinkTo>
      </:footer>
    </Error>
  `,
  context: {}
});

export let ErrorWithSize = () => ({
  template: hbs`
    <Error @size="2xl">
      <:title>Not Found</:title>
      <:subtitle>Error 404</:subtitle>
      <:content>Some content message</:content>
      <:footer>
        <LinkTo @route="cloud">
          <Icon @type='chevron-left' @size='sm' aria-hidden='true' />
          Go back
        </LinkTo>
      </:footer>
    </Error>
  `,
  context: {}
});
