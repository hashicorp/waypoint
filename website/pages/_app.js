import './style.css'
import '@hashicorp/nextjs-scripts/lib/nprogress/style.css'
import NProgress from '@hashicorp/nextjs-scripts/lib/nprogress'
import useAnchorLinkAnalytics from '@hashicorp/nextjs-scripts/lib/anchor-link-analytics'
import Router from 'next/router'
import HashiHead from '@hashicorp/react-head'
import Head from 'next/head'
import createConsentManager from '@hashicorp/nextjs-scripts/lib/consent-manager'
import { ErrorBoundary } from '@hashicorp/nextjs-scripts/lib/bugsnag'
import { Provider as NextAuthProvider } from 'next-auth/client'
import ProductSubnav from 'components/subnav'
import Footer from '../components/footer'
import AuthIndicator from 'components/auth-indicator'
import AuthGate from 'components/auth-gate'
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
      <ConditionalAuthProvider session={pageProps.session}>
        <HashiHead
          is={Head}
          title={`${productName} by HashiCorp`}
          siteName={`${productName} by HashiCorp`}
          description="Waypoint allows developers to define their application build, deploy, and release lifecycle as code, with a consistent 'waypoint up' workflow."
          image="https://www.waypointproject.io/img/og-image.png"
          icon={[{ href: '/favicon.svg' }]}
        />
        <ProductSubnav />
        <div className="content">
          <Component {...pageProps} />
        </div>
        <Footer openConsentManager={openConsentManager} />
        <ConsentManager />
      </ConditionalAuthProvider>
    </ErrorBoundary>
  )
}

function ConditionalAuthProvider({ children, session }) {
  return process.env.HASHI_ENV === 'production' ? (
    <NextAuthProvider session={session}>
      <AuthGate>
        {children}
        <AuthIndicator />
      </AuthGate>
    </NextAuthProvider>
  ) : (
    <>{children}</>
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
