import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

export const I18N_PREFIX = 'helpers.options-for-icon-badge';
export const DEFAULT_ICON_CLASS = '';
export const DEFAULT_ICON_SIZE = 'md';
export const DEFAULT_ICON_TYPE = 'help-circle-fill';

export const SOURCE_TYPE_UNKNOWN = {
  UNKNOWN: { label: `${I18N_PREFIX}.unknown` },
};

export const SOURCE_CONSUL_VARIANT_MAPPING = {
  CREATING: { label: `${I18N_PREFIX}.creating`, iconType: 'loading' },
  DELETED: {
    label: `${I18N_PREFIX}.deleted`,
    iconType: 'minus-circle-fill',
    textColorClass: 'text--failure',
  },
  DELETING: { label: `${I18N_PREFIX}.deleting`, iconType: 'loading' },
  FAILED: { label: `${I18N_PREFIX}.failed`, iconType: 'alert-circle-fill', textColorClass: 'text--failure' },
  PENDING: { label: `${I18N_PREFIX}.pending`, iconType: 'loading' },
  RESTORING: { label: `${I18N_PREFIX}.restoring`, iconType: 'loading' },
  RUNNING: {
    label: `${I18N_PREFIX}.running`,
    iconType: 'check-circle-fill',
    textColorClass: 'text--success',
  },
  UNSET: { label: `${I18N_PREFIX}.unset`, iconType: 'alert-triangle', textColorClass: 'text--warning' },
  UPDATING: { label: `${I18N_PREFIX}.updating`, iconType: 'loading' },
};

export const SOURCE_HVN_VARIANT_MAPPING = {
  CREATING: { label: `${I18N_PREFIX}.creating`, iconType: 'loading' },
  DELETED: {
    label: `${I18N_PREFIX}.deleted`,
    iconType: 'minus-circle-fill',
    textColorClass: 'text--failure',
  },
  DELETING: { label: `${I18N_PREFIX}.deleting`, iconType: 'loading' },
  FAILED: { label: `${I18N_PREFIX}.failed`, iconType: 'alert-circle-fill', textColorClass: 'text--failure' },
  STABLE: { label: `${I18N_PREFIX}.stable`, iconType: 'check-circle-fill', textColorClass: 'text--success' },
  UNSET: { label: `${I18N_PREFIX}.unset`, iconType: 'alert-triangle', textColorClass: 'text--warning' },
};

export const SOURCE_HVN_PEERING_VARIANT_MAPPING = {
  ACCEPTED: { label: `${I18N_PREFIX}.accepted`, iconType: 'loading' },
  ACTIVE: { label: `${I18N_PREFIX}.active`, iconType: 'check-circle-fill', textColorClass: 'text--success' },
  CREATING: { label: `${I18N_PREFIX}.creating`, iconType: 'loading' },
  DELETING: { label: `${I18N_PREFIX}.deleting`, iconType: 'loading' },
  EXPIRED: { label: `${I18N_PREFIX}.expired`, iconType: 'alert-triangle', textColorClass: 'text--warning' },
  FAILED: { label: `${I18N_PREFIX}.failed`, iconType: 'alert-triangle', textColorClass: 'text--warning' },
  PENDING_ACCEPTANCE: { label: `${I18N_PREFIX}.pending_acceptance`, iconType: 'loading' },
  REJECTED: { label: `${I18N_PREFIX}.rejected`, iconType: 'minus-square-fill' },
  UNSET: { label: `${I18N_PREFIX}.unset`, iconType: 'alert-triangle', textColorClass: 'text--warning' },
};

export const SOURCE_REGION_VARIANT_MAPPING = {
  AWS: { iconSize: 'md', iconType: 'logo-aws-color' },
};

export const SOURCES = {
  CONSUL: SOURCE_CONSUL_VARIANT_MAPPING,
  HVN: SOURCE_HVN_VARIANT_MAPPING,
  HVN_PEERING: SOURCE_HVN_PEERING_VARIANT_MAPPING,
  REGION: SOURCE_REGION_VARIANT_MAPPING,
};

export const SOURCE_SCALE = Object.keys(SOURCES);
export const VARIANT_SCALE = Object.values(SOURCES).reduce((variants, variantMap) => {
  return [...variants, ...Object.keys(variantMap)];
}, []);

export function optionForIconBadge([source = '', variant = '']) {
  let selected = SOURCES[source.toUpperCase()] || {};
  let option = selected[variant.toUpperCase()] || SOURCE_TYPE_UNKNOWN;

  if (source) {
    assert(
      `@source must be one of the following: ${SOURCE_SCALE.join(', ')}, receieved: ${source}`,
      SOURCE_SCALE.includes(source.toUpperCase())
    );
  }

  return {
    iconClass: DEFAULT_ICON_CLASS,
    iconSize: DEFAULT_ICON_SIZE,
    iconType: DEFAULT_ICON_TYPE,
    ...option,
  };
}

export default buildHelper(optionForIconBadge);
