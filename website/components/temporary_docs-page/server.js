// -----------------------------------------------------
//                This code is LOCKED
//
// If any changes are needed to this code, or if this code
// is needed in any other projects, instead of changing or
// using it, instead we must complete this task as a prerequisite
//
// https://app.asana.com/0/1100423001970639/1195001770724993
//
// ------------------------------------------------------

import fs from 'fs'
import path from 'path'
import existsSync from 'fs-exists-sync'
import readdirp from 'readdirp'
import lineReader from 'line-reader'
import matter from 'gray-matter'
import { safeLoad } from 'js-yaml'
import renderToString from 'next-mdx-remote/render-to-string'
import markdownDefaults from '@hashicorp/nextjs-scripts/markdown'
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

export async function getStaticMdxPaths(root) {
  const files = await readdirp.promise(root, { fileFilter: ['*.mdx'] })

  return files.map(({ path: p }) => {
    return {
      params: {
        page: p
          .replace(/\.mdx$/, '')
          .split('/')
          .filter((p) => p !== 'index'),
      },
    }
  })
}

export async function renderPageMdx(root, pagePath, components) {
  // get the page being requested - figure out if its index page or leaf
  // prefer leaf if both are present
  const leafPath = path.join(root, `${pagePath}.mdx`)
  const indexPath = path.join(root, `${pagePath}/index.mdx`)
  let page

  if (existsSync(leafPath)) {
    page = fs.readFileSync(leafPath, 'utf8')
  } else if (existsSync(indexPath)) {
    page = fs.readFileSync(indexPath, 'utf8')
  } else {
    // NOTE: if we decide to let docs pages render dynamically, we should replace this
    // error with a straight 404, at least in production.
    throw new Error(
      `We went looking for "${leafPath}" and "${indexPath}" but neither one was found.`
    )
  }

  const { data: frontMatter, content } = matter(page)
  const mdxSource = await renderToString(content, {
    mdxOptions: markdownDefaults({
      resolveIncludes: path.join(process.cwd(), 'content/partials'),
    }),
    components,
  })

  return { mdxSource, frontMatter }
}

export function fastReadFrontMatter(p) {
  return new Promise((resolve) => {
    const fm = []
    readdirp(p, { fileFilter: '*.mdx' })
      .on('data', (entry) => {
        let lineNum = 0
        const content = []
        fm.push(
          new Promise((resolve2, reject2) => {
            lineReader.eachLine(
              entry.fullPath,
              (line) => {
                // if it has any content other than `---`, the file doesn't have front matter, so we close
                if (lineNum === 0 && !line.match(/^---$/)) return false
                // if it's not the first line and we have a bottom delimiter, exit
                if (lineNum !== 0 && line.match(/^---$/)) return false
                // now we read lines until we match the bottom delimiters
                content.push(line)
                // increment line number
                lineNum++
              },
              (err) => {
                if (err) return reject2(err)
                content.push(`__resourcePath: "${entry.path}"`)
                resolve2(safeLoad(content.slice(1).join('\n')), {
                  filename: entry.fullPath,
                })
              }
            )
          })
        )
      })
      .on('end', () => {
        Promise.all(fm).then((res) => resolve(res))
      })
  })
}
