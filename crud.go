package catrina

type (
	CRUD interface {
		Insert(values []Value) (id Value, e error)
		Select(id Value) (Object, error)
		SelectWhereFields(fields []string, values []Value) (<-chan Row, error)
		SelectWhereRange(field string, min, max Value) (<-chan Row, error)
		SelectWhereExpression(expr string, values []Value) (<-chan Row, error)
		Update(id Value, values []Value) error
		Delete(id Value) error
	}

	Row struct {
		Result Object
		Error  error
	}

	// some syntactic sugar
	Value interface{}
	Object interface{}
)
