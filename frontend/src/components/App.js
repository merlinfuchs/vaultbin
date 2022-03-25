import Router from "preact-router";
import PasteRouter from "./PasteRoute";

export default function App() {
    return (
        <Router>
            <PasteRouter path="/:pasteId?"/>
        </Router>
    )
}
