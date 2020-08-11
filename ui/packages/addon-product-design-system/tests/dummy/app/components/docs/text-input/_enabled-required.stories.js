import hbs from 'htmlbars-inline-precompile';
import Parent from './index.stories.js';

const CONFIG = {
  ...Parent,
  title: `${Parent.title} / Enabled / Required`,
};

/// "clean" (a.k.a., "pristine")
const clean_blurred = () => hbs`
  <Docs::TextInput
    @required={{true}}
  />
`;
const clean_hovered = () => hbs`
  <Docs::TextInput
    @hovered={{true}}
    @required={{true}}
  />
`;
const clean_focused = () => hbs`
  <Docs::TextInput
    @focused={{true}}
    @required={{true}}
  />
`;
const clean_focused_hovered = () => hbs`
  <Docs::TextInput
    @focused={{true}}
    @hovered={{true}}
    @required={{true}}
  />
`;

/// "dirty"
const dirty_blurred = () => hbs`
  <Docs::TextInput
    @dirty={{true}}
    @required={{true}}
  />
`;
const dirty_hovered = () => hbs`
  <Docs::TextInput
    @dirty={{true}}
    @hovered={{true}}
    @required={{true}}
  />
`;
const dirty_focused = () => hbs`
  <Docs::TextInput
    @dirty={{true}}
    @focused={{true}}
    @required={{true}}
  />
`;
const dirty_focused_hovered = () => hbs`
  <Docs::TextInput
    @dirty={{true}}
    @focused={{true}}
    @hovered={{true}}
    @required={{true}}
  />
`;

/// ------------------------------------------------------------ ///
/// Story Metadata
/// ------------------------------------------------------------ ///
clean_blurred.story = { name: '(clean) blurred' };
clean_hovered.story = { name: '(clean) :hover' };
clean_focused.story = { name: '(clean) :focus' };
clean_focused_hovered.story = { name: '(clean) :focus:hover' };

dirty_blurred.story = { name: '(dirty) blurred' };
dirty_hovered.story = { name: '(dirty) :hover' };
dirty_focused.story = { name: '(dirty) :focus' };
dirty_focused_hovered.story = { name: '(dirty) :focus:hover' };

// stories module exports
export {
  CONFIG as default,

  clean_blurred,
  clean_hovered,
  clean_focused,
  clean_focused_hovered,

  dirty_blurred,
  dirty_hovered,
  dirty_focused,
  dirty_focused_hovered,
};
