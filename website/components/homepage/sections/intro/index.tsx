import Features, { FeaturesProps } from 'components/homepage/features'
import Terminal, {
  TerminalLine,
  TerminalToken,
} from 'components/homepage/terminal'
import s from './style.module.css'

interface ColumnProps {
  heading: React.ReactNode
  description: string
  features: FeaturesProps
}

export default function SectionIntro({
  columnLeft,
  columnRight,
}: {
  columnLeft: ColumnProps
  columnRight: ColumnProps
}) {
  return (
    <div className={s.intro}>
      <div className={s.column}>
        <h2 className={s.heading}>{columnLeft.heading}</h2>
        <p className={s.description}>{columnLeft.description}</p>
        <div className={s.terminal}>
          <Terminal>
            <TerminalLine>
              <TerminalToken color="teal">~</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="fushia">$</TerminalToken> waypoint up
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="fushia">»</TerminalToken>{' '}
              <TerminalToken color="green">Deploying . . .</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="fushia">✓</TerminalToken>{' '}
              <TerminalToken color="green">
                Kubernetes client connected to
                https://kubernetes.example.com:6443
              </TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="fushia">✓</TerminalToken>{' '}
              <TerminalToken color="green">Created deployment</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="fushia">✓</TerminalToken>{' '}
              <TerminalToken color="green">
                Deployment successfully rolled out! The deploy was successful! A
                Waypoint deployment URL is shown below. This can be used
                internally to check your deployment and is not meant for
                external traffic. You can manage this hostname using
                &quot;waypoint hostname&quot; Deployment URL:
                https://immensely-guided-stag--v5.waypoint.run
              </TerminalToken>
            </TerminalLine>
          </Terminal>
        </div>
        <Features items={columnLeft.features} />
      </div>
      <div className={s.column}>
        <h2 className={s.heading}>{columnRight.heading}</h2>
        <p className={s.description}>{columnRight.description}</p>
        <div className={s.terminal}>
          <Terminal
            tabs={[
              {
                label: 'Build',
                content: (
                  <>
                    <TerminalLine>
                      <TerminalToken color="green">build</TerminalToken> &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'  '}use{' '}
                      <TerminalToken color="green">
                        &quot;pack&quot;
                      </TerminalToken>{' '}
                      &#123;&#125;
                    </TerminalLine>
                    <TerminalLine>{'  '}registry &#123;</TerminalLine>
                    <TerminalLine>
                      {'    '}use{' '}
                      <TerminalToken color="green">
                        &quot;docker&quot;
                      </TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}image ={' '}
                      <TerminalToken color="green">
                        &quot;nodejs-example&quot;
                      </TerminalToken>
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}tag ={' '}
                      <TerminalToken color="green">
                        &quot;latest&quot;
                      </TerminalToken>
                    </TerminalLine>
                    <TerminalLine>{'      '}local = true</TerminalLine>
                    <TerminalLine>{'    '}&#125;</TerminalLine>
                    <TerminalLine>{'  '}&#125;</TerminalLine>
                    <TerminalLine>&#125;</TerminalLine>
                  </>
                ),
              },
              {
                label: 'Deploy',
                content: (
                  <>
                    <TerminalLine>
                      <TerminalToken color="green">deploy</TerminalToken> &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'  '}use{' '}
                      <TerminalToken color="green">
                        &quot;kubernetes&quot;
                      </TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'    '}probe_path = &quot;/&quot;
                    </TerminalLine>
                    <TerminalLine>{'  '}&#125;</TerminalLine>
                    <TerminalLine>&#125;</TerminalLine>
                  </>
                ),
              },
              {
                label: 'Release',
                content: (
                  <>
                    <TerminalLine>
                      <TerminalToken color="green">release</TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'  '}use{' '}
                      <TerminalToken color="green">
                        &quot;kubernetes&quot;
                      </TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>{'    '}load_balancer = true</TerminalLine>
                    <TerminalLine>{'    '}port = 3000</TerminalLine>
                    <TerminalLine>{'  '}&#125;</TerminalLine>
                    <TerminalLine>&#125;</TerminalLine>
                  </>
                ),
              },
            ]}
          />
        </div>
        <Features items={columnRight.features} />
      </div>
    </div>
  )
}
