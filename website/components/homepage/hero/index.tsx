import { useEffect, useState } from 'react'
import classNames from 'classnames'
import InlineSvg from '@hashicorp/react-inline-svg'
import Mask from './Mask'
import Logos from './Logos'
import s from './style.module.css'

export default function Hero() {
  const [loaded, setLoaded] = useState(false)
  useEffect(() => {
    setTimeout(() => {
      setLoaded(true)
    }, 250)
  }, [])
  return (
    <header className={s.hero}>
      <div className={s.heroInner}>
        <h1 className={s.heroHeading}>
          Easy application deployment for <em>Kubernetes</em> and <em>ECS</em>
        </h1>
        <p className={s.heroDescription}>
          Waypoint is an application deployment tool for Kubernetes, ECS, and
          many other platforms. It allows developers to deploy, manage, and
          observe their applications through a consistent abstraction of the
          underlying infrastructure.
        </p>
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
