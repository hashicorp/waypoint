import { useMemo, useState, useEffect } from 'react'
import VERSION from 'data/version.js'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'
import { productName, productSlug } from 'data/metadata'
import { packageManagersByOs, getStartedLinks } from 'data/downloads'
import ReleaseInformation from 'components/downloader/release-information'
import {
  sortPlatforms,
  detectOs,
  sortAndFilterReleases,
} from 'components/downloader/utils/downloader'
import DownloadCards from 'components/downloader/cards'
import styles from './style.module.css'

export default function DownloadsPage({
  currentVersionReleaseData,
  allVersionsReleaseData,
}) {
  const sortedDownloads = useMemo(
    () => sortPlatforms(currentVersionReleaseData),
    [currentVersionReleaseData]
  )
  const osKeys = Object.keys(sortedDownloads)
  const [osIndex, setSelectedOsIndex] = useState()

  const tabData = Object.keys(sortedDownloads).map((osKey) => ({
    os: osKey,
    packageManagers: packageManagersByOs[osKey] || null,
  }))

  useEffect(() => {
    // if we're on the client side, detect the default platform only on initial render
    const index = osKeys.indexOf(detectOs(window.navigator.platform))
    setSelectedOsIndex(index)
  }, [])

  return (
    <div className={styles.root}>
      <h1>Download {productName}</h1>
      <HashiHead is={Head} title={`Downloads | ${productName} by HashiCorp`} />
      <DownloadCards
        brand="blue"
        defaultTabIdx={osIndex}
        tabData={tabData}
        downloads={sortedDownloads}
        version={VERSION}
        logo={
          <img
            className={styles.logo}
            alt="Waypoint"
            src={require('./img/waypoint-logo.svg')}
          />
        }
        tutorialLink={{
          label: 'View Tutorials at HashiCorp Learn',
          href: 'https://learn.hashicorp.com/waypoint',
        }}
      />

      <div className="g-container">
        <div className={styles.gettingStarted}>
          <h2>Get Started</h2>
          <p>
            Follow step-by-step tutorials on AWS, Azure, GCP, and localhost.
          </p>
          <div className={styles.links}>
            {getStartedLinks.map((link) => (
              <a href={link.href} key={link.href}>
                {link.label}
              </a>
            ))}
          </div>
        </div>
      </div>

      <ReleaseInformation
        brand="blue"
        productId="waypoint"
        productName={productName}
        releases={allVersionsReleaseData}
        latestVersion={currentVersionReleaseData.version}
      />
    </div>
  )
}

function fetchVersionRelease(version) {
  return fetch(
    `https://releases.hashicorp.com/waypoint/${version}/index.json`,
    {
      headers: {
        'Cache-Control': 'no-cache',
      },
    }
  ).then((r) => r.json())
}

function fetchAllReleases() {
  return fetch(`https://releases.hashicorp.com/waypoint/index.json`, {
    headers: {
      'Cache-Control': 'no-cache',
    },
  })
    .then((res) => res.json())
    .then((data) => {
      const latestReleases = sortAndFilterReleases(Object.keys(data.versions))
      const releases = latestReleases.map((releaseVersion) => ({
        ...sortPlatforms(data.versions[releaseVersion]),
        version: releaseVersion,
      }))
      return releases
    })
}

export async function getStaticProps() {
  return Promise.all([fetchVersionRelease(VERSION), fetchAllReleases()])
    .then((result) => {
      return {
        props: {
          currentVersionReleaseData: result[0],
          allVersionsReleaseData: result[1],
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
