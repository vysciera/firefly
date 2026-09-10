import "htmx.org"
import "./firefly.css"

import { mountIslands } from "./islands"

function boot() {
    mountIslands()
}

if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot, { once: true })
} else {
    boot()
}

document.addEventListener("htmx:after:swap", (event) => {
    const detail = (event as CustomEvent).detail
    const target = detail?.ctx?.target

    if (target instanceof HTMLElement) {
        mountIslands(target)
    }
})
