import React from 'react'
import GradientBox from '../gradient-box'
import s from './style.module.css'

interface NumberedBlockProps {
  index: string
  heading: string
  children: React.ReactNode
}

export default function NumberedBlock({
  index,
  heading,
  children,
}: NumberedBlockProps) {
  return (
    <div className={s.numberedBlock}>
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
