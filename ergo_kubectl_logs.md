## ergo kubectl logs

Print the logs for a container in a pod

### Synopsis

Print the logs for a container in a pod or specified resource. If the pod has only one container, the container name is optional.

```
ergo kubectl logs [-f] [-p] (POD | TYPE/NAME) [-c CONTAINER]
```

### Examples

```
  # Return snapshot logs from pod nginx with only one container
  kubectl logs nginx
  
  # Return snapshot logs from pod nginx with multi containers
  kubectl logs nginx --all-containers=true
  
  # Return snapshot logs from all containers in pods defined by label app=nginx
  kubectl logs -l app=nginx --all-containers=true
  
  # Return snapshot of previous terminated ruby container logs from pod web-1
  kubectl logs -p -c ruby web-1
  
  # Begin streaming the logs of the ruby container in pod web-1
  kubectl logs -f -c ruby web-1
  
  # Begin streaming the logs from all containers in pods defined by label app=nginx
  kubectl logs -f -l app=nginx --all-containers=true
  
  # Display only the most recent 20 lines of output in pod nginx
  kubectl logs --tail=20 nginx
  
  # Show all logs from pod nginx written in the last hour
  kubectl logs --since=1h nginx
  
  # Show logs from a kubelet with an expired serving certificate
  kubectl logs --insecure-skip-tls-verify-backend nginx
  
  # Return snapshot logs from first container of a job named hello
  kubectl logs job/hello
  
  # Return snapshot logs from container nginx-1 of a deployment named nginx
  kubectl logs deployment/nginx -c nginx-1
```

### Options

```
      --all-containers                     Get all containers' logs in the pod(s).
  -c, --container string                   Print the logs of this container
  -f, --follow                             Specify if the logs should be streamed.
  -h, --help                               help for logs
      --ignore-errors                      If watching / following pod logs, allow for any errors that occur to be non-fatal
      --insecure-skip-tls-verify-backend   Skip verifying the identity of the kubelet that logs are requested from.  In theory, an attacker could provide invalid log content back. You might want to use this if your kubelet serving certificates have expired.
      --limit-bytes int                    Maximum bytes of logs to return. Defaults to no limit.
      --max-log-requests int               Specify maximum number of concurrent logs to follow when using by a selector. Defaults to 5. (default 5)
      --pod-running-timeout duration       The length of time (like 5s, 2m, or 3h, higher than zero) to wait until at least one pod is running (default 20s)
      --prefix                             Prefix each log line with the log source (pod name and container name)
  -p, --previous                           If true, print the logs for the previous instance of the container in a pod if it exists.
  -l, --selector string                    Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.
      --since duration                     Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used.
      --since-time string                  Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used.
      --tail int                           Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided. (default -1)
      --timestamps                         Include timestamps on each line in the log output
```

### Options inherited from parent commands

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/runner/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --config string                  The ergo config file to use (default "/home/runner/.ergo/config/ergo.yml")
      --context string                 The name of the kubeconfig context to use
      --debug                          Prints the stack trace if an error occurs
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --match-server-version           Require server version to match client version
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --profile string                 Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex) (default "none")
      --profile-output string          Name of the file to write the profile to (default "profile.pprof")
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --silent                         Run in silent mode and prevents any ergo log output except panics & fatals
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
      --warnings-as-errors             Treat warnings received from the server as errors and exit with a non-zero exit code
```

### SEE ALSO

* [ergo kubectl](ergo_kubectl.md)	 - Kubectl controls the Kubernetes cluster manager

###### Auto generated by spf13/cobra on 10-Sep-2022
