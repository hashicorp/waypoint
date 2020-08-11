import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Tabs',
  component: 'Tabs',
};

export let TabsWithAnchors = () => ({
  template: hbs`
    <Tabs>
      <a class="active" href='/peering-connections'>
        Overview
      </a>
      <a href='/peering-connections'>
        Networks
      </a>
      <a class="active">
        Peering Connections
      </a>
    </Tabs>
  `,
  context: {},
});
