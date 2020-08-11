import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Menu',
  component: 'Menu',
};

// add stories by adding more exported functions
export let Menu = () => ({
  template: hbs`
    <Menu as |M|>
      <M.Trigger>
        {{#if M.isOpen}}
          This is an open menu!
        {{else}}
          Open the menu
        {{/if}}
      </M.Trigger>
      <M.Content>
        Put some content here
      </M.Content>
    </Menu>
    `,
  context: {
    // add items to the component rendering context here
  },
});
