import { createApp } from "vue"

import Counter from "./islands/Counter.vue"

export function mountIslands(): void {
    const element =
        document.querySelector<HTMLElement>('[data-vue="counter"]')

    console.log("firefly: island element", element)

    if (!element) {
        return
    }

    createApp(Counter).mount(element)

    console.log("firefly: counter mounted")
}
