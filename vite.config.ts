import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import { fileURLToPath } from "node:url"

export default defineConfig({
    plugins: [
        vue(),
    ],

    build: {
        outDir: "web/dist",
        emptyOutDir: true,

        lib: {
            entry: fileURLToPath(
                new URL("./web/src/main.ts", import.meta.url),
            ),
            formats: ["es"],
            fileName: "firefly",
            cssFileName: "firefly",
        },
    },
})
