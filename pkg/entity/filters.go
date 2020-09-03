package entity

type IntFilter struct {
	Le int
	Ge int
	Lt int
	Gt int
	Ne int
}

type StringFilter struct {
	StartsWith string
	EndsWith   string
	Contains   string
	NotContain string
}

type FruitFilter struct {
	Ids   *IntFilter
	Names *StringFilter
}

type RackFilter struct {
	Ids   *IntFilter
	Cages *StringFilter
}
