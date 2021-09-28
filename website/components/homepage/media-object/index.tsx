import classNames from 'classnames'
import Button from '@hashicorp/react-button'
import InlineSvg from '@hashicorp/react-inline-svg'
import s from './style.module.css'

export interface MediaObjectProps {
  icon: string
  heading: string
  description: string | React.ReactNode
  link?: {
    url: string
    text: string
  }
  stacked?: boolean
}

export default function MediaObject({
  icon,
  heading,
  description,
  link,
  stacked = false,
}: MediaObjectProps): JSX.Element {
  return (
    <div
      className={classNames(s.mediaObject, {
        [s.mediaObjectStacked]: stacked,
      })}
    >
      {icon && (
        <div className={s.mediaObjectIcon}>{<InlineSvg src={icon} />}</div>
      )}
      <div className={s.mediaObjectBody}>
        <h3 className={s.mediaObjectHeading}>{heading}</h3>
        <p className={s.mediaObjectDescription}>{description}</p>
        {link && (
          <div className={s.mediaObjectAnchor}>
            <Button
              url={link.url}
              title={link.text}
              theme={{
                variant: 'tertiary-neutral',
              }}
              linkType="inbound"
            />
          </div>
        )}
      </div>
    </div>
  )
}
