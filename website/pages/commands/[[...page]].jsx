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
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
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
  try {
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
  } catch (err) {
    console.log('xxxxxxxxxxxxxxxxxxx')
    console.log('Failed to generate static props:', params.page, err.message)
    console.log('xxxxxxxxxxxxxxxxxxx')
    return { notFound: true }
  }
}
