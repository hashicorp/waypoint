import VerticalTextBlockList from '@hashicorp/react-vertical-text-block-list'
import SectionHeader from '@hashicorp/react-section-header'
import Head from 'next/head'
import { productName, productSlug } from 'data/metadata'

export default function CommunityPage() {
  return (
    <div id="p-community">
      <Head>
        <title key="title">Community | {productName} by HashiCorp</title>
      </Head>
      <SectionHeader
        headline="Community"
        description={`${productName} is an open source project with a growing community. There are active, dedicated users willing to help you through various mediums.`}
        use_h1={true}
      />
      <VerticalTextBlockList
        data={[
          {
            header: 'IRC',
            body: `#${productSlug} on freenode`,
          },
          {
            header: 'Announcement List',
            body:
              '[HashiCorp Announcement Google Group](https://groups.google.com/group/hashicorp-announce)',
          },
          {
            header: 'Bug Tracker',
            body: `[Issue tracker on GitHub](https://github.com/hashicorp/${productSlug}/issues). Please only use this for reporting bugs. Do not ask for general help here. Use IRC or the mailing list for that.`,
          },
          {
            header: 'Training',
            body:
              'Paid [HashiCorp training courses](https://www.hashicorp.com/training) are also available in a city near you. Private training courses are also available.',
          },
        ]}
      />
    </div>
  )
}
