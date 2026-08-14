# AIoT Backend Service

Go 后端服务，基于 Gin + GORM + Wire 构建。

## 常用命令

```bash
make init       # 安装开发工具（wire, mockgen, swag）
make gen        # 生成 GORM Gen 查询代码
make swag       # 生成 Swagger 文档
make test       # 运行测试并生成覆盖率报告
make build      # 构建可执行文件
make bootstrap  # 启动 Docker 依赖 + 数据库迁移 + 启动服务
```

## Swagger 文档生成

```bash
# 确保先安装 swag
make init

# 生成 swagger 文档到 docs/ 目录
make swag
```

生成的文件：
- `docs/swagger.yaml` — OpenAPI 2.0 规范（YAML）
- `docs/swagger.json` — OpenAPI 2.0 规范（JSON）
- `docs/docs.go` — 嵌入式 Go 文件，支持静态托管

启动服务后访问：`http://localhost:8000/swagger/index.html`

## License

MIT