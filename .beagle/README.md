# git

```bash
git remote add upstream git@github.com:kiali/kiali.git

git fetch upstream

git merge v1.28.1
```

## build

```bash
rm -rf console
docker run --rm \
-v $PWD/:/go/src/github.com/kiali/kiali \
-w /go/src/github.com/kiali/kiali \
registry.cn-qingdao.aliyuncs.com/wod/kiali-ui:v1.28.0 \
sh -c "mkdir -p console && cp -r /www/* console/"

export GOARCH=arm64
make build
```
