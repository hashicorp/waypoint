import ReleaseDetail from './release';

export default class ReleaseIdDetail extends ReleaseDetail {
  renderTemplate(): void {
    this.render('workspace/projects/project/app/release', {
      into: 'workspace/projects/project',
    });
  }
}
