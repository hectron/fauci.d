package vaccines

type Vaccine int

const Null = -1

const (
	Pfizer Vaccine = iota
	Moderna
	JJ
)

var vaccineToGuid = map[Vaccine]string{
	Pfizer:  "a84fb9ed-deb4-461c-b785-e17c782ef88b",
	Moderna: "779bfe52-0dd8-4023-a183-457eb100fccc",
	JJ:      "784db609-dc1f-45a5-bad6-8db02e79d44f",
}

func (v Vaccine) Guid() string {
	return vaccineToGuid[v]
}

func (v Vaccine) String() string {
	switch v {
	case Pfizer:
		return "pfizer"
	case Moderna:
		return "moderna"
	case JJ:
		return "jj"
	default:
		return ""
	}
}
