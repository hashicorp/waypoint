import { useState } from 'react'
import styles from './FeaturesList.module.css'
import Button from '@hashicorp/react-button'
import Carousel from 'nuka-carousel'

export default function Features({ features }) {
  return (
    <>
      <FeaturesCarousel features={features} />
      <FeaturesList features={features} />
    </>
  )
}

function FeaturesList({ features }) {
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

function FeaturesCarousel({ features }) {
  return (
    <div className={styles.featuresCarousel}>
      <Carousel
        renderCenterRightControls={() => null}
        renderCenterLeftControls={() => null}
        slideWidth={1.0}
        defaultControlsConfig={{
          pagingDotsContainerClassName: styles.pagingDots,
        }}
        cellSpacing={40}
        getControlsContainerStyles={(key) => {
          switch (key) {
            case 'BottomCenter':
              return {
                top: 0,
              }
          }
        }}
      >
        {features.map((feature, stableIdx) => (
          <div key={stableIdx}>
            <Feature Element="div" id={stableIdx} title={feature.title} active>
              {feature.description}
            </Feature>
            <div className={styles.terminalWrapper}>
              {features[stableIdx].content}
            </div>
          </div>
        ))}
      </Carousel>
    </div>
  )
}

function Feature({
  children,
  title,
  active,
  onClick,
  learnMoreLink,
  id,
  Element = 'li',
}) {
  return (
    <Element className={active ? styles.activeFeature : styles.feature}>
      {onClick ? (
        <button
          className={styles.heading}
          onClick={() => onClick(id)}
          aria-expanded={active}
          aria-controls={`feature-${id}`}
        >
          {title}
        </button>
      ) : (
        <span className={styles.heading}>{title}</span>
      )}
      <div className={styles.body} id={`feature-${id}`} aria-hidden={!active}>
        <p>{children}</p>
        {learnMoreLink && (
          <Button
            url={learnMoreLink}
            className={styles.learnMoreLink}
            title="Learn more"
            linkType="inbound"
            theme={{
              variant: 'tertiary-neutral',
              brand: 'terraform',
            }}
          />
        )}
      </div>
    </Element>
  )
}
