package component

type AssociationType int

const (
	Extension      = 1 << iota // 0x01
	Implementation = 1 << iota // 0x02
	Composition    = 1 << iota // 0x04
	Dependency     = 1 << iota // 0x08
	supportedType  = Extension | Implementation | Composition | Dependency
)
