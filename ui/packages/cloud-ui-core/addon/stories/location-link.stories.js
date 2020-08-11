import hbs from 'htmlbars-inline-precompile';
import {
  HASHICORP_LINK_CONSUL_CLUSTER,
} from 'cloud-ui/components/location-link/consts';

export default {
  title: 'LocationLink',
  component: 'LocationLink',
};

export let LocationLink = () => ({
  template: hbs`
    <LocationLink @link={{link}} />
  `,
  context: {
    link: {
      location: {
        projectId: 'project_id',
      },
      type: HASHICORP_LINK_CONSUL_CLUSTER,
      uuid: 'some-uuid',
    }
  },
});

export let LocationLinkWithName = () => ({
  template: hbs`
    <LocationLink @link={{link}} />
  `,
  context: {
    link: {
      location: {
        projectId: 'project_id',
      },
      type: HASHICORP_LINK_CONSUL_CLUSTER,
      name: 'some name',
    }
  },
});

export let LocationLinkMissing = () => ({
  template: hbs`
    <LocationLink @link={{link}} />
  `,
  context: {
    link: {
      location: {
        projectId: 'project_id',
      },
      type: 'some.unknown.type',
      uuid: 'some-uuid',
    }
  },
});
