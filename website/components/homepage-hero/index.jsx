import Button from '@hashicorp/react-button'

export default function HomepageHero({ title, description, links }) {
  return (
    <div className="g-homepage-hero">
      <div className="g-grid-container">
        <span className="eyebrow">
          Welcome to the <span className="red">INTERNAL BETA</span> for
          HashiCorp Waypoint! This is a confidential internal only beta. No
          details should be shared externally.
        </span>
        <h1 data-testid="heading" className="g-type-display-1">
          {title}
        </h1>
        <div className="content-and-links">
          <p data-testid="content" className="g-type-body-large">
            {description}
          </p>
          <div data-testid="links" className="links">
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
