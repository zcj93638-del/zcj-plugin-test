import * as mcp from '@deepseek-ai/dsh-mcp-client'
import { TypertRemoteService } from '@deepseek-ai/dsh-typert-protocol'
import contract from './client.cjs'
import { executable, callSettings } from './native.js'

export const name = 'trim-seek'
class SeekSettings extends TypertRemoteService {
  constructor(ctx) { super(ctx, 'trimSeekSettings', {namespace: 'trimSeek'}) }
  async read() {
    const {values} = await callSettings('seek_settings_read')
    const {configured} = await callSettings('seek_status')
    return {address: values.address, username: values.username, configured}
  }
  async login(fields) {
    await callSettings('seek_settings_update', {set: fields})
    return this.read()
  }
}
export function apply(ctx) {
  ctx.plugin(mcp, {transport: 'stdio', serverName: 'trim-seek', command: executable(), args: [], toolCallTimeoutMs: 1810000})
  ctx.inject(['typert'], sub => {
    new SeekSettings(sub)
    const dispose = sub.typert.register({package: contract.pkg, face: 'host', schemas: [],
      model: {services: [], objects: [], events: []}, invocations: contract.descriptors})
    sub.effect(() => () => dispose(), 'trim-seek settings')
  })
}
