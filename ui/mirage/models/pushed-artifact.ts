import { Model, belongsTo } from 'miragejs';
import { PushedArtifact } from 'waypoint-pb';

export default Model.extend({
  application: belongsTo(),
  build: belongsTo({ inverse: 'pushedArtifact' }),
  component: belongsTo({ inverse: 'owner' }),
  status: belongsTo({ inverse: 'owner' }),
  workspace: belongsTo(),

  toProtobuf(): PushedArtifact {
    let result = new PushedArtifact();

    result.setApplication(this.application?.toProtobufRef());
    // TODO: result.setArtifact
    result.setBuild(this.build?.toProtobuf());
    result.setBuildId(this.build?.id);
    result.setComponent(this.component?.toProtobuf());
    result.setId(this.id);
    result.setJobId(this.jobId);
    result.setSequence(this.sequence);
    result.setStatus(this.status?.toProtobuf());
    result.setTemplateData(this.templateData);
    result.setWorkspace(this.workspace?.toProtobufRef());

    return result;
  },
});
