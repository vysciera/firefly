import { createApp } from "vue"

import Counter from "./islands/Counter.vue"

const islands = {
    counter: Counter,
}

type IslandName = keyof typeof islands

export function mountIslands(root: ParentNode = document): void {
    const elements: HTMLElement[] = []

    if (
        root instanceof HTMLElement &&
        root.matches("[data-vue]")
    ) {
        elements.push(root)
    }

    elements.push(
        ...root.querySelectorAll<HTMLElement>("[data-vue]"),
    )

    for (const element of elements) {
        if (element.dataset.vueMounted === "true") {
            continue
        }

        const name = element.dataset.vue as IslandName | undefined

        if (!name) {
            continue
        }

        const component = islands[name]

        if (!component) {
            console.warn(`firefly: unknown Vue island "${name}"`)
            continue
        }

        createApp(component).mount(element)

        element.dataset.vueMounted = "true"
    }
}
