package handler

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type DataHandler struct {
	dataService port.DataService
}

func NewDataHandler(dataService port.DataService) port.DataHandler {
	return &DataHandler{
		dataService: dataService,
	}
}

// docapi begin get_current_dbms
// docapi method GET
// docapi summary Get the current DBMS
// docapi desc Get the current database management system that Vertex is using.
// docapi tags Admin/Data
// docapi response 200 {string} The current DBMS.
// docapi end

func (h *DataHandler) GetCurrentDbms(c *router.Context) {
	c.JSON(h.dataService.GetCurrentDbms())
}

// docapi begin migrate_to_dbms
// docapi method POST
// docapi summary Migrate to a DBMS
// docapi desc Migrate Vertex to the given database management system.
// docapi tags Admin/Data
// docapi body {MigrateToBody} The DBMS to migrate to.
// docapi response 204
// docapi response 400
// docapi response 500
// docapi end

type MigrateToBody struct {
	Dbms string `json:"dbms"`
}

func (h *DataHandler) MigrateTo(c *router.Context) {
	var body MigrateToBody
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	err = h.dataService.MigrateTo(types.DbmsName(body.Dbms))
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToMigrateToNewDbms,
			PublicMessage:  "Migration to the new DBMS failed.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}