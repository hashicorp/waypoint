import { useEffect, useState } from 'react'

const useScrollPosition = () => {
  const [scrollPosition, setScrollPosition] = useState(0)
  const [ticking, setTicking] = useState(false)

  useEffect(() => {
    window.addEventListener('scroll', handleScroll)

    return () => {
      window.removeEventListener('scroll', handleScroll)
    }
  }, [ticking])

  function handleScroll() {
    if (!ticking) {
      window.requestAnimationFrame(function () {
        setScrollPosition(window.scrollY)
        setTicking(false)
      })
      setTicking(true)
    }
  }
  return scrollPosition
}
export default useScrollPosition
