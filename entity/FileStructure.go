package entity

type FileStructure struct {
	ID              int64
	FullName        string
	SourceIBAN      string
	DestinationIBAN string
	Amount          float32
}
