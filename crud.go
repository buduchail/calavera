package catrina

type (
	CRUD interface {
		Insert(values []Value) (id Value, e error)
		Select(id Value) (Object, error)
		SelectWhereFields(fields []string, values []Value) <-chan Object
		SelectWhereRange(field string, min, max Value) <-chan Object
		SelectWhereExpression(expr string, values []Value) <-chan Object
		Update(id Value, values []Value) error
		Delete(id Value) error
	}

	// some syntactic sugar
	Value interface{}
	Object interface{}
)
