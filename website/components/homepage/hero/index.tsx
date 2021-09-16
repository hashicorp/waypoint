import { useEffect, useState } from 'react'
import classNames from 'classnames'
import Mask from './Mask'
import Logos from './Logos'
import s from './style.module.css'

interface HeroProps {
  heading: React.ReactNode
  description: string
}

export default function Hero({ heading, description }: HeroProps): JSX.Element {
  const [loaded, setLoaded] = useState(false)
  useEffect(() => {
    setTimeout(() => {
      setLoaded(true)
    }, 250)
  }, [])
  return (
    <header className={s.hero}>
      <div className={s.heroInner}>
        <h1 className={s.heroHeading}>{heading}</h1>
        <p className={s.heroDescription}>{description}</p>
      </div>
      <div
        className={classNames(s.heroGraphic, {
          [s.visible]: loaded,
        })}
      >
        <div className={s.heroMask}>
          <Mask />
        </div>
        <div className={s.heroLogos}>
          <Logos />
        </div>
      </div>
    </header>
  )
}
