import ReleaseDetail from './release';
import { GetReleaseRequest, Release, Ref } from 'waypoint-pb';

interface ReleaseModelParams {
  release_id: string;
}

export default class ReleaseIdDetail extends ReleaseDetail {
  renderTemplate() {
    this.render('workspace/projects/project/app/release', {
      into: 'workspace/projects/project',
    });
  }

  async model(params: ReleaseModelParams): Promise<Release.AsObject> {
    let ref = new Ref.Operation();
    ref.setId(params.release_id);
    let req = new GetReleaseRequest();
    req.setRef(ref);

    let release: Release = await this.api.client.getRelease(req, this.api.WithMeta());

    return release.toObject();
  }
}
