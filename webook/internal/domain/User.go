package domain

// User 领域对象，是 DDD 中的 entity
// BO(Business Object)
type User struct {
	Id       int64
	Email    string
	Password string
}
