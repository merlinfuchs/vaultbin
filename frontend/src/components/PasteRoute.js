import style from "./PasteRoute.css";
import PasteEditor from "./PasteEditor";
import PasteMenu from "./PasteMenu";
import PasteStats from "./PasteStats"
import {useEffect, useRef, useState} from "preact/hooks";
import {route} from "preact-router";
import PasteViewer from "./PasteViewer";
import hljs from "highlight.js/lib/common";

const supportedLanguage = hljs.listLanguages().filter(hljs.getLanguage)

function guessProgrammingLanguage(value) {
    return new Promise(resolve => {
        const match = supportedLanguage.reduce(
            (previous, next) => {
                const result = hljs.highlight(value, {
                    language: next,
                    ignoreIllegals: false
                })

                if (result.relevance > previous.relevance) {
                    return {...result, language: next}
                }

                return previous
            },
            {relevance: 0, value}
        )

        resolve(match.language || null)
    })
}

export default function PasteRoute({pasteId}) {
    const [content, setContent] = useState("");
    const [viewCount, setViewCount] = useState(null);

    const detectLanguage = useRef(true)
    const [language, setLanguage] = useState(null);

    useEffect(() => {
        if (pasteId) {
            fetch(`/api/pastes/${pasteId}`)
                .then(async resp => {
                    if (resp.ok) {
                        const data = await resp.json()
                        setContent(data.data.content);
                        setViewCount(data.data.view_count)
                        setLanguage(data.data.language)
                        detectLanguage.current = false;
                    } else {
                        route("/")
                    }
                })
        } else {
            setViewCount(null);
            detectLanguage.current = true;
        }
    }, [pasteId])

    useEffect(() => {
        if (detectLanguage) {
            guessProgrammingLanguage(content).then(result => {
                setLanguage(result)
            })
        }
    }, [content])

    return (
        <div className={style.container} path="/:pasteId?">
            <div className={style.center}>
                {pasteId ? (
                    <PasteViewer content={content}/>
                ) : (
                    <PasteEditor disabled={!!pasteId} value={content} onChange={setContent}/>
                )}
            </div>
            <div className={style.menu}>
                <PasteMenu content={content} setContent={setContent} pasteId={pasteId} language={language}/>
            </div>
            <div className={style.stats}>
                <PasteStats viewCount={viewCount} language={language}/>
            </div>
        </div>
    )
}