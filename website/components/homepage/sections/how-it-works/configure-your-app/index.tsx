import NumberedBlock from 'components/homepage/numbered-block'
import MediaObject from 'components/homepage/media-object'
import s from './style.module.css'

export default function ConfigureYourApp() {
  return (
    <div className={s.root}>
      <div className={s.content}>
        <NumberedBlock index="1" heading="Configure your app for Waypoint">
          <MediaObject
            icon={require('../../icons/edit-pencil.svg?include')}
            heading="Writing waypoint.hcl files"
            description="Your waypoint.hcl file defines how Waypoint builds, deploys, and releases a project."
          />
          <MediaObject
            icon={require('../../icons/layout.svg?include')}
            heading="Sample Waypoint files"
            description="View sample waypoint.hcl files to see how straight-forward it is to configure your deployments"
          />
        </NumberedBlock>
      </div>
      <div className={s.media}>
        <img
          src={require('../img/configure-your-app-for-waypoint.png')}
          alt=""
        />
      </div>
    </div>
  )
}
