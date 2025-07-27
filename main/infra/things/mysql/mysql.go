package mysql

type MySQL interface {
	findWithRetry(object interface{})
}