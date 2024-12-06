package components

type IAlchemyComponent interface {
	Setup(component string) (func() error, error)
}
