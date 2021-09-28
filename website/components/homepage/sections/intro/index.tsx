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
              <TerminalToken color="green">$ waypoint up</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="white">» Deploying . . .</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="white">✓</TerminalToken>{' '}
              <TerminalToken color="green">
                Kubernetes client connected to
                https://kubernetes.example.com:6443
              </TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="white">✓</TerminalToken>{' '}
              <TerminalToken color="green">Created deployment</TerminalToken>
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="white">✓</TerminalToken>{' '}
              <TerminalToken color="green">
                Deployment successfully rolled out!
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
                      <TerminalToken color="teal">
                        &quot;pack&quot;
                      </TerminalToken>{' '}
                      &#123;&#125;
                    </TerminalLine>
                    <TerminalLine>{'  '}registry &#123;</TerminalLine>
                    <TerminalLine>
                      {'    '}use{' '}
                      <TerminalToken color="teal">
                        &quot;docker&quot;
                      </TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}
                      <TerminalToken color="green">image</TerminalToken> ={' '}
                      <TerminalToken color="teal">
                        &quot;nodejs-example&quot;
                      </TerminalToken>
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}
                      <TerminalToken color="green">tag</TerminalToken> ={' '}
                      <TerminalToken color="teal">
                        &quot;latest&quot;
                      </TerminalToken>
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}
                      <TerminalToken color="green">local</TerminalToken> = true
                    </TerminalLine>
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
                      <TerminalToken color="teal">
                        &quot;kubernetes&quot;
                      </TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'    '}
                      <TerminalToken color="green">
                        probe_path
                      </TerminalToken> ={' '}
                      <TerminalToken color="teal">&quot;/&quot;</TerminalToken>
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
                      <TerminalToken color="teal">
                        &quot;kubernetes&quot;
                      </TerminalToken>{' '}
                      &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'    '}
                      <TerminalToken color="green">
                        load_balancer
                      </TerminalToken>{' '}
                      = true
                    </TerminalLine>
                    <TerminalLine>
                      {'    '}
                      <TerminalToken color="green">port</TerminalToken> = 3000
                    </TerminalLine>
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
