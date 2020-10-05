import classNames from 'classnames'
import styles from './Terminal.module.css'

/**
 * A Terminal is a simple component representing the presentation
 * of a terminal in a static state.
 *
 * Animated versions of the terminal are available in higher order components,
 * but they work by manipulating the passed props down to this terminal to
 * represent the active state in a given frame.
 *
 * Example Usage:
 *
 *  <Terminal
 *    lines={[
 *      {
 *        code: '» Building . . . . . . . . . . . . .',
 *      },
 *      {
 *        color: 'gray',
 *        code: 'Creating new buildpack-based image using builder:',
 *        indent: 1,
 *      },
 *      {
 *        color: 'gray',
 *        code: 'heroku/buildpacks:18',
 *        indent: 1,
 *      },
 *      {
 *        color: 'navy',
 *        code: '✓ Creating pack client',
 *        indent: 1,
 *      },
 *      {
 *        color: 'white',
 *        code: '⠴ Building image',
 *      },
 *    ]}
 *  />
 */
export default function Terminal({ lines }) {
  return (
    <div className={styles.terminal}>
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
