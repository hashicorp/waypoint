import { useState } from 'react'
import styles from './FeaturesList.module.css'
import Terminal from 'components/terminal'
import Button from '@hashicorp/react-button'

const features = [
  {
    id: 0,
    title: 'Preview Urls',
    description: 'Get publicly accessible preview URLs per-development',
    learnMoreLink: 'https://waypointproject.io/docs/url',
    content: (
      <Terminal
        theme="dark"
        title="waypoint-up.txt"
        lines={[
          { code: '$ waypoint deploy' },
          { code: '' },
          { code: '» Deploying...', color: 'white' },
          { code: '✓ Deployment successfully rolled out!', color: 'navy' },
          { code: '' },
          { code: '» Releasing...', color: 'white' },
          { code: '✓ Service successfully configured!', color: 'navy' },
          { code: '' },
          { code: '» Pruning old deployments...', color: 'white' },
          { code: 'Deployment: 01EJYN2P2FG4N9CTXAYGGZW9W0', color: 'white' },
          { code: '' },
          {
            code:
              'The deploy was successful! A Waypoint deployment URL is shown below.',
            color: 'white',
          },
          { code: '' },
          {
            code: 'Release URL: http://35.230.13.162:8080',
            color: 'white',
          },
          {
            code:
              'Deployment URL: https://immensely-guided-stag--v18.alpha.waypoint.run',
            color: 'white',
          },
        ]}
      />
    ),
  },
  {
    id: 1,
    title: 'Live Exec',
    description: 'Execute a command in the context of a running application',
    content: (
      <Terminal
        theme="dark"
        title="waypoint-up.txt"
        lines={[
          { code: '$ waypoint exec' },
          { code: '# Connected to 01EJY7JXPX1DZF36FMNJHQS8MH', color: 'white' },
        ]}
      />
    ),
  },
  {
    id: 2,
    title: 'Application Logs',
    description: 'View log output for running applications and deployments',
    learnMoreLink: 'https://waypointproject.io/docs/logs',
    content: (
      <Terminal
        theme="dark"
        title="waypoint-up.txt"
        lines={[
          { code: '$ waypoint logs' },
          { code: '[11] Puma starting in cluster mode...', color: 'white' },
          {
            code:
              '[11] * Version 3.11.2 (ruby 2.6.6-p146), codename: Love Song',
            color: 'white',
          },
          { code: '[11] * Min threads: 5, max threads: 5', color: 'white' },
          { code: '[11] * Environment: production', color: 'white' },
          { code: '[11] * Process workers: 2', color: 'white' },
          { code: '[11] * Preloading application', color: 'white' },
          { code: '[11] * Listening on tcp://0.0.0.0:3000', color: 'white' },
        ]}
      />
    ),
  },
  {
    id: 3,
    title: 'Web UI',
    description:
      'View projects and applications being deployed by Waypoint in a web interface',
    content: <img src="https://placehold.it/500x500" alt="placeholder" />,
  },
  {
    id: 4,
    title: 'CI/CD and Version Control Integration',
    description:
      'Integrate easily with existing CI/CD providers and version control providers like GitHub',
    content: (
      <Terminal theme="dark" title="waypoint-up.txt" lines={[{ code: '' }]} />
    ),
  },
  {
    id: 5,
    title: 'Extensible Plugin Interface',
    description:
      'Easily extend Waypoint with custom support for platforms, build processes, and release systems.',
    content: (
      <Terminal theme="dark" title="waypoint-up.txt" lines={[{ code: '' }]} />
    ),
  },
]

export default function FeaturesList() {
  const [activeFeature, setActiveFeature] = useState(0)
  return (
    <div className={styles.features}>
      <ul className={styles.options}>
        {features.map((feature) => (
          <Feature
            key={feature.id}
            title={feature.title}
            active={feature.id === activeFeature}
            onClick={() => setActiveFeature(feature.id)}
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

function Feature({ children, title, active, onClick, learnMoreLink }) {
  return (
    <li className={active ? styles.activeFeature : styles.feature}>
      <button className={styles.heading} onClick={onClick}>
        {title}
      </button>
      <div className={styles.body}>
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
