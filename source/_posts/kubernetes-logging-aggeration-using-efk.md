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

### 程序的规范日志

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

3. 部署示例应用,使用镜像 `kuops/log-example-app`,程序会同时往 `/data/log/log-app/info.log`和标准输出来输出日志, [程序源码](https://raw.githubusercontent.com/kuops/kuops.github.io/master/files/kubernetes-logging-aggeration-using-efk/main.go),

```bash
kubectl apply -f https://raw.githubusercontent.com/kuops/kuops.github.io/master/files/kubernetes-logging-aggeration-using-efk/log-app.yaml
```

### 部署 efk

首先对 yaml 进行修改: 

```bash
# 下载 yaml
curl -LO https://raw.githubusercontent.com/kuops/kuops.github.io/master/files/kubernetes-logging-aggeration-using-efk/elasticsearch.yaml
curl -LO https://raw.githubusercontent.com/kuops/kuops.github.io/master/files/kubernetes-logging-aggeration-using-efk/filebeat.yaml
curl -LO https://raw.githubusercontent.com/kuops/kuops.github.io/master/files/kubernetes-logging-aggeration-using-efk/kibana.yaml

#############################
# elasticsearch.yaml 可以修改 
#############################
# service kind 里的 spec.ports.nodeport
...
  ports:
  - name: http
    port: 9200
    # nodeport ，固定为 31001
    nodePort: 31001
    protocol: TCP
  type: NodePort
###########################
# filbeat.yaml 修改以下地方
###########################
# daemonset kind 里的 spec.template.spec 中的 containers.volumeMounts 和 volumes
...
          volumeMounts:
	  ...
	  # 这里，如果没有单独为 docker 使用磁盘，修改为 `/var/lib/docker/containers`
          - name: varlibdockercontainers
            mountPath: /data/docker/containers
            readOnly: true
	  ...
	  # 这里, 是标准输出日志的 path , 路径是 rancher local volume 的 configmap 中的路径
          - name: filelog
            mountPath: /opt/local-path-provisioner
            readOnly: true
	...
        volumes:
	# 这里，如果没有单独为 docker 使用磁盘，修改为 `/var/lib/docker/containers`
        - name: varlibdockercontainers
          hostPath:
            path: /data/docker/containers
	...
	# 这里, 是标准输出日志的 path , 路径是 rancher local volume 的 configmap 中的路径
        - name: filelog
          hostPath:
            path: /opt/local-path-provisioner
###########################
# kibana.yaml 修改以下地方
###########################
# service kind 里的 spec.ports.nodeport
...
  ports:
    - port: 5601
      nodePort: 31002
      protocol: TCP
      name: http
      targetPort: 5601
  selector:
    k8s-app: kibana
  type: NodePort
```

部署 elasticsearch, filebeat, kibana;

```
kubectl apply -f elasticsearch.yaml
kubectl apply -f filebeat.yaml
kubectl apply -f kibana.yaml
```

修改 istio 的 ingressgateway

```
# 使用 kubectl edit 修改 添加 hostPort 
kubectl edit -n istio-system  deployments.apps istio-ingressgateway
...
        ports:
        - containerPort: 15020
          protocol: TCP
        - containerPort: 80
          protocol: TCP
	  # 添加 hostPort
          hostPort: 80

```

添加 virtualservice

```bash
curl -LO https://raw.githubusercontent.com/kuops/kuops.github.io/master/files/kubernetes-logging-aggeration-using-efk/kibana-vs.yaml

# 修改 INGRESS_NODE_IP 为 ingressgateway 所在 node 的 IP, 例如 `kibana.10.7.0.101.nip.io`
...
hosts:
  - "kibana.INGRESS_NODE_IP.nip.io"
```

在 kibana 中, 使用 devtools 为 es 中添加 ingress Pipeline

```bash
PUT _ingest/pipeline/ingress
{
    "description": "ingress",
    "on_failure": [
      {
        "set": {
          "field": "_index",
          "value": "failed-{{ _index }}"
        }
      }
    ],
    "processors": [
      {
        "grok": {
          "field": "message",
          "patterns": [
            """\[%{TIMESTAMP_ISO8601:timestamp}\] "%{DATA:method} (?:%{URIPATH:uri_path}(?:%{URIPARAM:uri_param})?|%{DATA:})%{DATA:protocol}" %{NUMBER:status_code:int} %{DATA:response_flags} %{NUMBER:bytes_sent:int} %{NUMBER:bytes_received:int} %{NUMBER:duration:int} (%{NUMBER:upstream_service_time}|-) "%{IPORHOST:remote_ip}(?:,\s)?%{DATA:forwarded_for}" "%{DATA:user_agent}" "%{DATA:request_id}" "%{DATA:authority}" %{DATA:upstream_service}"""
          ]
        }
      },
      {
        "gsub": {
          "field": "uri_path",
          "pattern": "HTTP/1.1",
          "replacement": ""
        }
      },
      {
        "gsub": {
          "field": "uri_path",
          "pattern": "HTTP/2",
          "replacement": ""
        }
      },
      {
        "remove": {
          "field": "message"
        }
      }
    ]
}
```

在 es 中添加程序规范日志的 pipeline

```bash
PUT _ingest/pipeline/podfile
{
      "description": "podfile",
      "processors": [
        {
          "grok": {
            "field": "message",
            "patterns": [
              """%{TIMESTAMP_ISO8601:log_time} \[%{LOGLEVEL:log_level}\] %{GREEDYDATA:log_msg}"""
            ]
          }
        },
        {
          "remove": {
            "field": "message"
          }
        }
      ],
      "on_failure": [
        {
          "set": {
            "field": "_index",
            "value": "failed-{{ _index }}"
          }
        }
      ]
  }
```

查看程序规范输出日志

![](/images/kubernetes-logging-aggeration-using-efk/efk_podfile.jpg)


查看 ingress 日志

![](/images/kubernetes-logging-aggeration-using-efk/efk_ingress.jpg)

查看标准输出

![](/images/kubernetes-logging-aggeration-using-efk/efk_stdout.jpg)

