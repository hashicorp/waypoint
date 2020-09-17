import DocsPage from '@hashicorp/react-docs-page'
import order from '../data/commands-navigation.js'
import { frontMatter as data } from '../pages/commands/**/*.mdx'
import { productName, productSlug } from 'data/metadata'
import Head from 'next/head'
import Link from 'next/link'
import { createMdxProvider } from '@hashicorp/nextjs-scripts/lib/providers/docs'

const MDXProvider = createMdxProvider({ product: productName })

function CommandsLayoutWrapper(pageMeta) {
  function CommandsLayout(props) {
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
            category: 'commands',
            currentPage: props.path,
            data,
            order,
          }}
          resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/pages/${pageMeta.__resourcePath}`}
        />
      </MDXProvider>
    )
  }

  CommandsLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return CommandsLayout
}

export default CommandsLayoutWrapper
