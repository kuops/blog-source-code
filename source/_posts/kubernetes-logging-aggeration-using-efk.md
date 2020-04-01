---
title: 使用 efk 对 kubernetes 日志进程聚合
date: 2020-03-31 22:29:32
index_img: "/images/kubernetes-logging-aggeration-using-efk/efk_index.jpg"
tags:
- logging
categories:
- kubernetes
---

## 日志分类

在 kubernetes 中，日志大概可以分为三类:

### docker 标准输出

在 docker 的 `/var/lib/docker/containers/xxxx-json.log`，在 kubernets 中，`/var/log/containers/` 使用 `readlink` 指向 `/var/log/pods/<namespace>-<podname>/<containers>/x.log`，而 `/var/log/pods/` 又指向了 docker 的 json.log, 因为这种日志输出不构标准，每个程序的日志规范不太一样，对该日志只进行按行搜索，不进行分词操作。

### 程序的标准输出

这一类日志一般跟开发约定好规范,使用同一目录进行存放，如容器中的 `/data/logs/appname/*.log`, 而 `/data/logs` 可以作为 pod 的 volume 进行挂载, 因为该日志是一种标准规范，可以进行分词操作。

### ingress 访问日志

这一类日志，用来定位请求，一般用来做故障分析，包括域名，状态码，响应时间等信息。

##  日志搜集实践

### 准备阶段

1. 准备 kubernets 集群,并安装好网络插件。

```bash
# 步骤略
kubeadm init 
```

2. 准备 persistent volume ，这里为了方便使用 rancher 的 local-path-provisioner，默认的路径为 `/opt/local-path-provisioner`

```bash
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
```

3. 安装 istio ingress

```bash
istioctl manifest apply  --set values.global.proxy.accessLogFile="/dev/stdout" --set values.mixer.telemetry.enabled=false --set values.prometheus.enabled=false
```

3. 部署示例应用,这里暂且定义为 `/data/log/log-app/info.log` 镜像为 `kuops/log-example-app`, [源码链接](/files/kubernetes-logging-aggeration-using-efk/main.go.text)

```bash
kubectl apply -f 
```


### 部署 es



