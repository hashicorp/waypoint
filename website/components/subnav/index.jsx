import Subnav from '@hashicorp/react-subnav'
import { useRouter } from 'next/router'
import subnavItems from 'data/navigation'
import { productSlug } from 'data/metadata'

export default function ProductSubnav() {
  const router = useRouter()
  return (
    <Subnav
      titleLink={{
        text: productSlug,
        url: '/',
      }}
      ctaLinks={[
        {
          text: 'GitHub',
          url: `https://www.github.com/hashicorp/${productSlug}`,
        },
        { text: 'Download', url: '/downloads' },
      ]}
      currentPath={router.pathname}
      menuItemsAlign="right"
      menuItems={subnavItems}
      constrainWidth
    />
  )
}
