package InfoController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/webmtr/app/models"
	"github.com/iceking2nd/webmtr/global"
	"github.com/olekukonko/tablewriter"
	"strings"
)

func Info(c *gin.Context) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(200, gin.H{
			"aoolication": "webmtr",
			"version":     global.Version,
			"build_time":  global.BuildTime,
			"default_parameters": gin.H{
				"count":            models.DefParams.COUNT,
				"timeout":          models.DefParams.TIMEOUT,
				"interval":         models.DefParams.INTERVAL,
				"hop_sleep":        models.DefParams.HOP_SLEEP,
				"max_hops":         models.DefParams.MAX_HOPS,
				"max_unknown_hops": models.DefParams.MAX_UNKNOWN_HOPS,
				"ring_buffer_size": models.DefParams.RING_BUFFER_SIZE,
				"ptr_lookup":       models.DefParams.PTR_LOOKUP,
				"json_format":      models.DefParams.JsonFmt,
				"source_address":   models.DefParams.SrcAddr,
			},
		})

	default:
		tableContent := &strings.Builder{}
		table := tablewriter.NewWriter(tableContent)
		table.SetRowLine(true)
		table.Append([]string{"the number of pings sent", fmt.Sprintf("%d", models.DefParams.COUNT)})
		table.Append([]string{"ICMP echo request timeout", fmt.Sprintf("%v", models.DefParams.TIMEOUT)})
		table.Append([]string{"ICMP echo request interval", fmt.Sprintf("%v", models.DefParams.INTERVAL)})
		table.Append([]string{"wait time between pinging next hop", fmt.Sprintf("%v", models.DefParams.HOP_SLEEP)})
		table.Append([]string{"maximum number of hops", fmt.Sprintf("%d", models.DefParams.MAX_HOPS)})
		table.Append([]string{"maximum unknown host", fmt.Sprintf("%d", models.DefParams.MAX_UNKNOWN_HOPS)})
		table.Append([]string{"cached packet buffer size", fmt.Sprintf("%d", models.DefParams.RING_BUFFER_SIZE)})
		table.Append([]string{"disable DNS lookup", fmt.Sprintf("%v", models.DefParams.PTR_LOOKUP)})
		table.Append([]string{"output as JSON", fmt.Sprintf("%v", models.DefParams.JsonFmt)})
		table.Append([]string{"bind the outgoing socket to ADDRESS", fmt.Sprintf("%s", models.DefParams.SrcAddr)})
		table.Render()
		content := fmt.Sprintf(`
webmtr %s(%s)
Build at: %s
Default parameters:
%s
`, global.Version, global.GitCommit, global.BuildTime, tableContent.String())
		c.String(200, content)
	}
}
