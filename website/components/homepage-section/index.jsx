import styles from './HomepageSection.module.css'
import classNames from 'classnames'

export default function HomepageSection({ title, theme, children }) {
  return (
    <section
      className={classNames(styles.homepageSection, {
        [styles.light]: theme === 'light',
        [styles.gray]: theme === 'gray',
        [styles.dark]: theme === 'dark',
      })}
    >
      <div className={styles.gridContainer}>
        {title && <h2 className="g-type-display-2">{title}</h2>}
        {children}
      </div>
    </section>
  )
}
