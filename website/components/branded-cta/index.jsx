import styles from './branded-cta.module.css'
import Button from '@hashicorp/react-button'

export default function BrandedCta(props) {
  const { heading, content, links } = props

  return (
    <div className={styles.brandedCta}>
      <div className={`g-grid-container ${styles.contentContainer}`}>
        <h2
          data-testid="heading"
          className={`g-type-display-2 ${styles.heading}`}
        >
          {heading}
        </h2>
        <div className="content-and-links">
          <p
            data-testid="content"
            className={`g-type-body-large ${styles.content}`}
          >
            {content}
          </p>
          <div data-testid="links" className={styles.links}>
            {links.map((link, stableIdx) => {
              const buttonVariant = stableIdx === 0 ? 'primary' : 'secondary'
              const linkType = link.type || ''
              return (
                <Button
                  // eslint-disable-next-line react/no-array-index-key
                  key={stableIdx}
                  linkType={linkType}
                  theme={{
                    variant: buttonVariant,
                    brand: 'waypoint',
                    background: 'light',
                  }}
                  title={link.text}
                  url={link.url}
                />
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
