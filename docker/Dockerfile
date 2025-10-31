FROM golang:1.25

# 设置时区
ENV TZ=Asia/Shanghai

# 安装必要依赖（含 Neo4j 所需的 OpenJDK 和系统工具）
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    git \
    wget \
    curl \
    systemd \
    systemd-sysv \
    tzdata \
    openjdk-17-jdk \
    gnupg \
    apt-transport-https \
    && ln -fs /usr/share/zoneinfo/$TZ /etc/localtime \
    && dpkg-reconfigure -f noninteractive tzdata \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# -------------------------- 新增 Neo4j 安装步骤 --------------------------
# 1. 导入 Neo4j 官方 GPG 密钥（验证包完整性）
RUN curl -fsSL https://debian.neo4j.com/neotechnology.gpg.key | gpg --dearmor -o /usr/share/keyrings/neo4j.gpg

# 2. 添加 Neo4j 官方 APT 源（适配 Debian 12 "bookworm"）
RUN echo "deb [signed-by=/usr/share/keyrings/neo4j.gpg] https://debian.neo4j.com stable 5" | tee /etc/apt/sources.list.d/neo4j.list

# 3. 更新源并安装 Neo4j（指定版本可替换为 5.x 具体版本，如 neo4j=5.18.0）
RUN apt-get update && apt-get install -y --no-install-recommends \
    neo4j \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 4. 配置 Neo4j（可选：根据需求调整配置文件）
# 允许远程访问（默认仅本地 127.0.0.1，修改为 0.0.0.0 允许外部连接）
RUN sed -i 's/#server.default_listen_address=0.0.0.0/server.default_listen_address=0.0.0.0/' /etc/neo4j/neo4j.conf
# 可选：调整内存配置（根据容器资源修改，示例为 2G 堆内存）
RUN sed -i 's/#dbms.memory.heap.initial_size=512m/dbms.memory.heap.initial_size=2g/' /etc/neo4j/neo4j.conf
RUN sed -i 's/#dbms.memory.heap.max_size=512m/dbms.memory.heap.max_size=2g/' /etc/neo4j/neo4j.conf

# -------------------------- 原有步骤保留 --------------------------
# 拷贝执行命令，并执行初始化脚本
COPY /docker/bin/init_analyzer.sh /usr/bin
COPY /docker/.gitconfig /root/.gitconfig
RUN chmod +x /usr/bin/init_analyzer.sh && /usr/bin/init_analyzer.sh && rm -f /root/.gitconfig

# 拷贝 Neo4j 密码初始化脚本（与现有 init_analyzer.sh 一起拷贝）
COPY /docker/bin/neo4j_init.sh /usr/bin/
# 赋予执行权限
RUN chmod +x /usr/bin/neo4j_init.sh

# 暴露端口（保持不变）
EXPOSE 7474 7473 7687

# 启动命令改为执行初始化脚本（替代原 CMD ["neo4j", "console"]）
CMD ["/usr/bin/neo4j_init.sh"]