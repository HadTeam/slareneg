package GameOperator

import (
	"server/ApiProvider/pkg/DataOperator"
	_ "server/ApiProvider/pkg/DataOperator/Local"
)

var data DataOperator.DataSource

func init() {
	// TODO: Link data source

}
