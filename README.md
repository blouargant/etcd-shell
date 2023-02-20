# etcd-shell

ETCD Shell

### Tools
CLI is based on Cobra:  
https://github.com/spf13/cobra  

### Package builder
You need to install nfpm to generate linux packages:
https://nfpm.goreleaser.com/install/

```bash
echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
sudo apt update
sudo apt install nfpm
```