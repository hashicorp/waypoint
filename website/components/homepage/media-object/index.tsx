import classNames from 'classnames'
import InlineSvg from '@hashicorp/react-inline-svg'
import s from './style.module.css'

interface MediaObjectProps {
  icon: string
  heading: string
  description: string
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
        {link && (
          <div className={s.mediaObjectAnchor}>
            <a className={s.mediaObjectAnchorLink} href={link.url}>
              {link.text}
            </a>
            <RightArrowIcon />
          </div>
        )}
      </div>
    </div>
  )
}

function RightArrowIcon() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 20 20"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M3.334 10h13.333M11.666 5l5 5-5 5"
        stroke="currentColor"
        strokeWidth="1.5"
      />
    </svg>
  )
}
