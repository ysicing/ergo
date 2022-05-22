## ergo kubectl attach

Attach to a running container

### Synopsis

Attach to a process that is already running inside an existing container.

```
ergo kubectl attach (POD | TYPE/NAME) -c CONTAINER
```

### Examples

```
  # Get output from running pod mypod; use the 'kubectl.kubernetes.io/default-container' annotation
  # for selecting the container to be attached or the first container in the pod will be chosen
  kubectl attach mypod
  
  # Get output from ruby-container from pod mypod
  kubectl attach mypod -c ruby-container
  
  # Switch to raw terminal mode; sends stdin to 'bash' in ruby-container from pod mypod
  # and sends stdout/stderr from 'bash' back to the client
  kubectl attach mypod -c ruby-container -i -t
  
  # Get output from the first pod of a replica set named nginx
  kubectl attach rs/nginx
```

### Options

```
  -c, --container string               Container name. If omitted, use the kubectl.kubernetes.io/default-container annotation for selecting the container to be attached or the first container in the pod will be chosen
  -h, --help                           help for attach
      --pod-running-timeout duration   The length of time (like 5s, 2m, or 3h, higher than zero) to wait until at least one pod is running (default 1m0s)
  -q, --quiet                          Only print output from the remote session
  -i, --stdin                          Pass stdin to the container
  -t, --tty                            Stdin is a TTY
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

###### Auto generated by spf13/cobra on 22-May-2022