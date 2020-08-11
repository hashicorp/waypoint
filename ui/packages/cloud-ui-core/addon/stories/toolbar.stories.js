import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Toolbar',
  component: 'Toolbar',
};

// add stories by adding more exported functions
export let Toolbar = () => ({
  template: hbs`
  <Toolbar as |t|>
    <t.Filters>
      <input type="search" />
    </t.Filters>
    <t.Actions>
      <button> An Action </button>
      <a href="">Create a thing</a>
    </t.Actions>
  </Toolbar>`,
  context: {
    // add items to the component rendering context here
  }
});

export let ToolbarWithTable = () => ({
  template: hbs`
  <Toolbar as |t|>
    <t.Filters>
      <input type="search" />
    </t.Filters>
    <t.Actions>
      <button> An Action </button>
      <a href="">Create a thing</a>
    </t.Actions>
  </Toolbar>
  <table class="pdsTable">
    <thead>
      <tr>
        <th>
        Header one
        </th>
        <th>
        Header two
        </th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>
        Data one
        </td>
        <td>
        Data two
        </td>
      </tr>
    </tbody>
  </table>
  `,
  context: {
    // add items to the component rendering context here
  }
});
