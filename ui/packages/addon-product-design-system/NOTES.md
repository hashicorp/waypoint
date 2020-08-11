## pds/modules

### pds/modules/layout
See [every-layout.dev](https://every-layout.dev)

* division
  * stack
  * rows
  * grid
  * tiles
  * center
  * etc.
* alignment
  * top
  * middle
  * bottom
  * left
  * center
  * right
  * start
  * end
  * etc.

## Font Stacks

### SCSS Variables

Be careful to avoid collision with the following pre-existing variables.

variable             | Bulma | Nomad | Vault | Consul | Atlas
-------------------- | ----- | ----- | ----- | ------ | -----
`$body-family`       | x     |       |       |        |
`$code-family`       | x     |       |       |        |
`$family-code`       | x     |       |       |        |
`$family-monospace`  | x     |       |       |        |
`$family-primary`    | x     |       | x     |        |
`$family-sans-serif` | x     | x     |       |        | x
`$family-sans`       |       |       | x     |        |
`$family-secondary`  | x     |       |       |        |
`$typo-family-mono`  |       |       |       | x      |
`$typo-family-sans`  |       |       |       | x      |


### System Font Stack (sans-serif)

**TL;DR**: Bulma's definition will be used as a base because it contains the most
ubiquitous values.  However, we'll use `system-ui` instead of
`BlinkMacSystemFont, -apple-system`.

#### Why use `system-ui`?
The `system-ui` keyword works the same as `BlinkMacSystemFont, -apple-system`,
so it allows us to simplify the font-family definition.

* WebKit/Safari has supported `system-ui` [since 2017](https://bugs.webkit.org/show_bug.cgi?id=151493).
* Chrome [supports `system-ui`](https://www.chromestatus.com/feature/5640395337760768) as of v56.
* Edge (pre-chromium) doesn't matter, because it doesn't run on Mac.
* Firefox seems to support it, but there
  [might be quirks](https://bugzilla.mozilla.org/show_bug.cgi?id=1226042)
  in the internal font resolution algorithm.
    * The browser still internally resolves to one of the fallback fonts, but it
      may be less predictable than expected (may be only related to Linux).
    * We only care that the font resolves to either the OS-configured UI font _or_
      one of the listed fallback fonts. Even with the quirky behavior, that goal
      is still achieved.


#### Analysis of Product UI Font Stacks
- Consul matches Bulma
- Nomad and Atlas Match
- Vault is unique
- GitHub is unique


value               | Bulma | Nomad | Vault | Consul | Atlas | GitHub
------------------- | ----- | ----- | ----- | ------ | ----- | ------
system-ui,          |       |       | x     |        |       |
-apple-system,      | x     | x     | x     | x      | x     | x
BlinkMacSystemFont, | x     | x     | x     | x      | x     | x
'Segoe UI',         | x     | x     | x     | x      | x     | x
'Roboto',           | x     | x     | x     | x      | x     |
"Oxygen",           | x     |       | x     | x      |       |
Oxygen-sans,        |       | x     |       |        | x     |
"Ubuntu",           | x     | x     | x     | x      | x     |
"Cantarell",        | x     | x     | x     | x      | x     |
"Fira Sans",        | x     |       | x     | x      |       |
"Droid Sans",       | x     |       | x     | x      |       |
"Helvetica Neue",   | x     | x     | x     | x      | x     |
"Helvetica",        | x     |       |       | x      |       | x
"Arial",            | x     |       |       | x      |       | x
sans-serif          | x     | x     | x     | x      | x     | x
Apple Color Emoji   |       |       |       |        |       | x
Segoe UI Emoji      |       |       |       |        |       | x

```scss
/* Bulma */
$family-sans-serif: BlinkMacSystemFont, -apple-system, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", "Helvetica", "Arial", sans-serif;
$family-primary: $family-sans-serif;
$family-primary: $family-sans-serif;
$family-secondary: $family-sans-serif;
$body-family: $family-primary;

/* Nomad */
$family-sans-serif: -apple-system, BlinkMacSystemFont,'Segoe UI', Roboto, Oxygen-Sans, Ubuntu, Cantarell, 'Helvetica Neue', sans-serif;
/* Vault */
$family-sans: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
$family-primary: $family-sans;
/* Consul (doesn't use Bulma) */
$typo-family-sans: BlinkMacSystemFont, -apple-system, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', 'Helvetica', 'Arial', sans-serif;
/* Atlas */
$family-sans-serif: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif;

/* GitHub */
font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji;
```

- [Bulma variables](https://bulma.io/documentation/customize/variables/)
- Nomad [SCSS config](https://github.com/hashicorp/nomad/blob/master/ui/app/styles/core/variables.scss#L27-L28)
- Nomad uses Bulma and overrides `$family-sans-serif`
- I like the set of vars (weight, size, leading, etc.) that are defined in Nomad.
- Vault [SCSS config](https://github.com/hashicorp/vault/blob/master/ui/app/styles/utils/_bulma_variables.scss#L26-L29)
- Vault uses Bulma and overrides `$family-primary` with local `$family-sans` variable.
- Consul [SCSS config](https://github.com/hashicorp/consul/blob/master/ui-v2/app/styles/base/typography/base-variables.scss#L1-L3)
- Consul doesn't use Bulma
- Atlas [SCSS config](https://github.com/hashicorp/atlas/blob/master/frontend/atlas/app/styles/_variables.scss#L194-L195)
- Atlas uses Bulma and overrides `$family-sans-serif`
- GitHub fonts were determined by searched for `font-family` on their production stylesheet via browser dev tools.



### Monospace Font Stack
**TL;DR** - `monospace` should suffice until data says otherwise.
No need to overcomplicate it right now.

```scss
/* Bulma */
$family-monospace: monospace;
$family-code: $family-monospace;
$code-family: $family-code;

/* Vault */
$family-monospace: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;

/* GitHub (2020-04-02) */
font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace;
```

- Vault uses Bulma and overrides `$family-monospace`
