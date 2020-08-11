import headerPageObject from './nav/header';
import sectionPageObject from './nav/section';
import subheaderPageObject from './nav/subheader';

let containerSelector = '[ data-test-sidebar-nav-container ]';
export default {
  containerSelector,
  headerSelector: headerPageObject.containerSelector,
  sectionSelector: sectionPageObject.containerSelector,
  subheaderSelector: subheaderPageObject.containerSelector,
};
