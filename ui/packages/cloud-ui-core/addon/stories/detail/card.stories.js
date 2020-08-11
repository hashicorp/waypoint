import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Detail/Card',
  component: 'DetailCard',
};

export let DetailCard = () => ({
  template: hbs`
    <Detail::Card as |DC|>
      <DC.Header>
        <Typography @variant='h3'>
          <a href="#">
            Test Link
          </a>
        </Typography>
        <Typography @component='p'>
          <IconBadge
            @source='region'
            @variant="aws"
            @label="us-west-1"
          />
        </Typography>
      </DC.Header>
      <DC.Content>
        <Typography @component='p'>
          <IconBadge @source='hvn' @variant="STABLE" />
        </Typography>
      </DC.Content>
    </Detail::Card>
  `,
  context: {},
});