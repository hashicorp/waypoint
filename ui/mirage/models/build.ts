import { Model, belongsTo, hasMany } from 'ember-cli-mirage';
import { Build } from 'waypoint-pb';

export default Model.extend({
  application: belongsTo(),
  workspace: belongsTo(),
  component: belongsTo({ inverse: 'owner' }),
  status: belongsTo({ inverse: 'owner' }),
  deployments: hasMany(),
  pushedArtifact: belongsTo({ inverse: 'build' }),

  toProtobuf(): Build {
    let result = new Build();

    result.setApplication(this.application?.toProtobufRef());
    // TODO: result.setArtifact(...)
    result.setComponent(this.component?.toProtobuf());
    // TODO: result.setExtension(...)
    result.setId(this.id);
    result.setJobId(this.JobId);
    result.setSequence(this.sequence);
    result.setStatus(this.status?.toProtobuf());
    result.setTemplateData(this.templateData);
    result.setWorkspace(this.workspace?.toProtobufRef());

    for (let [key, value] of Object.entries<string>(this.labels ?? {})) {
      result.getLabelsMap().set(key, value);
    }

    return result;
  },
});
