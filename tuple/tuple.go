package tuple

type Tuple2 struct {
	Item1 interface{}
	Item2 interface{}
}

type Tuple3 struct {
	*Tuple2
	Item3 interface{}
}
type Tuple4 struct {
	*Tuple3
	Item4 interface{}
}
type Tuple5 struct {
	*Tuple4
	Item5 interface{}
}
type Tuple6 struct {
	*Tuple5
	Item6 interface{}
}

type Tuple7 struct {
	*Tuple6
	Item7 interface{}
}

func NewTuple2(item1, item2 interface{}) *Tuple2 {
	return &Tuple2{
		item1,
		item2,
	}
}

func NewTuple3(item1, item2, item3 interface{}) *Tuple3 {
	return &Tuple3{
		NewTuple2(item1, item2),
		item3,
	}
}
func NewTuple4(item1, item2, item3, item4 interface{}) *Tuple4 {
	return &Tuple4{
		NewTuple3(item1, item2, item3),
		item4,
	}
}

func NewTuple5(item1 interface{}, item2 interface{}, item3 interface{}, item4 interface{}, item5 interface{}) *Tuple5 {
	return &Tuple5{
		Tuple4: NewTuple4(item1, item2, item3, item4),
		Item5:  item5,
	}
}
func NewTuple6(item1, item2, item3, item4, item5, item6 interface{}) *Tuple6 {
	return &Tuple6{
		NewTuple5(item1, item2, item3, item4, item5),
		item6,
	}
}

func NewTuple7(item1, item2, item3, item4, item5, item6, item7 interface{}) *Tuple7 {
	return &Tuple7{
		NewTuple6(item1, item2, item3, item4, item5, item6),
		item7,
	}
}
