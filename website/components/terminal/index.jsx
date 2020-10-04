import classNames from 'classnames'
import styles from './Terminal.module.css'

export default function Terminal({ lines, className }) {
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
        <div className={styles.overflowWrapper}>
          <div className={styles.codeWrapper}>
            {lines &&
              lines.map((line) => (
                <>
                  <code
                    className={classNames({
                      [styles.navy]: line.color === 'navy',
                      [styles.gray]: line.color === 'gray',
                      [styles.white]: line.color === 'white',
                    })}
                  >
                    {line.indent &&
                      new Array(line.indent * 2)
                        .fill({})
                        .map(() => <>&nbsp;</>)}
                    {line.code}
                  </code>
                  <br />
                </>
              ))}
          </div>
        </div>
      </div>
    </div>
  )
}
