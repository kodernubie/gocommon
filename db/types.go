package db

type FieldOrder struct {
	Field string
	Order string
}

type FindOption struct {
	Skip  int
	Limit int
	Order []FieldOrder
}

type TableNameInformer interface {
	TableName() string
}

type BeforeInsertHandler interface {
	BeforeCreate()
}

type Connection interface {
	Save(obj interface{}) error
	Update(table interface{}, filter interface{}, update interface{}) error
	Delete(obj interface{}) error
	DeleteMany(table interface{}, filter interface{}) error
	FindById(out interface{}, id string) error
	FindOne(out interface{}, filter interface{}, options ...FindOption) error
	Find(out interface{}, filter interface{}, options ...FindOption) error
}

type ConnectionCreator func(configName string) (Connection, error)
