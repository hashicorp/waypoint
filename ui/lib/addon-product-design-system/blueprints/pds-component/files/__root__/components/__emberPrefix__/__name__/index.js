import Component from '@glimmer/component';
/**<% if (!dummy) { %>
 *
 * `<%= tagName %>` description here.
 *
 *
 * ```hbs
 * <<%= tagName %> />
 * ```
 *
 *<% } %>
 * @class <%= jsClass %>
 */
export default class <%= jsClass %> extends Component {
<% if (!dummy) { %>
  /**
   * Is component awesome?
   *
   * @argument isAwesome
   * @type { boolean }
   */
  isAwesome = true;
<% } %>
}
