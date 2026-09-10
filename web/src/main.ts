import "htmx.org"

import "./firefly.css"

import { mountIslands } from "./islands"

console.log("firefly: bundle loaded")

document.addEventListener(
    "DOMContentLoaded",
    () => {
        console.log("firefly: DOM ready")
        mountIslands()
    },
    { once: true },
)
