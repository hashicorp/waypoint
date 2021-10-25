import Button from '@hashicorp/react-button'
import Card, { CardProps } from 'components/card'
import Link from 'next/link'
import InlineSvg from '@hashicorp/react-inline-svg'
import useWaypointServiceStatus from 'lib/hooks/useWaypointServiceStatus'
import s from './style.module.css'

interface LinkProps {
  text: string
  url: string
}

interface FooterProps {
  heading: string
  description: string
  cards?: [CardProps, CardProps] // Require two cards
  navLinks?: Array<LinkProps>
  ctaLinks?: Array<LinkProps>
  openConsentManager: () => void
}

function FooterLink({ text, url }) {
  const isInternalLink = url && (url.startsWith('/') || url.startsWith('#'))
  if (isInternalLink) {
    return (
      <Link href={url}>
        <a>{text}</a>
      </Link>
    )
  }

  return <a href={url}>{text}</a>
}

export default function Footer({
  heading,
  description,
  cards,
  ctaLinks,
  navLinks,
  openConsentManager,
}: FooterProps) {
  const waypointServiceOK = useWaypointServiceStatus()
  return (
    <footer className={s.footer}>
      <div className={s.inner}>
        <div className={s.content}>
          <h2 className={s.contentTitle}>{heading}</h2>
          <p className={s.contentDescription}>{description}</p>
          {ctaLinks && ctaLinks.length
            ? ctaLinks.map((link) => {
                return (
                  <Button
                    key={link.url}
                    className={s.contentLink}
                    title={link.text}
                    url={link.url}
                    linkType="inbound"
                    theme={{
                      variant: 'tertiary',
                      brand: 'neutral',
                      background: 'dark',
                    }}
                  />
                )
              })
            : null}
        </div>

        {cards && cards.length ? (
          <div className={s.cards}>
            {cards.map((card) => {
              return (
                <Card
                  key={card.title}
                  link={card.link}
                  img={card.img}
                  eyebrow={card.eyebrow}
                  title={card.title}
                  description={card.description}
                />
              )
            })}
          </div>
        ) : null}

        <div className={s.bottom}>
          <div className={s.bottomMeta}>
            <InlineSvg src={require('./hashicorp-logo.svg?include')} />
            <p>Waypoint is maintained by HashiCorp, Inc.</p>
            <a href="https://github.com/hashicorp/waypoint/blob/main/.github/CODE_OF_CONDUCT.md">
              View Code of Conduct
            </a>
          </div>

          {navLinks && navLinks.length ? (
            <div className={s.bottomAnchors}>
              {navLinks.map((link) => {
                return (
                  <FooterLink key={link.text} text={link.text} url={link.url} />
                )
              })}
              <button onClick={openConsentManager}>Consent Manager</button>
              {waypointServiceOK && (
                <div className={s.status}>
                  <Link href="https://status.hashicorp.com">
                    <a className={s.normal}>All systems normal</a>
                  </Link>
                </div>
              )}
            </div>
          ) : null}
        </div>
        {!waypointServiceOK && (
          <a className={s.statusBanner} href="https://status.hashicorp.com">
            <span>
              The Waypoint URL service is currently experiencing an issue. View
              information on <b>status.hashicorp.com</b>
            </span>
          </a>
        )}
      </div>
    </footer>
  )
}
