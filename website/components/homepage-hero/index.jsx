import s from './style.module.css'
import Button from '@hashicorp/react-button'

export default function HomepageHero({ title, description, links }) {
  return (
    <div className={s.root}>
      <div className="g-grid-container">
        <span className={s.eyebrow}>
          Welcome to the <span className={s.red}>INTERNAL BETA</span> for
          HashiCorp Waypoint! This is a confidential internal only beta. No
          details should be shared externally.
        </span>
        <h1 className="g-type-display-1">{title}</h1>
        <div className={s.contentAndLinks}>
          <p className="g-type-body-large">{description}</p>
          <div className={s.links}>
            {links.map((link, index) => {
              const brand = index === 0 ? 'hashicorp' : 'neutral'
              const variant = index === 0 ? 'primary' : 'secondary'
              return (
                <Button
                  key={link.text}
                  title={link.text}
                  linkType={link.type}
                  url={link.url}
                  theme={{ variant, brand }}
                />
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
