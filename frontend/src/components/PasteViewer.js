import style from "./PasteViewer.scss"
import hljs from "highlight.js/lib/common";
import {useMemo} from "preact/hooks";
import 'highlight.js/styles/a11y-dark.css';

export default function PasteViewer({content}) {
    const html = useMemo(() => {
        const highlighted = hljs.highlightAuto(content);
        // this wraps each line a div and adds the line numbers
        // it's done this way because a line can wrap and take more space
        return highlighted.value
            .split(/\r\n|\r|\n/)
            .map((line, i) => (
                <div className={style.line}>
                    <div className={style.lineNumber}>{i + 1}</div>
                    <div className={style.lineCode} dangerouslySetInnerHTML={{__html: line}}/>
                </div>
            ));
    }, [content])

    return (
        <div className={style.container}>
            {html}
        </div>
    )
}