import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ModalDeleteConfirm',
  component: 'ModalDeleteConfirm',
};

// add stories by adding more exported functions
export let ModalDeleteConfirm = () => ({
  template: hbs`<ModalDeleteConfirm />`,
  context: {
    // add items to the component rendering context here
  }
});
