clear
go build -o db

export DB_DEFAULT_TYPE="mongo"
export DB_DEFAULT_URL="mongodb://localhost:27017"

./db