import styles from './StepsIndicator.module.css'
import classNames from 'classnames'

export default function StepsIndicator({ steps, activeIndex }) {
  return (
    <ul className={styles.stepsIndicator}>
      {steps.map((step, index) => (
        <li
          key={step.name}
          className={classNames({
            [styles.active]: index == activeIndex,
          })}
        >
          {step.name}
        </li>
      ))}
    </ul>
  )
}
