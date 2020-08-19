import DocsPage from '@hashicorp/react-docs-page'
import order from 'data/docs-navigation.js'
import { productName, productSlug } from 'data/metadata'
import { frontMatter as data } from '../pages/docs/**/*.mdx'
import { createMdxProvider } from '@hashicorp/nextjs-scripts/lib/providers/docs'
import Head from 'next/head'
import Link from 'next/link'

const MDXProvider = createMdxProvider({ product: productName })

function DocsLayoutWrapper(pageMeta) {
  function DocsLayout(props) {
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
            category: 'docs',
            currentPage: props.path,
            data,
            order,
          }}
          resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/pages/${pageMeta.__resourcePath}`}
        />
      </MDXProvider>
    )
  }

  DocsLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return DocsLayout
}

export default DocsLayoutWrapper
