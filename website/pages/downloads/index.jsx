import s from './style.module.css'
import VERSION from 'data/version.js'
import ProductDownloader from '@hashicorp/react-product-downloader'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'
import { productName, productSlug } from 'data/metadata'

export default function DownloadsPage({ releaseData }) {
  return (
    <div className={s.root}>
      <HashiHead is={Head} title={`Downloads | ${productName} by HashiCorp`} />
      <ProductDownloader
        product={productName}
        version={VERSION}
        releaseData={releaseData}
      />
    </div>
  )
}

export async function getStaticProps() {
  // NOTE: make sure to change "vault" here to your product slug
  return fetch(`https://releases.hashicorp.com/vault/${VERSION}/index.json`)
    .then((r) => r.json())
    .then((releaseData) => ({ props: { releaseData } }))
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
