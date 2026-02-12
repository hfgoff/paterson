#!/usr/bin/env bash
cd "$HOME" || exit
Ver="1.24.6"
FileName="go${Ver}.linux-armv6l.tar.gz"
wget https://dl.google.com/go/$FileName
sudo tar -C /usr/local -xvf $FileName
cat >> ~/.bashrc << 'EOF'
export GOPATH=$HOME/go
export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin
EOF
source "$HOME/.bashrc"