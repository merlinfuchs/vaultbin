import style from './PasteEditor.css';

export default function PasteEditor({disabled, value, onChange}) {
    return (
        <textarea disabled={disabled} spellCheck={false} className={style.textarea} value={value} onInput={e => onChange(e.target.value)}/>
    )
}