# kube-doom

## Usage

### General

```
Kube-doom is a new way of chaos engineering. helps you to kill pods or namespaces through Doom.

Usage:
  kube-doom [command]

Available Commands:
  help        Help about any command
  run         A brief description of your command

Flags:
      --config string   config file (default is $HOME/.kube-doom.yaml)
  -h, --help            help for kube-doom
  -t, --toggle          Help message for toggle

Use "kube-doom [command] --help" for more information about a command.
```

### Run

```
Usage:
  kube-doom run [flags]

Flags:
  -h, --help               help for run
  -m, --mode string        Mode: pods or namespaces are allowed (default "pods")
  -n, --namespace string   In which namespace do you want to kill Pods
```

If `-namespace` is not filled, all Pods in all namespaces are retrieved.


TODO:
* add golangci
* add MAKEFILE
* add DockerFile
* test