import DocsPage from '@hashicorp/react-docs-page'
import order from '../data/plugins-navigation.js'
import { productName, productSlug } from 'data/metadata'
import { frontMatter as data } from '../pages/plugins/**/*.mdx'
import Head from 'next/head'
import Link from 'next/link'
import { createMdxProvider } from '@hashicorp/nextjs-scripts/lib/providers/docs'

const MDXProvider = createMdxProvider({ product: productName })

function PluginsLayoutWrapper(pageMeta) {
  function PluginsLayout(props) {
    return (
      <MDXProvider>
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
            category: 'plugins',
            currentPage: props.path,
            data,
            order,
          }}
          resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/pages/${pageMeta.__resourcePath}`}
        />
      </MDXProvider>
    )
  }

  PluginsLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return PluginsLayout
}

export default PluginsLayoutWrapper
