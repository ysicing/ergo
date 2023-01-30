## ergo kubectl kustomize

Build a kustomization target from a directory or URL.

### Synopsis

Build a set of KRM resources using a 'kustomization.yaml' file. The DIR argument must be a path to a directory containing 'kustomization.yaml', or a git repository URL with a path suffix specifying same with respect to the repository root. If DIR is omitted, '.' is assumed.

```
ergo kubectl kustomize DIR [flags]
```

### Examples

```
  # Build the current working directory
  kubectl kustomize
  
  # Build some shared configuration directory
  kubectl kustomize /home/config/production
  
  # Build from github
  kubectl kustomize https://github.com/kubernetes-sigs/kustomize.git/examples/helloWorld?ref=v1.0.6
```

### Options

```
      --as-current-user          use the uid and gid of the command executor to run the function in the container
      --enable-alpha-plugins     enable kustomize plugins
      --enable-helm              Enable use of the Helm chart inflator generator.
  -e, --env stringArray          a list of environment variables to be used by functions
      --helm-command string      helm command (path to executable) (default "helm")
  -h, --help                     help for kustomize
      --load-restrictor string   if set to 'LoadRestrictionsNone', local kustomizations may load files from outside their root. This does, however, break the relocatability of the kustomization. (default "LoadRestrictionsRootOnly")
      --mount stringArray        a list of storage options read from the filesystem
      --network                  enable network access for functions that declare it
      --network-name string      the docker network to run the container in (default "bridge")
  -o, --output string            If specified, write output to this path.
      --reorder string           Reorder the resources just before output. Use 'legacy' to apply a legacy reordering (Namespaces first, Webhooks last, etc). Use 'none' to suppress a final reordering. (default "legacy")
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
      --disable-compression            If true, opt-out of response compression for all requests to the server
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

###### Auto generated by spf13/cobra on 30-Jan-2023
