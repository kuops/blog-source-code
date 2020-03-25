---
title: Fzf 快速入门
date: 2020-03-24 15:50:08
index_img: "/images/fzf_quick_start/fzf_index.png"
tags:
- fzf
categories:
- terminal-tools
---

<!-- markdownlint-disable MD013 -->
## 简介

fzf 是 go 语言写的一款命令行模糊查找器。

## 用法

使用 vim + fzf 快速查找当前目录下文件并打开

```bash
vim $(fzf)
```

使用方向键或者 `C-j`, `C-k` 或者 `C-n`, `C-p`选择项目,回车键确认项目

安装完毕之后, 默认绑定的快捷键:

- `C-t` 粘贴选择的目录或者文件到命令行

- `C-r` 打开历史命令记录，选择后粘贴到命令行

- `ALT-c` 进入选择的目录, 如果使用的是 iterm2 需要单独设定 `options` 按键为 `Esc+`

![iterm2_settings](/images/fzf_quick_start/fzf_iterm2_profiles_setting.png)

## 模糊补全

补全文件和目录:

- `COMMAND [DIRECTORY/][FUZZY_PATTERN]**<TAB>`

```bash
# Files under current directory
# - You can select multiple items with TAB key
vim **<TAB>

# Files under parent directory
vim ../**<TAB>

# Files under parent directory that match `fzf`
vim ../fzf**<TAB>

# Files under your home directory
vim ~/**<TAB>


# Directories under current directory (single-selection)
cd **<TAB>

# Directories under ~/github that match `fzf`
cd ~/github/fzf**<TAB>
```

补全进程 id

```bash
# Can select multiple processes with <TAB> or <Shift-TAB> keys
kill -9 <TAB>
```

补全 ssh 主机名(`~/.ssh/config`), 或者 telnet (`/etc/hosts`):

```bash
ssh **<TAB>
telnet **<TAB>
```

补全环境变量, 和别名

```bash
unset **<TAB>
export **<TAB>
unalias **<TAB>
```

设置 preview, 需要安装 `bat` 命令如果没有,可以用 `cat` 或者 `head` 代替, 添加到 `~/.zshrc`

```bash
FZF_DEFAULT_OPTS=
FZF_DEFAULT_OPTS+=" --preview-window 'right:60%'"
FZF_DEFAULT_OPTS+=" --bind alt-k:preview-up,alt-k:preview-up,alt-j:preview-down,alt-n:preview-down"
FZF_DEFAULT_OPTS+=" --preview 'bat --color=always --italic-text=always --style=numbers,changes,header --line-range :300 {} '"
export FZF_DEFAULT_OPTS
```

设置跟踪软连, 排除 git 目录

```bash
export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
```

## vim

安装以下两个插件:

```bash
Plug '/usr/local/opt/fzf'
Plug 'junegunn/fzf.vim'
```

常用命令

- `:Files` 查看当前目录的文件
- `:Buffers` 切换缓冲区
