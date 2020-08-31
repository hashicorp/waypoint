import { setJSONDoc } from '@storybook/addon-docs/ember';
import json from '../dist/storybook-docgen/index.json';

setJSONDoc(json);

let rootEl = document.getElementById('root');
rootEl.style.setProperty('margin', '20px');

let docsRootEl = document.getElementById('docs-root');
docsRootEl.classList.add('pdsDocs');
