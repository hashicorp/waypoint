import './style.css'
import '@hashicorp/platform-util/nprogress/style.css'
import NProgress from '@hashicorp/platform-util/nprogress'
import useAnchorLinkAnalytics from '@hashicorp/platform-util/anchor-link-analytics'
import Router from 'next/router'
import HashiHead from '@hashicorp/react-head'
import HashiStackMenu from '@hashicorp/react-hashi-stack-menu'
import AlertBanner from '@hashicorp/react-alert-banner'
import createConsentManager from '@hashicorp/react-consent-manager/loader'
import { ErrorBoundary } from '@hashicorp/platform-runtime-error-monitoring'
import ProductSubnav from 'components/subnav'
import Footer from '../components/footer'
import Error from './_error'
import alertBannerData, { ALERT_BANNER_ACTIVE } from 'data/alert-banner'

NProgress({ Router })
const { ConsentManager, openConsentManager } = createConsentManager({
  preset: 'oss',
})

const title = 'Waypoint by HashiCorp'
const description =
  'Waypoint is an open source solution that provides a modern workflow for build, deploy, and release across platforms.'

export default function App({ Component, pageProps }) {
  useAnchorLinkAnalytics()

  return (
    <ErrorBoundary FallbackComponent={Error}>
      <HashiHead
        title={title}
        siteName={title}
        description={description}
        image="https://www.waypointproject.io/img/og-image.png"
        icon={[{ href: '/_favicon.ico' }]}
      >
        <meta name="og:title" property="og:title" content={title} />
        <meta name="og:description" property="og:title" content={description} />
      </HashiHead>
      {ALERT_BANNER_ACTIVE && (
        <AlertBanner {...alertBannerData} product="waypoint" hideOnMobile />
      )}
      <HashiStackMenu />
      <ProductSubnav />
      <div className="content">
        <Component {...pageProps} />
      </div>
      <Footer openConsentManager={openConsentManager} />
      <ConsentManager />
    </ErrorBoundary>
  )
}
