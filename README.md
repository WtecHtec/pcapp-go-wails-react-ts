# 自定义渲染、主进程通信函数
命令：
```
wails generate module

```
注：函数名必须大写字母开头

# 渲染、主进程通信
主进程：
```
  // WailApp.ctx 必须是app 的content
	runtime.EventsEmit(WailApp.ctx, "login:code", uuid)

```
渲染：
```
  EventsOn('login:code', async (unid) => {
    console.log(unid)
  })

```