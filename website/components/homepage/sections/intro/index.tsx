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
        <Features items={columnRight.features} />
      </div>
    </div>
  )
}
