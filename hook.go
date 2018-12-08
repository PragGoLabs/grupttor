package grupttor

// Hook interface, by implementing hook interface you will
// be able to hang up with grupttor
type Hook interface {
	Init(grupttor *Grupttor)
}
