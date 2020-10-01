import styles from './IndicatedStepsList.module.css'
import StepsList from './steps-list'
import classNames from 'classnames'
import { useState } from 'react'

export default function IndicatedStepsList({ steps }) {
  const [indicatorIndex, setIndicatorIndex] = useState(0)

  return (
    <div className={styles.indicatedStepsList}>
      <div className={styles.indicatorWrapper}>
        <ul className={styles.indicator}>
          {steps.map((step, index) => (
            <li
              key={step.name}
              className={classNames({
                [styles.active]: index == indicatorIndex,
              })}
            >
              {step.name}
            </li>
          ))}
        </ul>
      </div>

      <StepsList
        steps={steps}
        onFocusedIndexChanged={(newStep) => {
          setIndicatorIndex(newStep)
        }}
      />
    </div>
  )
}
