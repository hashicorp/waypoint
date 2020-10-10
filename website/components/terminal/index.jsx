import { Fragment } from 'react'
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
export default function Terminal({ lines, title, noScroll }) {
  return (
    <div className={styles.terminal}>
      <div className={styles.titleBar}>
        <ul className={styles.windowControls}>
          <li></li>
          <li></li>
          <li></li>
        </ul>
        {title && <div className={styles.title}>{title}</div>}
      </div>
      <div className={styles.content}>
        <div
          className={
            noScroll ? styles.noScrollOverflowWrapper : styles.overflowWrapper
          }
        >
          <div className={styles.codeWrapper}>
            {lines &&
              lines.map((line, index) => (
                <Fragment key={index}>
                  <pre
                    className={classNames({
                      [styles.short]: line.short,
                      [styles.navy]: line.color === 'navy',
                      [styles.gray]: line.color === 'gray',
                      [styles.white]: line.color === 'white',
                    })}
                  >
                    {line.indent &&
                      new Array(line.indent * 2)
                        .fill({})
                        .map((_, index) => (
                          <Fragment key={index}>&nbsp;</Fragment>
                        ))}
                    {line.code}
                  </pre>
                </Fragment>
              ))}
          </div>
        </div>
      </div>
    </div>
  )
}
