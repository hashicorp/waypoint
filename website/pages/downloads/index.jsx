import VERSION, { packageManagers } from 'data/version.js'
import { productName, productSlug } from 'data/metadata'
import ProductDownloadsPage from '@hashicorp/react-product-downloads-page'
import { generateStaticProps } from '@hashicorp/react-product-downloads-page/server'
import styles from './style.module.css'

export default function DownloadsPage({ releases }) {
  return (
    <ProductDownloadsPage
      releases={releases}
      packageManagers={packageManagers}
      productName={productName}
      productId={productSlug}
      latestVersion={VERSION}
      getStartedDescription="Follow step-by-step tutorials on AWS, Azure, GCP, and localhost."
      getStartedLinks={[
        {
          label: 'Deploy to Docker',
          href:
            'https://learn.hashicorp.com/collections/waypoint/get-started-docker',
        },
        {
          label: 'Deploy to Kubernetes',
          href:
            'https://learn.hashicorp.com/collections/waypoint/get-started-kubernetes',
        },
        {
          label: 'Deploy to AWS',
          href: 'https://learn.hashicorp.com/collections/waypoint/deploy-aws',
        },
        {
          label: 'View all Waypoint tutorials',
          href: 'https://learn.hashicorp.com/waypoint',
        },
      ]}
      logo={
        <img
          className={styles.logo}
          alt="Waypoint"
          src={require('./img/waypoint-logo.svg')}
        />
      }
      product="waypoint"
      tutorialLink={{
        href: 'https://learn.hashicorp.com/waypoint',
        label: 'View Tutorials at HashiCorp Learn',
      }}
    />
  )
}

export const getStaticProps = () =>
  generateStaticProps({
    product: 'waypoint',
    latestVersion: VERSION,
  })
