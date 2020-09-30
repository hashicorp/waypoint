import styles from './Step.module.css'

export default function Step({ name, description, logos }) {
  return (
    <li className={styles.step}>
      <h4>{name}</h4>
      <div className={styles.description}>{description}</div>
      <img src={logos} />
    </li>
  )
}
