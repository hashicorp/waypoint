import s from './style.module.css'

export default function NestedNode(props) {
  const childrenArray = props.children.split('.')
  const removed = childrenArray.splice(0, 1)
  return (
    <>
      <span className={s.root}>{removed}.</span>
      {childrenArray.join('.')}
    </>
  )
}
