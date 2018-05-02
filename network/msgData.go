package network

type MsgData interface {
	Marshal() (dAtA []byte, err error)
	MarshalTo(dAtA []byte) (int, error)
	Unmarshal(dAtA []byte) error
	Size() (n int)
	Reset()
	String() string
	GetMsgName() string
}
