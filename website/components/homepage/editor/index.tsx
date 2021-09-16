import CodeBlock from '@hashicorp/react-code-block'
import s from './style.module.css'

interface EditorProps {
  code: string
  note: string
}

export default function Editor({ code, note }: EditorProps): JSX.Element {
  return (
    <>
      <div className={s.editor}>
        <CodeBlock
          options={{
            lineNumbers: true,
          }}
          theme="light"
          language="go"
          code={code}
        />
      </div>
      <p className={s.editorNote}>{note}</p>
    </>
  )
}
