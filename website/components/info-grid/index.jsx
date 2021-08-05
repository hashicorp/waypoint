import styles from './InfoGrid.module.css'

export default function InfoGrid({ items }) {
  return (
    <div className={styles.infoGrid}>
      <ul>
        {items.map((item) => (
          <li key={item.title}>
            <img src={item.icon} />
            <h4 className="g-type-display-4">{item.title}</h4>
            <p>{item.description}</p>
          </li>
        ))}
      </ul>
    </div>
  )
}
