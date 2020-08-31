import { productName, productSlug } from 'data/metadata'
import order from 'data/plugins-navigation.js'
import DocsPage from 'components/new-docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from 'components/new-docs-page/server'

const subpath = 'plugins'

function DocsLayout(props) {
  return (
    <DocsPage
      productName={productName}
      productSlug={productSlug}
      subpath={subpath}
      order={order}
      staticProps={props}
    />
  )
}

export async function getStaticPaths() {
  return generateStaticPaths(subpath)
}

export async function getStaticProps({ params }) {
  return generateStaticProps(subpath, productName, params)
}

export default DocsLayout
