import style from './PasteLineMarker.scss';

export default function PasteLineMarker({lineCount}) {
    const lineNumbers = []
    for (let i = 0; i < lineCount; i++) {
        lineNumbers.push(i)
    }

    return (
        <div className={style.marker}>
            {
                lineNumbers.map((_, i) => (
                    <div key={i}>{i}</div>
                ))
            }
        </div>
    )
}