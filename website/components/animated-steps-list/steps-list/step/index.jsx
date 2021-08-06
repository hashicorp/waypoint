import styles from './Step.module.css'
import { useState, useEffect } from 'react'
import { useInView } from 'react-intersection-observer'
import LogoList from 'components/logo-list'

export default function Step({
  name,
  description,
  logos,
  onInViewStatusChanged,
}) {
  const [ref, inView] = useInView({ threshold: 0.4 })
  const [inViewStatus, setInViewStatus] = useState(false)

  useEffect(() => {
    if (inView !== inViewStatus) {
      setInViewStatus(inView)
      onInViewStatusChanged(inView)
    }
  }, [inView, inViewStatus])

  return (
    <li className={styles.step} ref={ref}>
      <h4 className="g-type-display-4">{name}</h4>
      <div className={styles.description}>{description}</div>
      <LogoList className={styles.logoList} logos={logos} />
    </li>
  )
}
