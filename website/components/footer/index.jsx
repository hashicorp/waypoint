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
          <Link href="/commands">
            <a>CLI</a>
          </Link>
          <a href="https://learn.hashicorp.com/waypoint">Learn</a>
          <a href="https://hashicorp.com/privacy">Privacy</a>
          <Link href="/security">
            <a>Security</a>
          </Link>
          <Link href="/terms">
            <a>Terms of Use</a>
          </Link>
          <Link href="/copyright-policy">
            <a>Copyright Policy</a>
          </Link>
          <a href="/files/press-kit.zip">Press Kit</a>
          <a onClick={openConsentManager}>Consent Manager</a>
          <Link href="https://status.hashicorp.com">
            <a>URL Service Status</a>
          </Link>
        </div>
      </div>
    </footer>
  )
}
