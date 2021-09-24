import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetPushedArtifactRequest, Ref } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

interface DockerImageBadgeArgs {
  artifactId: string;
}

export default class DockerImageBadge extends Component<DockerImageBadgeArgs> {
  @service api!: ApiService;
  @tracked registry?: string;
  @tracked image?: string;
  @tracked tag?: string;

  constructor(owner: unknown, args: DockerImageBadgeArgs) {
    super(owner, args);

    this.checkArtifact();
  }

  clearUnicodeCharacters(string: string): string {
    let newString = string.replace(/[\u0006-\u00ff|\n]/, '');
    newString = newString.replace(/"/, '');
    return newString;
  }

  async checkArtifact(): Promise<void> {
    let ref = new Ref.Operation();
    ref.setId(this.args.artifactId);

    let artifactReq = new GetPushedArtifactRequest();
    artifactReq.setRef(ref);

    let resp = await this.api.client.getPushedArtifact(artifactReq, this.api.WithMeta());

    let textArtifact = atob(resp.getArtifact()?.getArtifact()?.getValue_asB64() || '');
    textArtifact = this.clearUnicodeCharacters(textArtifact);
    let arr = textArtifact.split(/\b[\u0006-\u0012]+\b/);

    this.image = arr[0];
    this.tag = arr[1];
  }
}
