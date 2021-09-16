import React from 'react'
import classNames from 'classnames'
import GradientBox from '../gradient-box'
import s from './style.module.css'

interface NumberedBlockProps {
  className?: string
  index: string
  heading: string
  children: React.ReactNode
}

export default function NumberedBlock({
  className,
  index,
  heading,
  children,
}: NumberedBlockProps) {
  return (
    <div className={classNames(s.numberedBlock, className)}>
      <header className={s.numberedBlockHeader}>
        <div className={s.numberedBlockIndex}>
          <GradientBox>{index}</GradientBox>
        </div>
        <h2 className={s.numberedBlockHeading}>{heading}</h2>
      </header>
      <div className={s.numberedBlockBody}>{children}</div>
    </div>
  )
}
