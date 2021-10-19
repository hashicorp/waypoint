import React, { useEffect, useRef } from 'react'
import { useInView } from 'react-intersection-observer'
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
  const [tabIndex, setTabIndex] = React.useState(0)
  const [isHovering, setIsHovering] = React.useState(false)
  const { ref, inView } = useInView({
    threshold: 1,
  })
  const handleTabsChange = (index) => {
    setTabIndex(index)
  }
  useInterval(
    () => {
      if (!tabs) return
      if (tabIndex >= tabs.length - 1) {
        setTabIndex(0)
      } else {
        setTabIndex((prevTabIndex) => prevTabIndex + 1)
      }
    },
    isHovering || !inView ? null : 3000
  )
  return (
    <div
      ref={ref}
      className={s.terminal}
      onMouseEnter={() => setIsHovering(true)}
      onMouseLeave={() => setIsHovering(false)}
    >
      {tabs ? (
        <Tabs index={tabIndex} onChange={handleTabsChange}>
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

// https://usehooks-typescript.com/react-hook/use-interval
function useInterval(callback: () => void, delay: number | null) {
  const savedCallback = useRef(callback)

  // Remember the latest callback if it changes.
  useEffect(() => {
    savedCallback.current = callback
  }, [callback])

  // Set up the interval.
  useEffect(() => {
    // Don't schedule if no delay is specified.
    if (delay === null) {
      return
    }

    const id = setInterval(() => savedCallback.current(), delay)

    return () => clearInterval(id)
  }, [delay])
}
