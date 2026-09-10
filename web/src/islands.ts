import { createApp } from "vue"
import Counter from "./islands/Counter.vue"

const islands = {
    counter: Counter,
}

type IslandName = keyof typeof islands

export function mountIslands(root: ParentNode = document) {
    const elements = Array.from(
        root.querySelectorAll<HTMLElement>("[data-vue]")
    )

    for (const element of elements) {
        if (element.dataset.vueMounted === "true") {
            continue
        }

        const name = element.dataset.vue as IslandName
        const component = islands[name]

        if (!component) {
            console.warn(`unknown Vue island: ${name}`)
            continue
        }

        createApp(component).mount(element)
        element.dataset.vueMounted = "true"
    }
}
