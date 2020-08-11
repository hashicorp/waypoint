import LinkComponent from '@ember/routing/link-component';

/**
 *
 * `CtaLink` extends the LinkTo component and can be used exactly the same way.
 *
 *
 * ```
 * <CtaLink @route="some.route">Some Link</CtaLink>
 * ```
 *
 * @class CtaLink
 *
 */

export default class CtaLinkComponent extends LinkComponent {
  classNames = ['cta-link'];
}
