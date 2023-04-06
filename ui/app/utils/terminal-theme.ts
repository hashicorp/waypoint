/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { ITheme } from 'xterm';

const lightTheme: ITheme = {
  foreground: 'rgb(0,0,0)',
  background: 'rgb(247,250,252)',
  cursor: 'rgb(0, 187, 187)',
  cyan: 'rgb(0, 187, 187)',
  brightBlue: 'rgb(85, 85, 255)',
  green: 'rgb(0, 187, 0)',
  magenta: 'rgb(187, 0, 187)',
  brightMagenta: 'rgb(187, 0, 187)',
  yellow: 'rgb(187, 187, 0)',
  brightYellow: 'rgb(187, 187, 0)',
  blue: 'rgb(0, 187, 187)',
  brightBlack: 'rgb(0,0,0)',
  selection: 'rgb(5,198,194)',
};

const darkTheme: ITheme = {
  foreground: 'rgb(247,250,252)',
  background: 'rgb(29,33,38)',
  cursor: 'rgb(0, 187, 187)',
  cyan: 'rgb(0, 187, 187)',
  brightBlue: 'rgb(85, 85, 255)',
  green: 'rgb(0, 187, 0)',
  magenta: 'rgb(187, 0, 187)',
  brightMagenta: 'rgb(187, 0, 187)',
  yellow: 'rgb(187, 187, 0)',
  brightYellow: 'rgb(187, 187, 0)',
  blue: 'rgb(0, 187, 187)',
  brightBlack: 'rgb(247,250,252)',
  selection: 'rgb(5,198,194)',
};

let terminalTheme = {
  light: lightTheme,
  dark: darkTheme,
};

export default terminalTheme;
