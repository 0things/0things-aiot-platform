# 生产集群证书组件

此目录管理集群级证书前置配置；`0things-test` 和 `0things-prod` 均不安装 cert-manager 控制器。命令在仓库根目录执行，要求 `kubectl` 已连接目标 k3s 集群。这里不会创建或修改 Cloudflare A 记录，也不会把 Token 保存到仓库。

## 1. 确认 Traefik

```bash
kubectl -n kube-system get pods
kubectl -n kube-system get svc traefik
kubectl get crd middlewares.traefik.io
```

如果系统 Pod 仍为 `ContainerCreating` 或 Traefik Service 不存在，先解决 k3s 镜像拉取问题，不继续安装应用入口。

## 2. 安装 cert-manager

`cert-manager.yaml` 是 k3s 原生 `HelmChart` 清单：k3s Helm Controller 根据它安装、升级并维护独立的 cert-manager Helm release。清单引用 Jetstack 官方 Helm 仓库中的 `v1.21.2` Chart（支持当前服务器报告的 Kubernetes `v1.36`），不把 cert-manager 放进 0things 应用 Chart：

```bash
kubectl apply -f deploy/helm/cluster/cert-manager.yaml
kubectl -n kube-system get helmchart cert-manager
kubectl -n kube-system wait --for=condition=Complete job/helm-install-cert-manager --timeout=10m
kubectl -n cert-manager rollout status deployment/cert-manager --timeout=10m
kubectl -n cert-manager rollout status deployment/cert-manager-webhook --timeout=10m
kubectl -n cert-manager rollout status deployment/cert-manager-cainjector --timeout=10m
kubectl wait --for=condition=Established crd/clusterissuers.cert-manager.io --timeout=120s
kubectl -n cert-manager get pods
```

升级时修改清单中的 `spec.version` 并重新 `kubectl apply`。安装失败时查看 `kubectl -n kube-system get jobs` 和 `kubectl -n kube-system describe helmchart cert-manager`。升级版本前核对 [cert-manager 支持的 Kubernetes 版本](https://cert-manager.io/docs/releases/)、[官方 Helm Chart](https://cert-manager.io/docs/installation/helm/) 与 [k3s Helm Controller 用法](https://docs.k3s.io/add-ons/helm)。

## 3. 配置 Cloudflare DNS-01

Cloudflare Token 只授权 `0thing.com` 区域的 `Zone - DNS - Edit` 和 `Zone - Zone - Read`。在本地、未被 Git 跟踪的 `deploy/helm/0things/values-prod.yaml` 中填写 `clusterIssuer.apiToken`。部署 prod 0things Chart 时，Helm 会在 `cert-manager` namespace 自动创建 `cloudflare-api-token` Secret 和 `ClusterIssuer`；test 不创建。不要提交包含 Token 的 values 文件，也不要将 `helm template` 或 `helm get values` 的完整输出公开；Helm release 记录中也会保存这个 Token。DNS-01 只在签发和续期时维护临时 TXT 验证记录。Cloudflare Token 权限说明见 [cert-manager Cloudflare 文档](https://cert-manager.io/docs/configuration/acme/dns01/cloudflare/)。

## 4. 部署应用和验证证书

确认前端生产镜像、生产 values 和数据库结构已准备好后，按 [0things Chart 部署说明](../0things/README.md) 安装 `0things-prod`。随后检查：

```bash
kubectl -n 0things-prod get certificate,ingress
kubectl -n cert-manager get secret cloudflare-api-token
kubectl wait --for=condition=Ready clusterissuer/letsencrypt-cloudflare --timeout=120s
kubectl -n 0things-prod wait --for=condition=Ready certificate/0thing-com --timeout=10m
kubectl -n 0things-prod get secret aiot-platform-tls
```

若证书未就绪，查看 `kubectl -n 0things-prod describe certificate 0thing-com` 和同 namespace 的 `order,challenge`。公网 A 记录及 HTTPS 路由需要单独确认；本配置不会更改它们。
