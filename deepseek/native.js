import { execFile } from 'node:child_process'
import { fileURLToPath } from 'node:url'

export function executable(platform = process.platform, arch = process.arch) {
  const cpu = {x64: 'amd64', arm64: 'arm64'}[arch]
  if (!cpu || !['darwin', 'linux', 'win32'].includes(platform)) {
    throw new Error(`Trim Seek 不支持 ${platform}/${arch}`)
  }
  const file = platform === 'win32'
    ? `seek-windows/bin/seek-mcp-windows-${cpu}.exe`
    : `seek/bin/seek-mcp-${platform}-${cpu}`
  return fileURLToPath(new URL(`../plugins/${file}`, import.meta.url))
}

// The settings page uses a one-shot request. Credentials only go over stdin.
export function callSettings(tool, args = {}) {
  return new Promise((resolve, reject) => {
    const child = execFile(executable(), [], {timeout: 15000, windowsHide: true}, (error, output) => {
      if (error) return reject(new Error('Seek 程序未能完成请求，请稍后重试。'))
      try {
        const {result} = JSON.parse(output)
        if (result.isError) throw new Error(result.content[0].text)
        resolve(result.structuredContent)
      } catch (err) { reject(err) }
    })
    child.stdin.on('error', reject)
    child.stdin.end(JSON.stringify({jsonrpc: '2.0', id: 1, method: 'tools/call', params: {name: tool, arguments: args}}) + '\n')
  })
}
