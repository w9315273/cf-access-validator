# cf-access-validator

Cloudflare Access JWT 验证器 —— 适用于 `nginx auth_request`。  
用于在自建服务前增加 Cloudflare Zero Trust 的验证环节。  

---

## 特性
> [!NOTE]
> - 轻量级 Go 程序，无外部依赖  
> - 支持多个 App / AUD 映射  
> - 可集成到 LuCI / 自建 Web 服务  
> - 提供三种使用方式：  
> - 作为 OpenWrt IPK 包  
> - 作为 Docker 镜像  
> - 直接集成到 OpenWrt 源码编译  

## 安装方式
> [!NOTE]
> ### 1. OpenWrt .ipk
> ```
> 系统 → 软件包 → 上传软件包 → 选择.ipk → 上传并安装
> ```

---

> [!NOTE]
> ### 2. OpenWrt 源码集成
> ```
> cd package
> ```
> ```
> git clone https://github.com/w9315273/cf-access-validator
> ```
> ```
> make menuconfig   # 勾选 Network/cf-auth
> ```

---

> [!NOTE]
> ### 3. Docker
> ```
> services:
>   cf-auth:
>     image: ghcr.io/w9315273/cf-access-validator:latest
>     container_name: cf-auth
>     networks:
>       net:
>         ipv4_address: 172.18.100.2
> #    ports:
> #      - "9000:9000"    # 可选, 建议不映射到宿主机, 仅容器内访问
>     environment:
>       ADDR: "0.0.0.0:9000"    # 监听地址:端口
>       TEAM_DOMAIN: ${TEAM_DOMAIN}   # 必填, 替换为你的 Cloudflare Access 团队域名, 例如 your-team.cloudflareaccess.com
>       APP_MAP: "blog=${AUD_BLOG};grafana=${AUD_GRAFANA}"    # 必填, 替换为你的应用名称和对应的AUD值, 多个应用用分号分隔, 例如 blog=xxxx;grafana=xxxx
>     restart: always
> 
> networks:
>   net:
>     driver: bridge
>     ipam:
>       config:
>         - subnet: 172.18.100.0/24
>           gateway: 172.18.100.1

> [!TIP]
> - 建议使用.env文件存储环境变量
> - 创建一个名为 .env 的文件, 内容如下:
> - TEAM_DOMAIN=Team.cloudflareaccess.com
> - AUD_BLOG=xxxx
> - AUD_GRAFANA=xxxx

## 配置
> [!NOTE]
> 
> ### /etc/config/cf-auth
> ```
> config cf_auth 'A'
>     option enabled '1'
>     option team_domain '<team>.cloudflareaccess.com'
>     option addr '127.0.0.1:9000'
> 
> config app
>     option instance 'A'
>     option name 'openwrt'
>     list aud '<AUD>'
> ```

### Nginx 示例
> [!NOTE]
> ```
> location = /__cf_auth__ {
>     internal;
>     proxy_http_version 1.1;
>     proxy_pass http://127.0.0.1:9000/validate;
>     proxy_set_header Cf-Access-Jwt-Assertion $http_cf_access_jwt_assertion;
>     proxy_set_header X-Required-App "openwrt";
> }
> 
> location / {
>     auth_request /__cf_auth__;
>     proxy_pass xxx;
>     proxy_set_header xxx;
> }
> ```
[<img src="https://github.githubassets.com/images/icons/emoji/unicode/1f60e.png" width="20"/>](https://github.com/w9315273/cf-access-validator/wiki/20250831_21:55)
