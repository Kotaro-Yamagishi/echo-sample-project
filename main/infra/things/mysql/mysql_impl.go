package mysql


type MySQLImpl struct {
	sqlHandler SqlHandler
}

func NewMySQLImpl(sqlHandler SqlHandler) MySQL {
	return &MySQLImpl{
		sqlHandler: sqlHandler,
	}
}

func (m *MySQLImpl) findWithRetry(object interface{}) {
	m.sqlHandler.FindAll(object)
}