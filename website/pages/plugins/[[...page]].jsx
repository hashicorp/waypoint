import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from '@hashicorp/react-docs-page/server'

const NAV_DATA_FILE = 'data/plugins-nav-data.json'
const CONTENT_DIR = 'content/plugins'
const basePath = 'plugins'
const product = { slug: productSlug, name: productName }

export default function DocsLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
    />
  )
}

export async function getStaticPaths() {
  return {
    fallback: false,
    paths: await generateStaticPaths({
      basePath,
      product,
      navDataFile: NAV_DATA_FILE,
      localContentDir: CONTENT_DIR,
    }),
  }
}

export async function getStaticProps({ params }) {
  return {
    props: await generateStaticProps({
      basePath,
      product,
      navDataFile: NAV_DATA_FILE,
      localContentDir: CONTENT_DIR,
      params,
    }),
    revalidate: 10,
  }
}
