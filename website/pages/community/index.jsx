import styles from './style.module.css'
import classNames from 'classnames'
import VerticalTextBlockList from '@hashicorp/react-vertical-text-block-list'
import SectionHeader from '@hashicorp/react-section-header'
import Head from 'next/head'

export default function CommunityPage() {
  return (
    <div className={styles.communityPage}>
      <Head>
        <title key="title">Community | Waypoint by HashiCorp</title>
      </Head>
      <div
        className={classNames(styles.sectionHeaderWrapper, 'g-grid-container')}
      >
        <SectionHeader
          headline="Community"
          description="Waypoint is a newly-launched open source project. The project team depends on the community’s engagement and feedback. Get involved today."
          use_h1={true}
        />
      </div>
      <VerticalTextBlockList
        product="waypoint"
        data={[
          {
            header: 'Community Forum',
            body:
              '<a href="https://discuss.hashicorp.com/c/waypoint">Waypoint Community Forum</a>',
          },
          {
            header: 'Bug Tracker',
            body:
              '<a href="https://github.com/hashicorp/waypoint/issues">Issue tracker on GitHub</a>. Please only use this for reporting bugs. Do not ask for general help here; use the Community Form for that.',
          },
        ]}
      />
    </div>
  )
}
