clear
go build -o config

export CONFIG_STRING="test string"
export CONFIG_INT=123
export CONFIG_FLOAT=123.456
export CONFIG_BOOL=true

./config