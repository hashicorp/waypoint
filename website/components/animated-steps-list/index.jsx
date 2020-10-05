import styles from './AnimatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import FramedTerminal from 'components/framed-terminal'
import { useState } from 'react'

export default function AnimatedStepsList({ steps }) {
  const [indicatorIndex, setIndicatorIndex] = useState(0)
  return (
    <div className={styles.animatedStepsList}>
      <div className={styles.indicatorWrapper}>
        <StepsIndicator steps={steps} activeIndex={indicatorIndex} />
      </div>

      <StepsList
        className={styles.stepsList}
        steps={steps}
        onFocusedIndexChanged={(newStep) => {
          setIndicatorIndex(newStep)
        }}
      />

      <div className={styles.terminalWrapper}>
        <FramedTerminal
          frame={8}
          lines={[
            {
              frames: 1,
              code: [
                '» Building . . .',
                '» Building . . . . . .',
                '» Building . . . . . . . . . ',
                '» Building . . . . . . . . . . . .',
              ],
            },
            {
              frames: 1,
              color: 'gray',
              code: 'Creating new buildpack-based image using builder:',
              indent: 1,
            },
            {
              frames: 1,
              color: 'gray',
              code: 'heroku/buildpacks:18',
              indent: 1,
            },
            {
              frames: 2,
              color: 'navy',
              code: '✓ Creating pack client',
              indent: 1,
            },
            {
              frames: 1,
              color: 'white',
              code: '⠴ Building image',
            },
          ]}
        />
      </div>
    </div>
  )
}
