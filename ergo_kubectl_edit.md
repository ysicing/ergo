## ergo kubectl edit

Edit a resource on the server

### Synopsis

Edit a resource from the default editor.

 The edit command allows you to directly edit any API resource you can retrieve via the command-line tools. It will open the editor defined by your KUBE_EDITOR, or EDITOR environment variables, or fall back to 'vi' for Linux or 'notepad' for Windows. You can edit multiple objects, although changes are applied one at a time. The command accepts file names as well as command-line arguments, although the files you point to must be previously saved versions of resources.

 Editing is done with the API version used to fetch the resource. To edit using a specific API version, fully-qualify the resource, version, and group.

 The default format is YAML. To edit in JSON, specify "-o json".

 The flag --windows-line-endings can be used to force Windows line endings, otherwise the default for your operating system will be used.

 In the event an error occurs while updating, a temporary file will be created on disk that contains your unapplied changes. The most common error when updating a resource is another editor changing the resource on the server. When this occurs, you will have to apply your changes to the newer version of the resource, or update your temporary saved copy to include the latest resource version.

```
ergo kubectl edit (RESOURCE/NAME | -f FILENAME)
```

### Examples

```
  # Edit the service named 'registry'
  kubectl edit svc/registry
  
  # Use an alternative editor
  KUBE_EDITOR="nano" kubectl edit svc/registry
  
  # Edit the job 'myjob' in JSON using the v1 API format
  kubectl edit job.v1.batch/myjob -o json
  
  # Edit the deployment 'mydeployment' in YAML and save the modified config in its annotation
  kubectl edit deployment/mydeployment -o yaml --save-config
  
  # Edit the deployment/mydeployment's status subresource
  kubectl edit deployment mydeployment --subresource='status'
```

### Options

```
      --allow-missing-template-keys   If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      --field-manager string          Name of the manager used to track field ownership. (default "kubectl-edit")
  -f, --filename strings              Filename, directory, or URL to files to use to edit the resource
  -h, --help                          help for edit
  -k, --kustomize string              Process the kustomization directory. This flag can't be used together with -f or -R.
  -o, --output string                 Output format. One of: (json, yaml, name, go-template, go-template-file, template, templatefile, jsonpath, jsonpath-as-json, jsonpath-file).
      --output-patch                  Output the patch if the resource is edited.
  -R, --recursive                     Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
      --save-config                   If true, the configuration of current object will be saved in its annotation. Otherwise, the annotation will be unchanged. This flag is useful when you want to perform kubectl apply on this object in the future.
      --show-managed-fields           If true, keep the managedFields when printing objects in JSON or YAML format.
      --subresource string            If specified, edit will operate on the subresource of the requested object. Must be one of [status]. This flag is alpha and may change in the future.
      --template string               Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --validate string               Must be one of: strict (or true), warn, ignore (or false).
                                      		"true" or "strict" will use a schema to validate the input and fail the request if invalid. It will perform server side validation if ServerSideFieldValidation is enabled on the api-server, but will fall back to less reliable client-side validation if not.
                                      		"warn" will warn about unknown or duplicate fields without blocking the request if server-side field validation is enabled on the API server, and behave as "ignore" otherwise.
                                      		"false" or "ignore" will not perform any schema validation, silently dropping any unknown or duplicate fields. (default "strict")
      --windows-line-endings          Defaults to the line ending native to your platform.
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
