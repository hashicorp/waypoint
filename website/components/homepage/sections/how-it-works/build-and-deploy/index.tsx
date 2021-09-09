import NumberedBlock from 'components/homepage/numbered-block'
import MediaObject from 'components/homepage/media-object'
import s from './style.module.css'

export default function BuildAndDeploy() {
  return (
    <div className={s.root}>
      <div className={s.content}>
        <NumberedBlock index="2" heading="Build and deploy">
          <MediaObject
            icon={require('../../icons/file-plus.svg?include')}
            heading="One simple command"
            description="Perform the build, deploy, and release steps for the app all from one simple command. Or instrument your Waypoint deployments through Remote or Git operations"
          />
        </NumberedBlock>
      </div>
      <div className={s.media}>
        <img src={require('../img/build-and-deploy.png')} alt="" />
      </div>
    </div>
  )
}
