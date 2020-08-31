import path from 'path'
import {
  getStaticMdxPaths,
  fastReadFrontMatter,
  renderPageMdx,
} from 'lib/mdx-remote-docs'
import generateComponents from './components'

export async function generateStaticPaths(subpath) {
  const paths = await getStaticMdxPaths(
    path.join(process.cwd(), 'content', subpath)
  )

  return { paths, fallback: false }
}

export async function generateStaticProps(subpath, productName, params) {
  const docsPath = path.join(process.cwd(), 'content', subpath)
  const pagePath = params.page ? params.page.join('/') : '/'

  // get frontmatter from all other pages in the category, for the sidebar
  const allFrontMatter = await fastReadFrontMatter(docsPath)

  // render the current page path markdown
  const { mdxSource, frontMatter } = await renderPageMdx(
    docsPath,
    pagePath,
    generateComponents(productName)
  )

  return {
    props: {
      data: allFrontMatter.map((p) => {
        p.__resourcePath = `docs/${p.__resourcePath}`
        return p
      }),
      mdxSource,
      frontMatter,
      pagePath: `/docs/${pagePath}`,
    },
  }
}
