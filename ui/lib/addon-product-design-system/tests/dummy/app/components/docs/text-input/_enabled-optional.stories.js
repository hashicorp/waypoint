import hbs from 'htmlbars-inline-precompile';
import Parent from './index.stories.js';

const CONFIG = {
  ...Parent,
  title: `${Parent.title} / Enabled / Optional`,
};

// :valid
const valid_blurred = () => hbs`
  <Docs::TextInput />
`;
const valid_hovered = () => hbs`
  <Docs::TextInput
    @hovered={{true}}
  />
`;
const valid_focused = () => hbs`
  <Docs::TextInput
    @focused={{true}}
  />
`;
const valid_focused_hovered = () => hbs`
  <Docs::TextInput
    @focused={{true}}
    @hovered={{true}}
  />
`;

// visually "invalid"
const invalid_blurred = () => hbs`
  <Docs::TextInput
    @invalid={{true}}
  />
`;
const invalid_hovered = () => hbs`
  <Docs::TextInput
    @hovered={{true}}
    @invalid={{true}}
  />
`;
const invalid_focused = () => hbs`
  <Docs::TextInput
    @focused={{true}}
    @invalid={{true}}
  />
`;
const invalid_focused_hovered = () => hbs`
  <Docs::TextInput
    @focused={{true}}
    @hovered={{true}}
    @invalid={{true}}
  />
`;

/// ------------------------------------------------------------ ///
/// Story Metadata
/// ------------------------------------------------------------ ///
valid_blurred.story = { name: '(valid) blurred' };
valid_hovered.story = { name: '(valid) :hover' };
valid_focused.story = { name: '(valid) :focus' };
valid_focused_hovered.story = { name: '(valid) :focus:hover' };

invalid_blurred.story = { name: '(invalid) blurred' };
invalid_hovered.story = { name: '(invalid) :hover' };
invalid_focused.story = { name: '(invalid) :focus' };
invalid_focused_hovered.story = { name: '(invalid) :focus:hover' };

// stories module exports
export {
  CONFIG as default,

  valid_blurred,
  valid_hovered,
  valid_focused,
  valid_focused_hovered,

  invalid_blurred,
  invalid_hovered,
  invalid_focused,
  invalid_focused_hovered,
};
