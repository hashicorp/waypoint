import DocsPage from '@hashicorp/react-docs-page'
import order from 'data/docs-navigation.js'
import { productName, productSlug } from 'data/metadata'
import { frontMatter as data } from '../pages/docs/**/*.mdx'
import { createMdxProvider } from '@hashicorp/nextjs-scripts/lib/providers/docs'
import Head from 'next/head'
import Link from 'next/link'

const MDXProvider = createMdxProvider({ product: productName })

function DocsLayout(props) {
  return (
    <MDXProvider>
      <DocsPage
        {...props}
        product={productSlug}
        head={{
          is: Head,
          title: `${props.frontMatter.page_title} | ${productName} by HashiCorp`,
          description: props.frontMatter.description,
          siteName: `${productName} by HashiCorp`,
        }}
        sidenav={{
          Link,
          category: 'docs',
          currentPage: props.path,
          data,
          order,
        }}
        resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/pages/${props.frontMatter.__resourcePath}`}
      />
    </MDXProvider>
  )
}

DocsLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

export default DocsLayout
