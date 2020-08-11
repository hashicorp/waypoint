import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'FlexGrid',
  component: 'FlexGrid',
};

export let FlexGrid = () => ({
  template: hbs`
    <FlexGrid as |G|>
      <G.Item @xs={{1}}>1</G.Item>
      <G.Item @xs={{10}}>10</G.Item>
      <G.Item @xs={{1}}>1</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{1}}>1</G.Item>
      <G.Item @xs={{10}}>10</G.Item>
      <G.Item @xs={{1}}>1</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{10}}>10</G.Item>
      <G.Item @xs={{2}}>2</G.Item>
    </FlexGrid>
  `,
  context: {},
});

export let FlexGridCascade = () => ({
  template: hbs`
    <FlexGrid as |G|>
      <G.Item @xs={{12}}>12</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{1}}>1</G.Item>
      <G.Item @xs={{11}}>11</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{2}}>2</G.Item>
      <G.Item @xs={{10}}>10</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{3}}>3</G.Item>
      <G.Item @xs={{9}}>9</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{4}}>4</G.Item>
      <G.Item @xs={{8}}>8</G.Item>
    </FlexGrid>
  `,
  context: {},
});

export let FlexGridOffset = () => ({
  template: hbs`
    <FlexGrid as |G|>
      <G.Item @xs={{1}} @xsOffset={{11}}>1</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{3}} @xsOffset={{9}}>3</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{5}} @xsOffset={{7}}>5</G.Item>
    </FlexGrid>
    <FlexGrid as |G|>
      <G.Item @xs={{7}} @xsOffset={{5}}>7</G.Item>
    </FlexGrid>
  `,
  context: {},
});

export let FlexGridReverse = () => ({
  template: hbs`
    <FlexGrid @reverse={{true}} as |G|>
      <G.Item @xs={{4}}>1</G.Item>
      <G.Item @xs={{4}}>2</G.Item>
      <G.Item @xs={{4}}>3</G.Item>
    </FlexGrid>
  `,
  context: {},
});
