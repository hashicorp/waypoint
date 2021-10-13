import Image, { ImageProps } from 'next/image'
import CallToAction from '@hashicorp/react-call-to-action'
import s from './style.module.css'

interface SectionCallToActionProps {
  features: Array<{
    media: ImageProps
    text: string | React.ReactNode
  }>
  heading: string
  content: string
  links: Array<{
    text: string
    url: string
  }>
}

export default function SectionCallToAction({
  features,
  heading,
  content,
  links,
}: SectionCallToActionProps): JSX.Element {
  return (
    <section className={s.root}>
      <ul className={s.featureList}>
        {features.map((feature, idx) => {
          const { media, text } = feature
          return (
            // Index is stable
            // eslint-disable-next-line react/no-array-index-key
            <li key={idx} className={s.feature}>
              <div className={s.featureMedia}>
                {media && <Image {...media} />}
              </div>
              <p className={s.featureText}>{text}</p>
            </li>
          )
        })}
      </ul>
      <CallToAction
        heading={heading}
        content={content}
        links={[...links]}
        product="waypoint"
        variant="compact"
      />
    </section>
  )
}
