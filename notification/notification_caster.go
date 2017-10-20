package notification

type Caster interface {
	cast(Origin, Object) (Object, error)
}
