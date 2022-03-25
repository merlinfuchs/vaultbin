import style from './PasteStats.scss';
import {faEye} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {useMemo} from "preact/compat";
import hljs from "highlight.js/lib/common";

export default function PasteStats({viewCount, language}) {
    if (!viewCount && !language) return <div />

    const languageObj = useMemo(() => hljs.getLanguage(language), [language])

    return (
        <div className={style.stats}>
            {language ? (
                <div className={style.stat}>{languageObj ? languageObj.name : language}</div>
            ) : <div />}
            {viewCount ? (
                <div className={style.stat}>
                    <FontAwesomeIcon icon={faEye}/>
                    <span>{viewCount}</span>
                </div>
            ) : <div />}
        </div>
    )
}