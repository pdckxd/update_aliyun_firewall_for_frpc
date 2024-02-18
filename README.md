## Introduce
Running on router to change SourceCidrIp for all matched security rules in aliyun

## How to build for AMD64 platform
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"

## How to run
```
export ALIBABA_CLOUD_ACCESS_KEY_ID=""
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=""
# location tag is suffix of description field in security role, for example: <location tag> should be 'company' if rule description is 'ssh - company'
./update_aliyun_firewall_for_frpc <location tag>
```
