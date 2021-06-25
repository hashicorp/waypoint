import VERSION, { packageManagers } from 'data/version.js'
import HashiHead from '@hashicorp/react-head'
import { productName, productSlug } from 'data/metadata'
import ProductDownloader from '@hashicorp/react-product-downloader'
import styles from './style.module.css'

export default function DownloadsPage({ releases }) {
  return (
    <span className={styles.downloads}>
      <HashiHead title={`Downloads | ${productName} by HashiCorp`} />
      <ProductDownloader
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
    </span>
  )
}

export async function getStaticProps() {
  return fetch(`https://releases.hashicorp.com/waypoint/index.json`, {
    headers: {
      'Cache-Control': 'no-cache',
    },
  })
    .then((res) => res.json())
    .then((result) => {
      return {
        props: {
          releases: result,
        },
      }
    })
    .catch(() => {
      throw new Error(
        `--------------------------------------------------------
        Unable to resolve version ${VERSION} on releases.hashicorp.com from link
        <https://releases.hashicorp.com/${productSlug}/${VERSION}/index.json>. Usually this
        means that the specified version has not yet been released. The downloads page
        version can only be updated after the new version has been released, to ensure
        that it works for all users.
        ----------------------------------------------------------`
      )
    })
}
