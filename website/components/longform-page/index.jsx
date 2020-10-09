import classNames from 'classnames'
import styles from './LongformPage.module.css'

export default function LongformPage({ className, title, children }) {
  return (
    <div className={classNames(styles.longformPage, className)}>
      <div className="g-container">
        <div className={styles.longformWrapper}>
          <h2>{title}</h2>
          {children}
        </div>
      </div>
    </div>
  )
}
