import styles from './HomepageHero.module.css'
import Button from '@hashicorp/react-button'

export default function HomepageHero({ title, description, subtitle, links }) {
  return (
    <div className={styles.homepageHero}>
      <div className={styles.gridContainer}>
        <div className={styles.content}>
          <h1>{title}</h1>
          <p className={styles.subtitle}>{subtitle}</p>
          <p className={styles.description}>{description}</p>
          <div className={styles.links}>
            {links.map((link, index) => {
              const brand = index === 0 ? 'waypoint' : 'neutral'
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
