import { useState } from 'react'
import { useInView } from 'react-intersection-observer'
import Typical from 'react-typical'
import InlineSvg from '@hashicorp/react-inline-svg'
import NumberedBlock from 'components/homepage/numbered-block'
import Features, { FeaturesProps } from 'components/homepage/features'
import usePrefersReducedMotion from 'lib/hooks/usePrefersReducedMotion'
import classNames from 'classnames'
import s from './style.module.css'

interface BuildAndDeployProps {
  heading: string
  features: FeaturesProps
}

export default function BuildAndDeploy({
  heading,
  features,
}: BuildAndDeployProps): JSX.Element {
  const prefersReducedMotion = usePrefersReducedMotion()
  const [typeFinished, setTypeFinished] = useState(false)
  const { ref, inView } = useInView({
    triggerOnce: true,
  })
  return (
    <div className={s.root}>
      <div className={s.content}>
        <NumberedBlock index="2" heading={heading}>
          <Features items={features} />
          <InlineSvg className={s.logos} src={require('./logos.svg?include')} />
        </NumberedBlock>
      </div>

      <div className={s.media}>
        <div
          ref={ref}
          className={classNames({
            [s.active]: typeFinished || prefersReducedMotion,
          })}
        >
          <div className={s.terminalContainer}>
            <pre className={s.terminal}>
              <code>
                <span className={s.terminalLine}>
                  <span className={s.terminalTilde}>~</span>
                </span>
                <span className={s.terminalLine}>
                  <span className={s.terminalToken}>❯</span>{' '}
                  {prefersReducedMotion
                    ? 'waypoint up'
                    : inView && (
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
                    <InlineSvg src={require('./github-repo.svg?include')} />{' '}
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
