import Card, { CardProps } from 'components/card'
import Link from 'next/link'
import InlineSvg from '@hashicorp/react-inline-svg'
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
}: FooterProps) {
  return (
    <footer className={s.footer}>
      <div className={s.inner}>
        <div className={s.content}>
          <h2 className={s.contentTitle}>{heading}</h2>
          <p className={s.contentDescription}>{description}</p>
          {ctaLinks && ctaLinks.length
            ? ctaLinks.map((link) => {
                return (
                  <p className={s.contentLink} key={link.text}>
                    <FooterLink url={link.url} text={link.text} />
                    <RightArrowIcon />
                  </p>
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
            {/* TODO: COC link */}
            <Link href="/">
              <a>View Code of Conduct</a>
            </Link>
          </div>

          {navLinks && navLinks.length ? (
            <div className={s.bottomAnchors}>
              {navLinks.map((link) => {
                return (
                  <FooterLink key={link.text} text={link.text} url={link.url} />
                )
              })}
            </div>
          ) : null}
        </div>
      </div>
    </footer>
  )
}

function RightArrowIcon() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 20 20"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M3.334 10h13.333M11.666 5l5 5-5 5"
        stroke="#62D4DC"
        strokeWidth="1.5"
      />
    </svg>
  )
}
