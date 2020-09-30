import styles from './StepsList.module.css'

export default function StepsList({ steps }) {
  return (
    <ul className={styles.stepsList}>
      {steps.map((step) => (
        <li key={step.name}>
          <h4>{step.name}</h4>
          <div className={styles.description}>{step.description}</div>
          <img src={step.logos} />
        </li>
      ))}
    </ul>
  )
}
