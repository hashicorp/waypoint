import classNames from 'classnames'
import styles from './Terminal.module.css'

export default function Terminal({ code, className }) {
  return (
    <div className={classNames(styles.terminal, className)}>
      <div className={styles.titleBar}>
        <ul className={styles.windowControls}>
          <li></li>
          <li></li>
          <li></li>
        </ul>
      </div>
      <div className={styles.content}>
        <code>{code}</code>
      </div>
    </div>
  )
}
