// Harness loads the browser entry with its module loader. The host imports the
// same endpoint descriptors so field validation and wire names stay in sync.
const pkg = 'dsh-trim-seek'
function recordCodec(name, types) {
  return {mode: 'strict', typeSymbol: `${pkg}#${name}`, create: () => ({parse(value) {
    const result = {}
    for (const [key, type] of Object.entries(types)) {
      if (!value || typeof value[key] !== type) throw new Error(`无效的设置字段 ${key}`)
      result[key] = value[key]
    }
    return result
  }})}
}
const settings = recordCodec('Settings', {address: 'string', username: 'string', configured: 'boolean'})
const credentials = recordCodec('Credentials', {address: 'string', username: 'string', password: 'string'})
const descriptors = ['read', 'login'].map(method => ({
  id: `${pkg}#trimSeek/${method}`, service: 'trimSeekSettings', namespace: 'trimSeek', method,
  invocation: {kind: 'direct'}, result: settings,
  parameters: method === 'read' ? [] : [{name: 'fields', wire: 'fields', source: 'json', codec: credentials}],
}))

function client(require) {
  const {createElement: h, useState, useEffect} = require('react')
  function Settings({invoke}) {
    const [address, setAddress] = useState('')
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [message, show] = useState('正在读取连接…')
    const [busy, setBusy] = useState(true)
    useEffect(() => {
      let live = true
      invoke('read').then(value => {
        if (!live) return
        setAddress(value.address); setUsername(value.username); show(value.configured ? '已配置登录会话' : '')
      }).catch(err => { if (live) show(err.message) }).finally(() => { if (live) setBusy(false) })
      return () => { live = false }
    }, [invoke])
    async function submit(event) {
      event.preventDefault(); setBusy(true); show('正在登录…')
      try {
        const value = await invoke('login', {address, username, password})
        setAddress(value.address); show('登录成功')
      } catch (err) { show(err.message) }
      finally { setPassword(''); setBusy(false) }
    }
    function field(title, value, change, options) {
      return h('label', {style: {display: 'grid', gap: 6}}, title,
        h('input', {required: true, disabled: busy, value, onChange: e => change(e.target.value),
          style: {padding: 8, font: 'inherit', border: '1px solid #888', borderRadius: 6}, ...options}))
    }
    return h('section', {style: {maxWidth: 540}}, h('h2', null, 'Trim Seek'),
      h('p', null, '与本机 Codex 的 Seek 插件共享 NAS 连接。'),
      h('form', {onSubmit: submit, style: {display: 'grid', gap: 12}},
        field('NAS 地址', address, setAddress, {placeholder: 'http://192.168.9.118:5666/seek/'}),
        h('small', null, '仅填 IP 时默认使用 HTTP 和端口 5666。'),
        field('Seek 账号', username, setUsername, {autoComplete: 'username'}),
        h('small', null, '当前测试阶段暂用 Seek 账号和密码登录，后续将改用 API Key。'),
        field('密码', password, setPassword, {type: 'password', autoComplete: 'current-password'}),
        h('button', {type: 'submit', disabled: busy, style: {justifySelf: 'start', padding: '6px 18px'}}, '登录'),
        h('output', {'aria-live': 'polite'}, message)))
  }
  return {inject: ['slots', 'remote'], apply(ctx) {
    const mounted = ctx.remote.$mount({package: pkg, descriptors})
    mounted.catch(() => {}) // surface the error in the settings page
    ctx.effect(() => () => { mounted.then(dispose => dispose()).catch(() => {}) }, 'trim-seek remote')
    async function invoke(method, ...args) {
      await mounted
      const result = await ctx.get('remote.trimSeek')[method](...args)
      if (!result.ok) throw new Error(result.error?.message || '连接请求失败')
      return result.value
    }
    ctx.slots.inject('settings.section', () => ctx.slots.register({name: 'settings.section',
      id: 'trim-seek', order: 14, label: () => 'Trim Seek', inject: () => ({invoke})}, Settings))
  }}
}
if (typeof window !== 'undefined' && window.__ModuleLoader__) {
  window.__ModuleLoader__.load({id: pkg, factory: client})
} else { module.exports = {pkg, descriptors} }
