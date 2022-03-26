import style from './PasteMenu.scss';
import {route} from "preact-router";
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {faFileArrowUp, faFileCirclePlus, faFileLines, faFilePen} from "@fortawesome/free-solid-svg-icons";
import {faGithub} from "@fortawesome/free-brands-svg-icons";
import {useEffect} from "preact/hooks";

export default function PasteMenu({pasteId, content, setContent, language}) {
    const canSave = !pasteId && content.length !== 0

    function onSave() {
        if (!canSave) return

        fetch("/api/pastes", {
            method: "POST",
            body: JSON.stringify({
                content,
                language
            }),
            headers: {"Content-Type": "application/json"}
        })
            .then(async resp => {
                if (resp.ok) {
                    const data = await resp.json()
                    route(`/${data.data.id}`)
                } else {
                    alert("Failed to create paste :(")
                }
            })
    }

    function onCreateNew() {
        setContent("")
        route("/")
    }

    function onDuplicate() {
        if (!pasteId) return
        route("/")
    }

    function onOpenRaw() {
        if (!pasteId) return
        window.location.href = `/api/pastes/${pasteId}/raw`;
    }

    useEffect(() => {
        function onKeyDown(e) {
            if (!e.ctrlKey) return;
            console.log(e.key)
            switch (e.key) {
                case "s":
                    e.preventDefault();
                    onSave();
                    break;
                case "n":
                    e.preventDefault();
                    onCreateNew();
                    break;
                case "d":
                    e.preventDefault();
                    onDuplicate()
                    break;
                case "r":
                    e.preventDefault();
                    if (e.shiftKey) onOpenRaw();
                    break;
                case "R":
                    e.preventDefault();
                    if (e.shiftKey) onOpenRaw();
                    break;
            }
        }

        document.addEventListener("keydown", onKeyDown);
        return () => document.removeEventListener("keydown", onKeyDown)
    }, [pasteId, content, setContent, language])

    return (
        <div className={style.menu}>
            <div className={style.button}>
                <FontAwesomeIcon icon={faFileArrowUp} alt=""
                                 className={`${!canSave ? style.disabled : ''}`}
                                 onClick={onSave} style={{paddingRight: "7px"}}/>
                <div className={style.description}>
                    <div className={style.cleanup}/>
                    <div className={style.inner}>
                        <div className={style.title}>Save</div>
                        <div className={style.shortcut}>ctrl + s</div>
                    </div>
                </div>
            </div>
            <div className={style.button}>
                <FontAwesomeIcon icon={faFileCirclePlus} alt="" onClick={onCreateNew}/>
                <div className={style.description}>
                    <div className={style.cleanup}/>
                    <div className={style.inner}>
                        <div className={style.title}>New</div>
                        <div className={style.shortcut}>ctrl + n</div>
                    </div>
                </div>
            </div>
            <div className={style.button}>
                <FontAwesomeIcon icon={faFilePen} alt="" className={`${!pasteId ? style.disabled : ''}`}
                                 onClick={onDuplicate}/>
                <div className={style.description}>
                    <div className={style.cleanup}/>
                    <div className={style.inner}>
                        <div className={style.title}>Duplicate</div>
                        <div className={style.shortcut}>ctrl + d</div>
                    </div>
                </div>
            </div>
            <div className={style.button}>
                <FontAwesomeIcon icon={faFileLines} alt="" className={`${!pasteId ? style.disabled : ''}`}
                                 onClick={onOpenRaw}/>
                <div className={style.description}>
                    <div className={style.cleanup}/>
                    <div className={style.inner}>
                        <div className={style.title}>Raw</div>
                        <div className={style.shortcut}>ctrl + shift + r</div>
                    </div>
                </div>
            </div>
            <a className={style.button} href="https://github.com/merlinfuchs/vaultbin" target="_blank">
                <FontAwesomeIcon icon={faGithub} alt="" style={{paddingLeft: "6px"}}/>
                <div className={style.description}>
                    <div className={style.cleanup}/>
                    <div className={style.inner}>
                        <div className={style.title}>Github</div>
                        <div className={style.shortcut}>Support this project</div>
                    </div>
                </div>
            </a>
        </div>
    )
}