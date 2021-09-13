import React from 'react'
import { Tabs, TabList, Tab, TabPanels, TabPanel } from '@reach/tabs'
import classNames from 'classnames'
import s from './style.module.css'

interface TerminalTabProps {
  label: string
  content: React.ReactNode
}
interface TerminalProps {
  tabs?: Array<TerminalTabProps>
  children?: React.ReactNode
}

function Terminal({ tabs, children }: TerminalProps): JSX.Element {
  return (
    <div className={s.terminal}>
      {tabs ? (
        <Tabs>
          <div className={s.terminalHeader}>
            <TabList className={s.terminalTabList}>
              {tabs.map((tab, idx) => (
                <React.Fragment key={tab.label}>
                  <Tab>{tab.label}</Tab>
                  {idx < tabs.length - 1 ? (
                    <span
                      aria-hidden={true}
                      className={s.terminalTabListSeperator}
                    >
                      {'>'}
                    </span>
                  ) : null}
                </React.Fragment>
              ))}
            </TabList>
          </div>
          <TabPanels className={s.terminalTabPanels}>
            {tabs.map((tab) => (
              <TabPanel key={tab.label}>{tab.content}</TabPanel>
            ))}
          </TabPanels>
        </Tabs>
      ) : (
        <div className={s.terminalContent}>
          <pre className={s.active}>
            <code>{children}</code>
          </pre>
        </div>
      )}
    </div>
  )
}

function TerminalLine({
  children,
}: {
  children: React.ReactNode
}): JSX.Element {
  return <span className={s.terminalLine}>{children}</span>
}

function TerminalToken({
  children,
  color = 'white',
}: {
  children: React.ReactNode
  color?: 'white' | 'fushia' | 'teal' | 'green'
}): JSX.Element {
  return (
    <span className={classNames(s.terminalToken, s[color])}>{children}</span>
  )
}

export default Terminal
export { TerminalLine, TerminalToken }
