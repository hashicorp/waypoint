import styles from './Step.module.css'
import { useState } from 'react'
import { useInView } from 'react-intersection-observer'

export default function Step({
  name,
  description,
  logos,
  logosAlt,
  onInViewStatusChanged,
}) {
  const [ref, inView] = useInView({ threshold: 0.4 })
  const [inViewStatus, setInViewStatus] = useState(false)
  if (inView != inViewStatus) {
    setInViewStatus(inView)
    onInViewStatusChanged(inView)
  }
  return (
    <li className={styles.step} ref={ref}>
      <h4>{name}</h4>
      <div className={styles.description}>{description}</div>
      <img src={logos} alt={logosAlt} />
    </li>
  )
}
