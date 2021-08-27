import { ITheme } from 'xterm';

const lightTheme: ITheme = {
  foreground: 'rgb(0,0,0)',
  background: 'rgb(247,250,252)',
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

let terminalTheme = {
  light: lightTheme,
  // TODO: dark mode: {}
};

export default terminalTheme;
