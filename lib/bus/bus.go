package bus

import "fmt"

type StopDistance struct {
	RouteID     string
	DirectionID int
	FromID      string
	ToID        string
	Distance    float64
}

func (sd *StopDistance) String() string {
	return fmt.Sprintf("%s – Direction %d – From %s – To %s = %f metres", sd.RouteID, sd.DirectionID, sd.FromID, sd.ToID, sd.Distance)
}
