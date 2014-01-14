package vivom

type Row interface {
	Table() string

	// Columns returns column names of the table, starting with the primary key.
	Columns() []string
}

type SelectableRow interface {
	Row

	ScanValues() []interface{}
}

type InsertableRow interface {
	Row

	GetID() int
	SetID(int)
	Validate() error
	Values() []interface{}
}

type SelectableRows interface {
	Row

	Next() SelectableRow
}
