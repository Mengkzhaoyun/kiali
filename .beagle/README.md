# build

```bash
docker buildx build \
  --build-args VERSION=v1.22.1
  --tag registry.cn-qingdao.aliyuncs.com/wod/kiali-arm64:v1.22.1 \
  --platform linux/arm64 \
  --file .beagle/arm.dockerfile .
```
