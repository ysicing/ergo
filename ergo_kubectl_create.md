## ergo kubectl create

Create a resource from a file or from stdin

### Synopsis

Create a resource from a file or from stdin.

 JSON and YAML formats are accepted.

```
ergo kubectl create -f FILENAME
```

### Examples

```
  # Create a pod using the data in pod.json
  kubectl create -f ./pod.json
  
  # Create a pod based on the JSON passed into stdin
  cat pod.json | kubectl create -f -
  
  # Edit the data in docker-registry.yaml in JSON then create the resource using the edited data
  kubectl create -f docker-registry.yaml --edit -o json
```

### Options

```
      --allow-missing-template-keys    If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      --dry-run string[="unchanged"]   Must be "none", "server", or "client". If client strategy, only print the object that would be sent, without sending it. If server strategy, submit server-side request without persisting the resource. (default "none")
      --edit                           Edit the API resource before creating
      --field-manager string           Name of the manager used to track field ownership. (default "kubectl-create")
  -f, --filename strings               Filename, directory, or URL to files to use to create the resource
  -h, --help                           help for create
  -k, --kustomize string               Process the kustomization directory. This flag can't be used together with -f or -R.
  -o, --output string                  Output format. One of: json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-as-json|jsonpath-file.
      --raw string                     Raw URI to POST to the server.  Uses the transport specified by the kubeconfig file.
  -R, --recursive                      Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
      --save-config                    If true, the configuration of current object will be saved in its annotation. Otherwise, the annotation will be unchanged. This flag is useful when you want to perform kubectl apply on this object in the future.
  -l, --selector string                Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
      --show-managed-fields            If true, keep the managedFields when printing objects in JSON or YAML format.
      --template string                Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --validate                       If true, use a schema to validate the input before sending it (default true)
      --windows-line-endings           Only relevant if --edit=true. Defaults to the line ending native to your platform.
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
* [ergo kubectl create clusterrole](ergo_kubectl_create_clusterrole.md)	 - Create a cluster role
* [ergo kubectl create clusterrolebinding](ergo_kubectl_create_clusterrolebinding.md)	 - Create a cluster role binding for a particular cluster role
* [ergo kubectl create configmap](ergo_kubectl_create_configmap.md)	 - Create a config map from a local file, directory or literal value
* [ergo kubectl create cronjob](ergo_kubectl_create_cronjob.md)	 - Create a cron job with the specified name
* [ergo kubectl create deployment](ergo_kubectl_create_deployment.md)	 - Create a deployment with the specified name
* [ergo kubectl create ingress](ergo_kubectl_create_ingress.md)	 - Create an ingress with the specified name
* [ergo kubectl create job](ergo_kubectl_create_job.md)	 - Create a job with the specified name
* [ergo kubectl create namespace](ergo_kubectl_create_namespace.md)	 - Create a namespace with the specified name
* [ergo kubectl create poddisruptionbudget](ergo_kubectl_create_poddisruptionbudget.md)	 - Create a pod disruption budget with the specified name
* [ergo kubectl create priorityclass](ergo_kubectl_create_priorityclass.md)	 - Create a priority class with the specified name
* [ergo kubectl create quota](ergo_kubectl_create_quota.md)	 - Create a quota with the specified name
* [ergo kubectl create role](ergo_kubectl_create_role.md)	 - Create a role with single rule
* [ergo kubectl create rolebinding](ergo_kubectl_create_rolebinding.md)	 - Create a role binding for a particular role or cluster role
* [ergo kubectl create secret](ergo_kubectl_create_secret.md)	 - Create a secret using specified subcommand
* [ergo kubectl create service](ergo_kubectl_create_service.md)	 - Create a service using a specified subcommand
* [ergo kubectl create serviceaccount](ergo_kubectl_create_serviceaccount.md)	 - Create a service account with the specified name

###### Auto generated by spf13/cobra on 14-May-2022
