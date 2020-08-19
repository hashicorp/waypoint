import './style.css'
import '@hashicorp/nextjs-scripts/lib/nprogress/style.css'

import NProgress from '@hashicorp/nextjs-scripts/lib/nprogress'
import createConsentManager from '@hashicorp/nextjs-scripts/lib/consent-manager'
import useAnchorLinkAnalytics from '@hashicorp/nextjs-scripts/lib/anchor-link-analytics'
import Router from 'next/router'
import HashiHead from '@hashicorp/react-head'
import Head from 'next/head'
import { ErrorBoundary } from '@hashicorp/nextjs-scripts/lib/bugsnag'
import MegaNav from '@hashicorp/react-mega-nav'
import ProductSubnav from '../components/subnav'
import Footer from 'components/footer'
import Error from './_error'
import { productName } from '../data/metadata'

NProgress({ Router })
const { ConsentManager, openConsentManager } = createConsentManager({
  preset: 'oss',
})

function App({ Component, pageProps }) {
  useAnchorLinkAnalytics()

  return (
    <ErrorBoundary FallbackComponent={Error}>
      <HashiHead
        is={Head}
        title={`${productName} by HashiCorp`}
        siteName={`${productName} by HashiCorp`}
        description="Write your description here"
        image="https://www.example.com/img/og-image.png"
        stylesheet={[
          {
            href:
              'https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700&display=swap',
          },
        ]}
        icon={[{ href: '/favicon.ico' }]}
        preload={[
          { href: '/fonts/klavika/medium.woff2', as: 'font' },
          { href: '/fonts/gilmer/light.woff2', as: 'font' },
          { href: '/fonts/gilmer/regular.woff2', as: 'font' },
          { href: '/fonts/gilmer/medium.woff2', as: 'font' },
          { href: '/fonts/gilmer/bold.woff2', as: 'font' },
          { href: '/fonts/metro-sans/book.woff2', as: 'font' },
          { href: '/fonts/metro-sans/regular.woff2', as: 'font' },
          { href: '/fonts/metro-sans/semi-bold.woff2', as: 'font' },
          { href: '/fonts/metro-sans/bold.woff2', as: 'font' },
          { href: '/fonts/dejavu/mono.woff2', as: 'font' },
        ]}
      />
      <MegaNav product={productName} />
      <ProductSubnav />
      <div className="content">
        <Component {...pageProps} />
      </div>
      <Footer openConsentManager={openConsentManager} />
      <ConsentManager />
    </ErrorBoundary>
  )
}

App.getInitialProps = async ({ Component, ctx }) => {
  let pageProps = {}

  if (Component.getInitialProps) {
    pageProps = await Component.getInitialProps(ctx)
  } else if (Component.isMDXComponent) {
    // fix for https://github.com/mdx-js/mdx/issues/382
    const mdxLayoutComponent = Component({}).props.originalType
    if (mdxLayoutComponent.getInitialProps) {
      pageProps = await mdxLayoutComponent.getInitialProps(ctx)
    }
  }

  return { pageProps }
}

export default App
