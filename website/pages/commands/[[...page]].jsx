import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from '@hashicorp/react-docs-page/server'

const NAV_DATA_FILE = 'data/commands-nav-data.json'
const CONTENT_DIR = 'content/commands'
const basePath = 'commands'

export default function DocsLayout(props) {
  console.log(process.env.ENABLE_VERSIONED_DOCS)
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
    />
  )
}

export async function getStaticPaths() {
  const paths = await generateStaticPaths({
    navDataFile: NAV_DATA_FILE,
    localContentDir: CONTENT_DIR,
    // new ----
    product: { name: productName, slug: productSlug },
    currentVersion: 'latest',
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
    basePath,
  })
  return {
    props,
  }
}
