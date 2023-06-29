# Ayachan 测试服务器 - V0.7.0-rc.1

一个可以由玩家上传临时谱面用于测试的Sonolus服务器

## 使用方法

### 编译

```shell
git clone https://github.com/6qhtsk/sonolus-test-server.git
cd sonolus-test-server
go build
./sonolus-test-server
```

- 在运行程序之前，请确保在程序运行文件夹中包含服务器资源的 `sonolus` 文件夹。
- 程序在运行开始时会读取 `PORT` 环境变量来决定服务访问的端口。可以通过设置环境变量来指定端口，例如：`export PORT=8000`。

### Github容器服务

```bash
docker pull ghcr.io/6qhtsk/sonolus-test-server:latest
docker run -p 8000:8000 -v /path/to/your/sonolus:/sonolus-test-server/sonolus ghcr.io/6qhtsk/sonolus-test-server:latest
```

- 在运行容器之前，请确保将宿主机器上的服务器资源文件夹挂载到 `/sonolus-test-server/sonolus` 文件夹下。
- 注意，docker容器开放了8000端口用于访问服务。

## API文档

### POST /upload

此接口用于上传谱面。

#### 请求参数

请求参数为一个表单，表单包含以下字段：

- `title` (string): 谱面标题，必填。
- `bgm` (file): 背景音乐文件，必填。
- `chart` (string): 谱面字符串，必填。
- `difficulty` (int): 难度等级，默认25级,所有低于（包含）0级的传入都将纠正为25级。
- `hidden` (bool): 是否隐藏谱面，默认非隐藏。
- `lifetime` (int64): 谱面有效期（秒），默认6小时（21600），所有低于0的有效期和错误的输入都将纠正为21600秒。

#### 响应

响应体为一个JSON对象。

- 出错时，返回值类似：

```json
{
	"code": 303,
	"description": "上传bgm格式有误",
	"detail": "the file you upload is png (MIME image/png), not audio"
}
```

其中：

- `code` (int): 错误代码。
- `description` (string): 错误描述。
- `detail` (string): 错误详情。

- 正确时，返回值类似：

```json
{
	"uid": 123456
}
```

其中：

- `uid` (int): 上传成功后的谱面ID。