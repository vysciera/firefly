import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import { fileURLToPath } from "node:url"

export default defineConfig({
    plugins: [
        vue(),
    ],

    define: {
        "process.env.NODE_ENV": JSON.stringify("production"),
    },

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
