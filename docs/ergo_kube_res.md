## ergo kube res

resource

### Options

```
      --context string      context to use for Kubernetes config
  -h, --help                help for res
      --kubeconfig string   kubeconfig file to use for Kubernetes config
  -o, --output string       prints the output in the specified format. Allowed values: table, json, yaml (default table)
  -s, --sortBy string       sort by cpu or memory (default "cpu")
```

### Options inherited from parent commands

```
      --config string   The ergo config file to use (default "/Users/ysicing/.ergo/config/ergo.yml")
      --debug           Prints the stack trace if an error occurs
      --silent          Run in silent mode and prevents any ergo log output except panics & fatals
```

### SEE ALSO

* [ergo kube](ergo_kube.md)	 - kube ops tools
* [ergo kube res node](ergo_kube_res_node.md)	 - node provides an overview of the node
* [ergo kube res pod](ergo_kube_res_pod.md)	 - pod provides an overview of the pod

###### Auto generated by spf13/cobra on 30-Apr-2022