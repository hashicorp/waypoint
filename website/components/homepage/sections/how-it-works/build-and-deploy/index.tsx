import { useState } from 'react'
import { useInView } from 'react-intersection-observer'
import Typical from 'react-typical'
import InlineSvg from '@hashicorp/react-inline-svg'
import NumberedBlock from 'components/homepage/numbered-block'
import MediaObject from 'components/homepage/media-object'
import classNames from 'classnames'
import s from './style.module.css'

export default function BuildAndDeploy() {
  const [typeFinished, setTypeFinished] = useState(false)
  const { ref, inView } = useInView({
    triggerOnce: true,
  })
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
        <div ref={ref} className={classNames({ [s.active]: typeFinished })}>
          <div className={s.terminalContainer}>
            <pre className={s.terminal}>
              <code>
                <span className={s.terminalLine}>
                  <span className={s.terminalTilde}>~</span>
                </span>
                <span className={s.terminalLine}>
                  <span className={s.terminalToken}>❯</span>{' '}
                  {inView && (
                    <Typical
                      steps={[
                        100,
                        'waypoint up',
                        500,
                        () => setTypeFinished(true),
                      ]}
                      wrapper="span"
                    />
                  )}
                </span>
                <span className={s.terminalLine}>
                  <span className={s.terminalNote}>
                    Building tech-blog with Pack...
                  </span>
                </span>
              </code>
            </pre>
            <p className={s.note}>Instantly deploy from the command line...</p>
          </div>

          <div className={s.deploymentContainer}>
            <InlineSvg
              className={s.arrow}
              src={require('./arrow.svg?include')}
            />
            <div className={s.deploymentWrapper}>
              <div className={s.deployment}>
                <div className={s.deploymentHeading}>
                  <span>Deploy from</span>
                  <div className={s.deploymentHeadingBranch}>
                    <InlineSvg src={require('./github.svg?include')} />{' '}
                    hashicorp/tech-blog
                  </div>
                </div>
                <span className={s.deploymentCommit}>
                  Last commit 3 seconds ago by{' '}
                  <img
                    src={require('./avatar.jpg')}
                    width={16}
                    height={16}
                    alt=""
                  />
                  <strong>@almonk</strong>
                </span>
              </div>
              <p className={s.note}>
                or connect to GitHub for automatic deploys
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
