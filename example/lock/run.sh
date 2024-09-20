clear
go build -o lock

export CACHE_DEFAULT_ADDR="localhost:6379"
export CACHE_DEFAULT_PASS=""
export CACHE_DEFAULT_DB="0"

./lock