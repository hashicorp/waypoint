import styles from './AnimatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import Terminal from 'components/terminal'
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
        <Terminal
          lines={[
            {
              code: '» Building . . . . . . . . . . . . .',
            },
            {
              color: 'gray',
              code: 'Creating new buildpack-based image using builder:',
              indent: 1,
            },
            {
              color: 'gray',
              code: 'heroku/buildpacks:18',
              indent: 1,
            },
            {
              color: 'navy',
              code: '✓ Creating pack client',
              indent: 1,
            },
            {
              color: 'white',
              code: '⠴ Building image',
            },
            {
              code: '',
            },
            {
              color: 'gray',
              code: 'Generated new Docker image: example-ruby:latest',
            },
            {
              color: 'gray',
              code:
                'Tagging Docker image: example-ruby:latest => gcr.io/wp-dev-277323/example-ruby:latest',
            },
            {
              color: 'white',
              code:
                'Docker image pushed: gcr.io/wp-dev-277323/example-ruby:latest',
            },
          ]}
        />
      </div>
    </div>
  )
}
