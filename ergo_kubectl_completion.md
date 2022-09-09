## ergo kubectl completion

Output shell completion code for the specified shell (bash, zsh or fish)

### Synopsis

Output shell completion code for the specified shell (bash, zsh, fish, or powershell). The shell code must be evaluated to provide interactive completion of kubectl commands.  This can be done by sourcing it from the .bash_profile.

 Detailed instructions on how to do this are available here:

    for macOS:
    https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/#enable-shell-autocompletion
  
    for linux:
    https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#enable-shell-autocompletion
  
    for windows:
    https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/#enable-shell-autocompletion
  
 Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2.

```
ergo kubectl completion SHELL
```

### Examples

```
  # Installing bash completion on macOS using homebrew
  ## If running Bash 3.2 included with macOS
  brew install bash-completion
  ## or, if running Bash 4.1+
  brew install bash-completion@2
  ## If kubectl is installed via homebrew, this should start working immediately
  ## If you've installed via other means, you may need add the completion to your completion directory
  kubectl completion bash > $(brew --prefix)/etc/bash_completion.d/kubectl
  
  
  # Installing bash completion on Linux
  ## If bash-completion is not installed on Linux, install the 'bash-completion' package
  ## via your distribution's package manager.
  ## Load the kubectl completion code for bash into the current shell
  source <(kubectl completion bash)
  ## Write bash completion code to a file and source it from .bash_profile
  kubectl completion bash > ~/.kube/completion.bash.inc
  printf "
  # Kubectl shell completion
  source '$HOME/.kube/completion.bash.inc'
  " >> $HOME/.bash_profile
  source $HOME/.bash_profile
  
  # Load the kubectl completion code for zsh[1] into the current shell
  source <(kubectl completion zsh)
  # Set the kubectl completion code for zsh[1] to autoload on startup
  kubectl completion zsh > "${fpath[1]}/_kubectl"
  
  
  # Load the kubectl completion code for fish[2] into the current shell
  kubectl completion fish | source
  # To load completions for each session, execute once:
  kubectl completion fish > ~/.config/fish/completions/kubectl.fish
  
  # Load the kubectl completion code for powershell into the current shell
  kubectl completion powershell | Out-String | Invoke-Expression
  # Set kubectl completion code for powershell to run on startup
  ## Save completion code to a script and execute in the profile
  kubectl completion powershell > $HOME\.kube\completion.ps1
  Add-Content $PROFILE "$HOME\.kube\completion.ps1"
  ## Execute completion code in the profile
  Add-Content $PROFILE "if (Get-Command kubectl -ErrorAction SilentlyContinue) {
  kubectl completion powershell | Out-String | Invoke-Expression
  }"
  ## Add completion code directly to the $PROFILE script
  kubectl completion powershell >> $PROFILE
```

### Options

```
  -h, --help   help for completion
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
