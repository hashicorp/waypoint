import { helper } from '@ember/component/helper';
import awsRegionList from 'cloud-ui-core/utils/aws-region-list';

export default helper(function awsRegions() {
  return awsRegionList;
});
