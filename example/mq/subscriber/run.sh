clear
go build -o subscriber

export MQ_DEFAULT_TYPE="rabbit"
export MQ_DEFAULT_URL="amqp://admin:maumasuk123@127.0.0.1:5672/db1"

./subscriber