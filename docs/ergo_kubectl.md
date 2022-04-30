## ergo kubectl

Kubectl controls the Kubernetes cluster manager

### Synopsis

kubectl controls the Kubernetes cluster manager.

 Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/

```
ergo kubectl [flags]
```

### Options

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/Users/ysicing/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
  -h, --help                           help for kubectl
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --match-server-version           Require server version to match client version
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --profile string                 Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex) (default "none")
      --profile-output string          Name of the file to write the profile to (default "profile.pprof")
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
      --warnings-as-errors             Treat warnings received from the server as errors and exit with a non-zero exit code
```

### Options inherited from parent commands

```
      --config string   The ergo config file to use (default "/Users/ysicing/.ergo/config/ergo.yml")
      --debug           Prints the stack trace if an error occurs
      --silent          Run in silent mode and prevents any ergo log output except panics & fatals
```

### SEE ALSO

* [ergo](ergo.md)	 - ergo, ergo, NB!
* [ergo kubectl alpha](ergo_kubectl_alpha.md)	 - Commands for features in alpha
* [ergo kubectl annotate](ergo_kubectl_annotate.md)	 - 更新一个资源的注解
* [ergo kubectl api-resources](ergo_kubectl_api-resources.md)	 - Print the supported API resources on the server
* [ergo kubectl api-versions](ergo_kubectl_api-versions.md)	 - Print the supported API versions on the server, in the form of "group/version"
* [ergo kubectl apply](ergo_kubectl_apply.md)	 - Apply a configuration to a resource by file name or stdin
* [ergo kubectl attach](ergo_kubectl_attach.md)	 - Attach 到一个运行中的 container
* [ergo kubectl auth](ergo_kubectl_auth.md)	 - Inspect authorization
* [ergo kubectl autoscale](ergo_kubectl_autoscale.md)	 - Auto-scale a deployment, replica set, stateful set, or replication controller
* [ergo kubectl certificate](ergo_kubectl_certificate.md)	 - 修改 certificate 资源.
* [ergo kubectl cluster-info](ergo_kubectl_cluster-info.md)	 - Display cluster information
* [ergo kubectl completion](ergo_kubectl_completion.md)	 - Output shell completion code for the specified shell (bash, zsh or fish)
* [ergo kubectl config](ergo_kubectl_config.md)	 - 修改 kubeconfig 文件
* [ergo kubectl cordon](ergo_kubectl_cordon.md)	 - 标记 node 为 unschedulable
* [ergo kubectl cp](ergo_kubectl_cp.md)	 - Copy files and directories to and from containers
* [ergo kubectl create](ergo_kubectl_create.md)	 - Create a resource from a file or from stdin
* [ergo kubectl debug](ergo_kubectl_debug.md)	 - Create debugging sessions for troubleshooting workloads and nodes
* [ergo kubectl delete](ergo_kubectl_delete.md)	 - Delete resources by file names, stdin, resources and names, or by resources and label selector
* [ergo kubectl describe](ergo_kubectl_describe.md)	 - 显示一个指定 resource 或者 group 的 resources 详情
* [ergo kubectl diff](ergo_kubectl_diff.md)	 - Diff the live version against a would-be applied version
* [ergo kubectl drain](ergo_kubectl_drain.md)	 - Drain node in preparation for maintenance
* [ergo kubectl edit](ergo_kubectl_edit.md)	 - 在服务器上编辑一个资源
* [ergo kubectl exec](ergo_kubectl_exec.md)	 - 在一个 container 中执行一个命令
* [ergo kubectl explain](ergo_kubectl_explain.md)	 - Get documentation for a resource
* [ergo kubectl expose](ergo_kubectl_expose.md)	 - Take a replication controller, service, deployment or pod and expose it as a new Kubernetes service
* [ergo kubectl get](ergo_kubectl_get.md)	 - 显示一个或更多 resources
* [ergo kubectl kustomize](ergo_kubectl_kustomize.md)	 - Build a kustomization target from a directory or URL.
* [ergo kubectl label](ergo_kubectl_label.md)	 - 更新在这个资源上的 labels
* [ergo kubectl logs](ergo_kubectl_logs.md)	 - 输出容器在 pod 中的日志
* [ergo kubectl options](ergo_kubectl_options.md)	 - 输出所有命令的层级关系
* [ergo kubectl patch](ergo_kubectl_patch.md)	 - Update fields of a resource
* [ergo kubectl plugin](ergo_kubectl_plugin.md)	 - Provides utilities for interacting with plugins
* [ergo kubectl port-forward](ergo_kubectl_port-forward.md)	 - Forward one or more local ports to a pod
* [ergo kubectl proxy](ergo_kubectl_proxy.md)	 - 运行一个 proxy 到 Kubernetes API server
* [ergo kubectl replace](ergo_kubectl_replace.md)	 - Replace a resource by file name or stdin
* [ergo kubectl rollout](ergo_kubectl_rollout.md)	 - Manage the rollout of a resource
* [ergo kubectl run](ergo_kubectl_run.md)	 - 在集群中运行一个指定的镜像
* [ergo kubectl scale](ergo_kubectl_scale.md)	 - Set a new size for a deployment, replica set, or replication controller
* [ergo kubectl set](ergo_kubectl_set.md)	 - 为 objects 设置一个指定的特征
* [ergo kubectl taint](ergo_kubectl_taint.md)	 - 更新一个或者多个 node 上的 taints
* [ergo kubectl top](ergo_kubectl_top.md)	 - Display resource (CPU/memory) usage
* [ergo kubectl uncordon](ergo_kubectl_uncordon.md)	 - 标记 node 为 schedulable
* [ergo kubectl version](ergo_kubectl_version.md)	 - 输出 client 和 server 的版本信息
* [ergo kubectl wait](ergo_kubectl_wait.md)	 - Experimental: Wait for a specific condition on one or many resources

###### Auto generated by spf13/cobra on 30-Apr-2022
