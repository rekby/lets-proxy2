#!/bin/bash
set -ev

cd output
ls

./lets-proxy_linux_amd64 --print-default-config > config_default.toml_


cp ../README.md ./README.md
unix2dos -n README.md README.txt

for FILE in `ls`; do
    expr match "$FILE" '^lets-proxy_' > /dev/null || continue
    echo "$FILE"
    if expr match "$FILE" '^lets-proxy_windows' > /dev/null; then
        mv "$FILE" lets-proxy.exe
        unix2dos -n config_default.toml_ config_default.toml
        zip "${FILE%.exe}.zip" lets-proxy.exe README.txt config_default.toml
        rm config_default.toml
        rm lets-proxy.exe
    else
        mv "$FILE" lets-proxy
        cp config_default.toml_ config_default.toml
        tar -zcvf "${FILE%.exe}.tar.gz" lets-proxy README.md config_default.toml
        rm config_default.toml
        rm lets-proxy
    fi
done

rm config_default.toml_ README.md README.txt
ls