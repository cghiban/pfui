package psui

const (
	PF_CMD = "/sbin/pfctl"
)

type PF struct {
}

func PFTables() (tables []string, err error) {
	// pfctl -s Tables
	//
	return tables, err
}

func PFTableShow(t string) (hosts []string, err error) {
	// pfctl -t nointernet  -T show
	// 203.0.113.1
	return hosts, err
}

func PFTableDeleteEntry(t, ip string) (err error) {
	// pfctl -t nointernet -T delete 203.0.11.3
	// 1/1 addresses deleted.

	return err
}

func PFTableAddEntry(t, ip string) (hosts []string, err error) {
	// pfctl -t spammers -T add 203.0.113.1
	//  pfctl -t nointernet -T add 203.0.113.1
	// 1 table created. /// -- optional, when we need to create the table
	// 1/1 addresses added.
	return hosts, err
}
