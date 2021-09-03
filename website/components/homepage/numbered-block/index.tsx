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
      <div className={s.numberedBlockIndex}>
        <GradientBox>{index}</GradientBox>
      </div>
      <div className={s.numberedBlockBody}>
        <h2 className={s.numberedBlockHeading}>{heading}</h2>
        {children}
      </div>
    </div>
  )
}
