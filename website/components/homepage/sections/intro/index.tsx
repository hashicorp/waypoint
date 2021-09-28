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
              <TerminalToken color="fushia">❯</TerminalToken> waypoint up
            </TerminalLine>
            <TerminalLine>
              <TerminalToken color="green">
                Building tech-blog with Pack...
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
                    <TerminalLine>build &#123;</TerminalLine>
                    <TerminalLine>
                      {'  '}use &quot;pack&quot; &#123;&#125;
                    </TerminalLine>
                    <TerminalLine>{'  '}registry &#123;</TerminalLine>
                    <TerminalLine>
                      {'    '}use &quot;docker&quot; &#123;
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}image = &quot;nodejs-example&quot;
                    </TerminalLine>
                    <TerminalLine>
                      {'      '}tag = &quot;latest&quot;
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
                    <TerminalLine>deploy &#123;</TerminalLine>
                    <TerminalLine>
                      {'  '}use &quot;kubernetes&quot; &#123;
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
                    <TerminalLine>release &#123;</TerminalLine>
                    <TerminalLine>
                      {'  '}use &quot;kubernetes&quot; &#123;
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
