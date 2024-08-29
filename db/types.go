package db

type Connection interface {
}

type ConnectionCreator func(configName string) (Connection, error)
