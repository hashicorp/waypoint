import s from './style.module.css'
import Link from 'next/link'

export default function Footer({ openConsentManager }) {
  return (
    <footer className={s.root}>
      <div className="g-container">
        <div className={s.left}>
          <Link href="/docs">
            <a>Docs</a>
          </Link>
          <Link href="/cli">
            <a>CLI</a>
          </Link>
          <a href="https://learn.hashicorp.com/waypoint">Learn</a>
          <a href="https://hashicorp.com/privacy">Privacy</a>
          <Link href="/security">
            <a>Security</a>
          </Link>
          <a href="/files/press-kit.zip">Press Kit</a>
          <a onClick={openConsentManager}>Consent Manager</a>
        </div>
      </div>
    </footer>
  )
}
