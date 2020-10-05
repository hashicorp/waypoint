import AnimatedTerminal from 'components/animated-terminal'

export default function SteppedAnimatedTerminal({ activeIndex, steps }) {
  return (
    <AnimatedTerminal
      frameLength={steps[activeIndex].frameLength}
      loop={steps[activeIndex].loop}
      lines={steps[activeIndex].lines}
    />
  )
}
