package p2file

type PayloadType int

const (
	PL_LS PayloadType = iota
	PL_LS_RES
	PL_GET
	PL_GET_RES
	PL_GET_RES_DONE
)

type Payload struct {
	Type       PayloadType
	DirList    []string
	TargetFile string
	Data       []byte
}
