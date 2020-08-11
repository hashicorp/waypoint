import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Select',
  component: 'Select',
};

// add stories by adding more exported functions
export let basic = () => ({
  template: hbs`<Select
    @options={{this.options}}
    @value={{this.value}}
    {{on 'change' (set this.value (get _ 'target.value'))}}
  />
  `,
  context: {
    options: ['first', 'second', 'third'],
    value: 'second'
  }
});

export let withObjects = () => ({
  template: hbs`<Select
    @options={{this.options}}
    @value={{this.value}}
    @valuePath='ordinal'
    @labelPath='number'
    {{on 'change' (set this.value (get _ 'target.value'))}}
  />
  `,
  context: {
    options: [
      {number: '1', ordinal: 'first'},
      {number: '2', ordinal:'second'},
      {number: '3', ordinal: 'third'},
    ],
    value: 'second',
  }
});
