#!/usr/bin/env bash

set -Eeuo pipefail

project_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
backup_dir="$project_dir/.docker-mirror-backup-$(date +%Y%m%d%H%M%S)"
mkdir -p "$backup_dir"

process_file() {
  local file="$1"
  local relative="${file#"$project_dir/"}"
  local backup="$backup_dir/$relative"

  mkdir -p "$(dirname "$backup")"
  cp "$file" "$backup"
  perl -0pi -e '
    s#(^\s*FROM\s+)\$\{REGISTRY\}/(golang|alpine)(:)#$1docker.m.daocloud.io/library/$2$3#gmi;
    s#(^\s*FROM\s+)(node|nginx|postgres|redis|nats|golang|alpine)(:)#$1docker.m.daocloud.io/library/$2$3#gmi;
    s#(^\s*FROM\s+)(svhd/logto)(:)#$1docker.m.daocloud.io/$2$3#gmi;
    s#(^\s*image:\s*)(mysql|nats|postgres|redis)(:)#$1docker.m.daocloud.io/library/$2$3#gmi;
    s#(^\s*image:\s*)(emqx/emqx|tdengine/tsdb)(:)#$1docker.m.daocloud.io/$2$3#gmi;
  ' "$file"
  echo "已处理: $relative"
}

while IFS= read -r -d '' file; do
  process_file "$file"
done < <(find "$project_dir" -path "$project_dir/.git" -prune -o -type f -iname 'Dockerfile*' -print0)

while IFS= read -r -d '' file; do
  process_file "$file"
done < <(find "$project_dir" -path "$project_dir/.git" -prune -o -type f \( \
  -iname 'docker-compose*.yml' -o -iname 'docker-compose*.yaml' -o \
  -iname 'compose*.yml' -o -iname 'compose*.yaml' \
  \) -print0)

echo "备份目录: ${backup_dir#"$project_dir/"}"
