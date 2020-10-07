import Subnav from '@hashicorp/react-subnav'
import { useRouter } from 'next/router'
import Link from 'next/link'
import subnavItems from 'data/navigation'
import { productSlug } from 'data/metadata'

// A regex to match the pathname that comes out of router.pathname
// when the path is a content page (e.g. /docs/[[...page]]) to strip
// out the part that we're not interested in.
const multipathRegex = /\/\[\[.*\]\]/g

export default function ProductSubnav() {
  const router = useRouter()
  return (
    <Subnav
      titleLink={{
        text: 'Waypoint',
        url: '/',
      }}
      ctaLinks={[
        {
          text: 'GitHub',
          url: `https://www.github.com/hashicorp/${productSlug}`,
        },
        {
          text: 'Download',
          url: '/downloads',
        },
      ]}
      currentPath={router.pathname.replace(multipathRegex, '')}
      menuItemsAlign="right"
      menuItems={subnavItems}
      constrainWidth
      Link={Link}
    />
  )
}
