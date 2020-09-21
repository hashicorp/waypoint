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
  staticProps: { mdxSource, data, frontMatter, pagePath },
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
        resourceURL={`https://github.com/hashicorp/${productSlug}/blob/master/website/content/docs/${frontMatter.__resourcePath}`}
      >
        {content}
      </DocsPageComponent>
    </>
  )
}
