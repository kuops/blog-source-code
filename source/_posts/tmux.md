---
title: tmux 快速入门
date: 2020-03-22 13:34:52
tags:
- tmux
categories:
- terminal-tools
---

## 简介 

tmux 的作者将其描述为终端多路复用器 (terminal multiplexer)。使用 tmux 的好处主要有两点:

- 多窗口操作，在一个终端分出多个面板
- 避免 SSH 连接不稳定，断开前台任务问题

tmux 所有命令都以前置按键方式进行触发（默认为 `C-b`），`C-b` 表示按住 `ctrl` 键和 ``b` 键:


## 常用操作


常用命令:

| 描述 | 操作 |
|:---|:---|
|启动新的会话| 终端输入 `tmux` |
|创建新的会话并指定名称|终端输入 `tmux new -s name`|
|查看会话列表|终端输入 `tmux ls`|
|重新连接会话|终端输入 `tmux a`|
|指定会话连接|终端输入 `tmux a -t num/name`|
|启动新的窗口|终端输入 `tmux neww -n name`|
|根据会话名称启动新的窗口|终端输入 `tmux new -s session-name -n name`|
|退出会话|终端输入 `tmux detach`|
|结束会话|终端输入 `tmux kill-session -t num/name`|


> 这里退出的意思是退出窗口使其后台运行，结束为 kill

如果在一个会话中可以使用命令模式, 按住前置键 + `:`, 进入命令模式:

| 描述 | 操作 |
|:---|:---|
|新建窗口|命令模式输入 `new-window -n name`|
|新建会话|命令模式输入 `new -s name`|
|切换会话|命令模式输入 `attach-session -t num/name`|


常用快捷键

| 描述 | 操作 |
|:---|:---|
|显示快捷键帮助|`prefix ?`|
|重命名当前会话|`prefix $`|
|切换到上一个会话|`prefix (`|
|切换到下一个会话|`prefix )`|
|退出会话，使其在后台运行|`prefix d`|
|切换窗口|`prefix [0-9]`|
|创建新的窗口|`prefix c`|
|重命名当前窗口|`prefix ,`|
|显示所有窗口的可选择列表|`prefix w`|
|结束窗口|`prefix  &`|
|水平分割面板|`prefix "`|
|垂直分割面板|`prefix %`|
|切换面板|`prefix 方向键`|
|显示面板编号|`prefix q`|
|关闭面板|`prefix x`|
|切换到下一个面板|`prefix o`|
|交换面板位置|`prefix }`|


## 复制模式

添加下面一行到 $HOME/.tmux.conf, 通过 vim 的快捷键实现浏览, 复制等操作;

```bash
setw -g mode-keys vi
```

| 描述 | 操作 |
|:---|:---|
|进入复制模式|`prefix [`|
|粘贴选择内容(buffer_0)|`prefix ]`|
|显示 buffer_0 的内容|命令行模式输入`show-buffer`|
|复制整个能见的内容到当前的 buffer|`capture-buffer`|
|列出所有的 buffer|`list-buffers`|
|选择用于粘贴的 buffer|`choose-buffer`|
|将 buffer 的内容复制到文件|save-buffer file.txt|

| vi     | emacs     | 功能                 |
| ------ | --------- | ---                  |
| ^      | M-m       | 跳转到一行开头       |
| Escape | C-g       | 放弃选择             |
| k      | Up        | 上移                 |
| j      | Down      | 下移                 |
| h      | Left      | 左移                 |
| l      | Right     | 右移                 |
| L      |           | 最后一行             |
| M      | M-r       | 中间一行             |
| H      | M-R       | 第一行               |
| $      | C-e       | 跳转到行尾           |
| :      | g         | 跳转至某一行         |
| C-d    | M-Down    | 下翻半页             |
| C-u    | M-Up      | 上翻半页             |
| C-f    | Page down | 下翻一页             |
| C-b    | Page up   | 上翻一页             |
| w      | M-f       | 下一个字符           |
| b      | M-b       | 前一个字符           |
| q      | Escape    | 退出                 |
| ?      | C-r       | 往上查找             |
| /      | C-s       | 往下查找             |
| n      | n         | 查找下一个           |
| Space  | C-Space   | 进入选择模式         |
| Enter  | M-w       | 确认选择内容, 并退出 |

## 自定义配置

自定义配置放在 ~/.tmux.conf 中，可以自行 github 搜索相关的配置。

正在使用的配置:

主题: [https://github.com/jimeh/tmux-themepack](https://github.com/jimeh/tmux-themepack)
快捷键映射: [https://github.com/square/maximum-awesome/](https://github.com/square/maximum-awesome)
