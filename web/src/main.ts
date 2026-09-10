import "htmx.org"

import "./firefly.css"

import { mountIslands } from "./islands"

function boot(): void {
    console.log("firefly: boot")
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
