clear
go build -o db

export DB_DEFAULT_TYPE="mongo"
export DB_DEFAULT_DSN="mongodb://localhost:27017"

./db