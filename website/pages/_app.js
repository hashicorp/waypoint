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
      <Footer
        openConsentManager={openConsentManager}
        heading="Using Waypoint"
        description="The best way to understand what Waypoint can enable for your projects is to give it a try."
        cards={[
          {
            link:
              'https://learn.hashicorp.com/collections/waypoint/get-started-kubernetes',
            img: '/img/get-started-kubernetes.png',
            eyebrow: 'Tutorial',
            title: 'Get Started - Kubernetes',
            description:
              'Build, deploy, and release applications to a Kubernetes cluster.',
          },
          {
            link:
              'https://learn.hashicorp.com/tutorials/waypoint/get-started-intro',
            img: '/img/intro-to-waypoint.png',
            eyebrow: 'Tutorial',
            title: 'Introduction to Waypoint',
            description:
              'Waypoint enables you to publish any application to any platform with a single file and a single command.',
          },
        ]}
        ctaLinks={[
          {
            text: 'Waypoint tutorials',
            url: 'https://learn.hashicorp.com/waypoint',
          },
          {
            text: 'Waypoint documentation',
            url: '/docs',
          },
        ]}
        navLinks={[
          {
            text: 'Documentation',
            url: '/docs',
          },
          {
            text: 'API Reference',
            url: '/',
          },
          {
            text: 'Tutorials',
            url: 'https://learn.hashicorp.com/waypoint',
          },
          {
            text: 'Integrations',
            url: '/',
          },
        ]}
      />
      <ConsentManager />
    </ErrorBoundary>
  )
}
