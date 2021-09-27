import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from '@hashicorp/react-docs-page/server'
import Placement from 'components/placement-table'
import NestedNode from 'components/nested-node'

const NAV_DATA_FILE = 'data/docs-nav-data.json'
const CONTENT_DIR = 'content/docs'
const basePath = 'docs'
const additionalComponents = { Placement, NestedNode }

export default function DocsLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
      additionalComponents={additionalComponents}
      showVersionSelect={process.env.ENABLE_VERSIONED_DOCS}
    />
  )
}

export async function getStaticPaths() {
  const paths = await generateStaticPaths({
    navDataFile: NAV_DATA_FILE,
    localContentDir: CONTENT_DIR,
    // new ----
    product: { name: productName, slug: productSlug },
    basePath,
  })
  return {
    fallback: true,
    paths,
  }
}

export async function getStaticProps({ params }) {
  const props = await generateStaticProps({
    navDataFile: NAV_DATA_FILE,
    localContentDir: CONTENT_DIR,
    product: { name: productName, slug: productSlug },
    params,
    additionalComponents,
    basePath,
  })
  return {
    props,
    revalidate: 10,
  }
}
