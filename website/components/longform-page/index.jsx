import classNames from 'classnames'
import styles from './LongformPage.module.css'

export default function LongformPage({ className, title, alert, children }) {
  return (
    <div className={classNames(styles.longformPage, className)}>
      <div className="g-container">
        <div className={styles.longformWrapper}>
          {alert && <div className={styles.alert}>{alert}</div>}
          <h2 className="g-type-display-2">{title}</h2>
          {children}
        </div>
      </div>
    </div>
  )
}
