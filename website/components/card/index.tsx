import s from './style.module.css'

export interface CardProps {
  link: string
  img?: string
  eyebrow?: string
  title: string
  description: string
}

export default function Card({
  link,
  img,
  eyebrow,
  title,
  description,
}: CardProps) {
  return (
    <article className={s.card}>
      <a className={s.link} href={link}>
        {img && (
          <div className={s.media}>
            <img src={img} alt="" />
          </div>
        )}
        <div className={s.content}>
          {eyebrow && <span className={s.eyebrow}>{eyebrow}</span>}
          <h3 className={s.title}>{title}</h3>
          <p className={s.description}>{description}</p>
          <span className={s.fauxLink}>
            View <RightArrowIcon />
          </span>
        </div>
      </a>
    </article>
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
        d="M3.33333 10H16.6667"
        stroke="black"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M11.6667 5L16.6667 10L11.6667 15"
        stroke="black"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}
