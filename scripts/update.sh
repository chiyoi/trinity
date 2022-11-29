cd $(readlink -f $(dirname $0)/..) || return

docker compose build || return
docker compose up -d || return