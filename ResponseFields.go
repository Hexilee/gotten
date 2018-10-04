package gotten

type (
	Status int
)

func (status Status) Int() int {
	return int(status)
}
