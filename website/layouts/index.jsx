import DocsPage from '@hashicorp/react-docs-page'
import Head from 'next/head'
import Link from 'next/link'
import { productName, productSlug } from 'data/metadata'

function DefaultLayoutWrapper(pageMeta) {
  function DefaultLayout(props) {
    return (
      <DocsPage
        {...props}
        product={productSlug}
        head={{
          is: Head,
          title: `${pageMeta.page_title} | ${productName} by HashiCorp`,
          description: pageMeta.description,
          siteName: `${productName} by HashiCorp`,
        }}
        sidenav={{
          Link,
          category: 'docs',
          currentPage: props.path,
          data: [],
          order: [],
          disableFilter: true,
        }}
        resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/pages/${pageMeta.__resourcePath}`}
      />
    )
  }

  DefaultLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return DefaultLayout
}

export default DefaultLayoutWrapper
