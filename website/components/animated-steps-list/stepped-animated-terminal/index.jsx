import { useState } from 'react'
import AnimatedTerminal from 'components/animated-terminal'

export default function SteppedAnimatedTerminal({ activeIndex, steps }) {
  const [currentIndex, setCurrentIndex] = useState(activeIndex)
  if (activeIndex != currentIndex) {
    setCurrentIndex(activeIndex)
  }

  return (
    <AnimatedTerminal
      frameLength={steps[activeIndex].frameLength}
      loop={steps[activeIndex].loop}
      lines={steps[activeIndex].lines}
    />
  )
}
