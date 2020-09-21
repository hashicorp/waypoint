import defaultMdxComponents from '@hashicorp/nextjs-scripts/lib/providers/docs'

export default function generateComponents(productName) {
  return defaultMdxComponents({ product: productName })
}
