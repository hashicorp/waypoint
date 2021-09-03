import s from './style.module.css'

function Section({ children }) {
  return <section className={s.section}>{children}</section>
}

function SectionHeading({ children }) {
  return <h2 className={s.sectionHeading}>{children}</h2>
}

function SectionDescription({ children }) {
  return <p className={s.sectionDescription}>{children}</p>
}

export default Section
export { SectionHeading, SectionDescription }
