import "htmx.org"
import "./firefly.css"

import { mountIslands } from "./islands"

function boot() {
    mountIslands()
}

if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot)
} else {
    boot()
}

document.addEventListener("htmx:after:swap", (event) => {
    const customEvent = event as CustomEvent
    const target =
        customEvent.detail?.target instanceof Element
            ? customEvent.detail.target
            : document

    mountIslands(target)
})
