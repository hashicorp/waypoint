import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import { getStaticGenerationFunctions } from '@hashicorp/react-docs-page/server'

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

const { getStaticPaths, getStaticProps } = getStaticGenerationFunctions({
  basePath: basePath,
  localContentDir: CONTENT_DIR,
  navDataFile: NAV_DATA_FILE,
  product: productSlug,
  strategy: process.env.STATIC_GENERATION_STRATEGY || 'fs',
  revalidate: 10,
})

export { getStaticPaths, getStaticProps }
