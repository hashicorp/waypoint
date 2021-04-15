import BaseTerminal from '@hashicorp/react-command-line-terminal'
import s from './style.module.css'

export default function Terminal({ title, lines }) {
  return (
    <div className={s.terminalWrapper}>
      <BaseTerminal product="waypoint" lines={lines} title={title} />
    </div>
  )
}
