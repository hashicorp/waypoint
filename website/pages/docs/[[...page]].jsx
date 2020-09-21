import { productName, productSlug } from 'data/metadata'
import order from 'data/docs-navigation.js'
import DocsPage from 'components/temporary_docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from 'components/temporary_docs-page/server'

const subpath = 'docs'

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
