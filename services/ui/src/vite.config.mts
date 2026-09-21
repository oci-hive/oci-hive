import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

const delayHmr = (ms: number) => ({
    name: 'delay-hmr',
    configureServer(server) {
        const originalSend = server.ws.send
        server.ws.send = function (...args: any[]) {
            setTimeout(() => originalSend.apply(this, args), ms)
        }
    },
})

export default defineConfig({
    plugins: [tailwindcss(), delayHmr(200)],

    server: {
        host: '0.0.0.0',
        port: 5173,
        hmr: {
            host: 'localhost',
        },
        forwardConsole: {
            unhandledErrors: true,
            logLevels: ['warn', 'error', 'debug', 'log', 'info'],
        },
    },

    build: {
        outDir: 'static/build',
        rollupOptions: {
            input: {
                main: './src/js/main.ts',
            },
            output: {
                entryFileNames: '[name].js',
                chunkFileNames: '[name].js',
                assetFileNames: '[name].[ext]',
            },
        },
    },

    resolve: {
        alias: {
            '@': '/src/js',
        },
    },
})
