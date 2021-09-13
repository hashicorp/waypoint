import MediaObject from '../media-object'
import Terminal, { TerminalLine, TerminalToken } from '../terminal'
import s from './style.module.css'

export default function Intro() {
  return (
    <div className={s.intro}>
      <div className={s.column}>
        <h2 className={s.heading}>
          Easy for <em>developers</em>
        </h2>
        <p className={s.description}>
          Waypoint enables developers to deploy, manage, and observe their
          applications to Kubernetes, ECS, and many other platforms through a
          consistent abstraction.
        </p>
        <div className={s.terminal}>
          <Terminal>
            <TerminalLine>
              <TerminalToken color="teal">~</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="fushia">❯</TerminalToken> waypoint up
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="green">
                Building tech-blog with Pack...
              </TerminalToken>
            </TerminalLine>
          </Terminal>
        </div>
        <MediaObject
          icon={require('../icons/eye.svg?include')}
          heading="Application-centric Abstraction"
          description="Specify the deployment needs with a simple and consistent abstraction without the underlying complexity."
        />
        <MediaObject
          icon={require('../icons/eye.svg?include')}
          heading="End-to-End Deployment Workflow"
          description="Build a complete end-to-end workflow with distinct build, deploy, release steps."
        />
      </div>
      <div className={s.column}>
        <h2 className={s.heading}>
          Powerful for <em>platform operators</em>
        </h2>
        <p className={s.description}>
          Waypoint enables operators to create PaaS workflows of Kubernetes,
          ECS, serverless applications
        </p>
        <div className={s.terminal}>
          <Terminal
            tabs={[
              {
                label: 'Build',
                content: <TerminalLine>Build</TerminalLine>,
              },
              {
                label: 'Deploy',
                content: <TerminalLine>Deploy</TerminalLine>,
              },
              {
                label: 'Release',
                content: <TerminalLine>Release</TerminalLine>,
              },
            ]}
          />
        </div>
        <MediaObject
          icon={require('../icons/eye.svg?include')}
          heading="Build-Deploy-Release Extensibility"
          description="Enable a pluggable framework, integrated with CI/CD pipelines, monitoring tools, and any other ecosystem tools. "
        />
        <MediaObject
          icon={require('../icons/eye.svg?include')}
          heading="PaaS Experience for Developers"
          description="Provide a consistent abstraction and unified workflow to scale across multiple platforms and clouds"
        />
      </div>
    </div>
  )
}
