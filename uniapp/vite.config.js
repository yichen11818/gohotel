import { cpSync, existsSync, mkdirSync } from 'node:fs'
import { relative, resolve } from 'node:path'
import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

const normalizeUrl = (value) => {
  if (typeof value !== 'string') return ''
  return value.trim().replace(/\/+$/, '')
}

function syncStaticAssets() {
  let projectRoot = ''
  let outputDir = ''

  return {
    name: 'sync-static-assets',
    configResolved(config) {
      projectRoot = config.root
      outputDir = resolve(projectRoot, config.build.outDir)
    },
    closeBundle() {
      const sourceDir = resolve(projectRoot, 'static')
      if (!existsSync(sourceDir) || !outputDir) {
        return
      }

      const targetDir = resolve(outputDir, 'static')
      mkdirSync(outputDir, { recursive: true })
      cpSync(sourceDir, targetDir, { recursive: true, force: true })
      console.log(
        `[sync-static-assets] copied ${relative(projectRoot, sourceDir)} to ${relative(projectRoot, targetDir)}`
      )
    }
  }
}

export default defineConfig(({ command, mode }) => {
  const apiBaseUrl = normalizeUrl(
    process.env.UNI_APP_API_BASE_URL ||
      process.env.VITE_UNI_APP_API_BASE_URL ||
      (command === 'serve' ? 'http://127.0.0.1:18080' : '')
  )

  return {
  define: {
    __APP_MODE__: JSON.stringify(mode || (command === 'build' ? 'production' : 'development')),
    __APP_API_BASE_URL__: JSON.stringify(apiBaseUrl)
  },
  resolve: {
    alias: {
      '@': resolve(__dirname)
    }
  },
  plugins: [uni(), syncStaticAssets()],
  server: {
    port: 5173,
    host: '0.0.0.0'
  }
}
})

















