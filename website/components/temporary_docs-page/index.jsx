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

import DocsPageComponent from '@hashicorp/react-docs-page'
import Head from 'next/head'
import Link from 'next/link'
import hydrate from 'next-mdx-remote/hydrate'
import generateComponents from './components'

export default function DocsPage({
  productName,
  productSlug,
  subpath,
  order,
  staticProps: { mdxSource, data, frontMatter, pagePath, filePath },
}) {
  const content = hydrate(mdxSource, {
    components: generateComponents(productName),
  })

  return (
    <>
      <DocsPageComponent
        product={productSlug}
        head={{
          is: Head,
          title: `${frontMatter.page_title} | ${productName} by HashiCorp`,
          description: frontMatter.description,
          siteName: `${productName} by HashiCorp`,
        }}
        sidenav={{
          Link,
          category: subpath,
          currentPage: pagePath,
          data: data,
          order,
        }}
        resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/content/docs/${filePath}`}
      >
        {content}
      </DocsPageComponent>
    </>
  )
}
