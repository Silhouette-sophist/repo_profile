#!/bin/bash
set -e

# 修改密码（关键：用 cypher-shell 执行 Cypher 语句）
change_password() {
  echo "正在修改 Neo4j 默认密码..."
  # 正确格式：通过 cypher-shell 传递用户名、旧密码、数据库和 Cypher 语句
  if ! cypher-shell -u neo4j -p neo4j --database "system" "ALTER CURRENT USER SET PASSWORD FROM 'neo4j' TO 'archergraph';"
  then
    echo "密码修改失败"
    exit 1
  fi
  echo "密码修改成功"
}

# 启动流程
neo4j start
change_password

# 检测neo4j状态
neo4j status