import styles from './StepsList.module.css'
import Step from './step'

export default function StepsList({ steps }) {
  return (
    <ul className={styles.stepsList}>
      {steps.map((step) => (
        <Step
          key={step.name}
          name={step.name}
          description={step.description}
          logos={step.logos}
        />
      ))}
    </ul>
  )
}
