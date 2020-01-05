## Dependencies

```
sudo apt-get update && sudo apt-get install -y \
    curl \
    gcc \
    gdb \
    nasm \
    strace \
    time

sudo snap install task

sudo bash -c 'curl -L https://golang.org/dl/go1.13.5.linux-amd64.tar.gz | tar -C /usr/local -xz'
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/bin
```
