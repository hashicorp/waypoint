import { useState } from 'react'
import styles from './FeaturesList.module.css'
import Button from '@hashicorp/react-button'

export default function FeaturesList({ features }) {
  const [activeFeature, setActiveFeature] = useState(0)
  return (
    <div className={styles.features}>
      <ul className={styles.options}>
        {features.map((feature, stableIdx) => (
          <Feature
            id={stableIdx}
            key={stableIdx}
            title={feature.title}
            active={stableIdx === activeFeature}
            onClick={setActiveFeature}
            learnMoreLink={feature.learnMoreLink}
          >
            {feature.description}
          </Feature>
        ))}
      </ul>
      <div className={styles.terminalWrapper}>
        {features[activeFeature].content}
      </div>
    </div>
  )
}

function Feature({ children, title, active, onClick, learnMoreLink, id }) {
  return (
    <li className={active ? styles.activeFeature : styles.feature}>
      <button
        className={styles.heading}
        onClick={() => onClick(id)}
        aria-expanded={active}
        aria-controls={`feature-${id}`}
      >
        {title}
      </button>
      <div className={styles.body} id={`feature-${id}`} aria-hidden={!active}>
        <p>{children}</p>
        {learnMoreLink && (
          <Button
            url={learnMoreLink}
            title="Learn more"
            linkType="inbound"
            theme={{
              variant: 'tertiary-neutral',
              brand: 'terraform',
            }}
          />
        )}
      </div>
    </li>
  )
}
