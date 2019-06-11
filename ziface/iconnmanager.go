package ziface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	GetConnection(connID uint32) (IConnection, error)
	Len() int
	ClearConn()
}
