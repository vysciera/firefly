import "htmx.org"

import "./firefly.css"

import { mountIslands } from "./islands"

function boot(): void {
    mountIslands()
}

if (document.readyState === "loading") {
    document.addEventListener(
        "DOMContentLoaded",
        boot,
        { once: true },
    )
} else {
    boot()
}
