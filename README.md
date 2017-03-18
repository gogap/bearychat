# Bearychat API


### Outgoing


下载代码

```
go get github.com/gogap/bearychat
```


#### Outgoing HTTP 

**直接启动现成的HTTP服务示例:**

一. 编译

```bash
cd $GOPATH/github.com/gogap/bearychat/outgoing/cmd/outgoing
go build
```

二. 配置

```bash
cp outgoing.conf.example outgoing.conf
```

样例配置文件：`outgoing.conf`

```hocon
{
    http.address=":3000"
    http.path="/triggers"

    outgoing {

        hello {
            word = "!hello"
            drivers = [gogap-auth, gogap-greeter]
            gogap-greeter = {
                name = "Robot GoGap"
                image="https://avatars2.githubusercontent.com/u/8731757?v=3&s=200"
            }

            gogap-auth = {
                token = "a5abc87ce0dcd169d4560385c2be9d3a"
            }
        }

        cmd {
            word = "!cmd"
            drivers = [gogap-auth, gogap-commands]

            gogap-auth = {
                token = "8831067e28290392313ca4d81356abe3"
            }

            gogap-commands = {
               timeout = 5s
               commands = {
                    ping = {
                        cmd = ping
                        cwd = /
                    }

                    ls = {
                        cmd = ls
                        cwd = /Users/zhengxujin/Downloads
                    }
                }
            }
        }
    }
}
```

`outgoing`这个程序默认会加载本地的 `outgoing.conf`, 如果您想指定某个配置文件，也可以使用`--config` 参数.

```bash
./outgoing run --config your-config.conf
```

config 文件采用的是 `hocon` 格式，同时兼容`JSON`，具体使用方法请参考：`https://github.com/go-akka/configuration`

三. 启动

```bash
./outgoing run
```

根据上面的配置信息，我们的服务监听地址为：`http://127.0.0.1:3000/triggers`

访问测试

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "token" : "8831067e28290392313ca4d81356abe3",
  "ts" : 1355517523,
  "text" : "!cmd ping -c 2 baidu.com",
  "trigger_word" : "!cmd",
  "subdomain" : "your_domain",
  "channel_name" : "your_channel",
  "user_name" : "Zeal"
}' "http://127.0.0.1:3000/triggers"
```

返回如下结果

```
{
    "text": "PING baidu.com (220.181.57.217): 56 data bytes\n64 bytes from 220.181.57.217: icmp_seq=0 ttl=51 time=4.935 ms\n64 bytes from 220.181.57.217: icmp_seq=1 ttl=51 time=6.317 ms\n\n--- baidu.com ping statistics ---\n2 packets transmitted, 2 packets received, 0.0% packet loss\nround-trip min/avg/max/stddev = 4.935/5.626/6.317/0.691 ms\n",
    "attachments": null
}
```

#### 自定义 Trigger

`Auth` Trigger样例

```go
package auth

import (
    "errors"

    "github.com/go-akka/configuration"
    "github.com/gogap/bearychat/outgoing"
)

type Auth struct {
    word  string
    token string
}

func init() {
    outgoing.RegisterTriggerDriver("gogap-auth", NewAuth)
}

func NewAuth(word string, config *configuration.Config) (outgoing.Trigger, error) {
    return &Auth{
        word:  word,
        token: config.GetString("token"),
    }, nil
}

func (p *Auth) Handle(req *outgoing.Request, resp *outgoing.Response) (err error) {

    if req.TriggerWord != p.word {
        err = errors.New("bad request trigger word in gogap-auth")
        return
    }

    if req.Token != p.token {
        err = errors.New("error auth token")
        return
    }

    return
}
```

> `gogap-auth`为唯一标识，不得与其他插件冲突

使用这个Trigger

我们在目录`$GOPATH/github.com/gogap/bearychat/outgoing/cmd/outgoing`下创建一个 `imports_mine.go`

加入如下内容

```go
package main

import (
    _ "github.com/gogap/bearychat/outgoing/triggers/auth"
    ... //此处可以添加更多的Trigger
)
```

然后再重新编译

```bash
go build
```

这个Tigger的驱动就添加完成了，我们在配置文件里按照如下配置即可:

```hocon

...

outgoing {

        hello {
            word = "!hello"
            drivers = [gogap-auth, ...] // 按顺序排列执行，如果某个驱动返回error,则中断

            ...

            gogap-auth = {
                token = "a5abc87ce0dcd169d4560385c2be9d3a"
            }
        }
}
```




### Incoming

- Request

```go
type Request struct {
    Text         string       `json:"text"`
    Notification string       `json:"notification"`
    Markdown     bool         `json:"markdown"`
    Channel      string       `json:"channel"`
    User         string       `json:"user"`
    Attachments  []Attachment `json:"attachments"`
}
```

- Response

```go
type Response struct {
    Code   int         `json:"code"`
    Error  string      `json:"error"`
    Result interface{} `json:"result"`
}
```

#### 使用方法

```go
package main

import (
    "fmt"

    "github.com/gogap/bearychat/incoming"
)

func main() {

    client := incoming.NewClient()

    req := incoming.Request{
        Text: "愿原力与你同在",
    }

    resp, err := client.Send("=ba7Ld", "ae28af7a66e16effe45f71365d6b21dc", &req)

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(resp)
}

```