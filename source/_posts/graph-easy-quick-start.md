---
title: graph-easy 快速入门
date: 2020-04-30 13:34:00
index_img: "/images/graph-easy-quick-start/graph-easy_index.jpg"
tags:
- graph-easy
categories:
- terminal-tools
---

<!-- markdownlint-disable MD013 -->
## 简介

一款命令行绘图工具

## 安装

```bash
brew install cpanminus
brew install graphviz
cpanm --mirror http://mirrors.163.com/cpan --mirror-only Graph::Easy
ln -s /usr/local/Cellar/perl/5.30.2_1/bin/graph-easy /usr/local/bin/graph-easy
```

## 基本用法

基本用法

```bash
echo "[step1] --> [step2]" |graph-easy
```

输出结果

```bash
+-------+     +-------+
| step1 | --> | step2 |
+-------+     +-------+
```

多个分支

```bash
echo "[step1] --> [step2-1],[step2-2],[step2-3]" |graph-easy
```

输出结果

```bash
+---------+     +---------+     +---------+
| step2-3 | <-- |  step1  | --> | step2-1 |
+---------+     +---------+     +---------+
                  |
                  |
                  v
                +---------+
                | step2-2 |
                +---------+
```

带注释的用法

```bash
echo "[step1] -- comment --> [step2] -- comment --> [step3]" |graph-easy
```

输出结果:

```bash

+-------+  comment   +-------+  comment   +-------+
| step1 | ---------> | step2 | ---------> | step3 |
+-------+            +-------+            +-------+
```

稍微复杂的可以以文件方式输出

```bash
# cat a.txt
[kube-node] -- request --> [127.0.0.1:8443]
[127.0.0.1:8443] -- upstream --> [kube-master1:6443]{origin: 127.0.0.1:8443; offset: 2,0}
[127.0.0.1:8443] -- upstream --> [kube-master2:6443]{origin: 127.0.0.1:8443; offset: 2,-2}
[127.0.0.1:8443] -- upstream --> [kube-master3:6443]{origin: 127.0.0.1:8443; offset: 2,2}

graph-easy a.txt
```

输出结果:

```bash
                                            upstream    +-------------------+
                           +--------------------------> | kube-master2:6443 |
                           |                            +-------------------+
                           |
                           |
                           |
+-----------+  request   +----------------+  upstream   +-------------------+
| kube-node | ---------> | 127.0.0.1:8443 | ----------> | kube-master1:6443 |
+-----------+            +----------------+             +-------------------+
                           |
                           |
                           |
                           |                upstream    +-------------------+
                           +--------------------------> | kube-master3:6443 |
                                                        +-------------------+
```

## 高级用法

参考 graph-easy-cn

```bash
https://weishu.gitbooks.io/graph-easy-cn/content/index.html
```

