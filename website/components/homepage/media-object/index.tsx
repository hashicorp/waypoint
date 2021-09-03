import classNames from 'classnames'
import InlineSvg from '@hashicorp/react-inline-svg'
import s from './style.module.css'

interface MediaObjectProps {
  icon: string
  heading: string
  description: string
  stacked?: boolean
}

export default function MediaObject({
  icon,
  heading,
  description,
  stacked = false,
}: MediaObjectProps) {
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
      </div>
    </div>
  )
}
