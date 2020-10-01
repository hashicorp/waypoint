import React, { useState, useEffect } from 'react'
import Terminal from 'components/terminal'

export default function SimpleAnimatedTerminal({ lines }) {
  const [seconds, setSeconds] = useState(0)
  useEffect(() => {
    let interval = setInterval(() => {
      setSeconds((seconds) => seconds + 1)
    }, 500)
    return () => clearInterval(interval)
  }, [seconds])
  return <Terminal lines={lines.slice(0, (seconds % lines.length) + 1)} />
}
