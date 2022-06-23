go9p
====

Fork of the original `rminnich/go9p` repo with the patches of `k8s.io/minikube`.


This is go9p done in a way that I can understand.

To install:

```
$ go build -o ufs-server cmd/ufs/main.go
$ ./ufs-server
$ GOOS=windows GOARCH=amd64 go build -o ufs-server-windows.exe ./cmd/ufs/main.go
PS> .\ufs-server-windows.exe
```
