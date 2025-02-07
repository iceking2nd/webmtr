package mtrController

import (
	"errors"
	"fmt"
	"github.com/iceking2nd/webmtr/app/models"
	"github.com/iceking2nd/webmtr/app/utils/APIResponse"
	"github.com/olekukonko/tablewriter"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tonobo/mtr/pkg/mtr"
)

func MTR(c *gin.Context) {
	params := models.DefParams
	if countStr, exist := c.GetQuery("count"); exist {
		if count, err := strconv.Atoi(countStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_COUNT_ERROR, c)
			return
		} else {
			params.COUNT = count
		}
	}
	if timeoutStr, exist := c.GetQuery("timeout"); exist {
		if timeout, err := time.ParseDuration(timeoutStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_TIMEOUT_ERROR, c)
			return
		} else {
			params.TIMEOUT = timeout
		}
	}
	if intervalStr, exist := c.GetQuery("interval"); exist {
		if interval, err := time.ParseDuration(intervalStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_INTERVAL_ERROR, c)
			return
		} else {
			params.INTERVAL = interval
		}
	}
	if hopSleepStr, exist := c.GetQuery("hop_sleep"); exist {
		if hopSleep, err := time.ParseDuration(hopSleepStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_HOP_SLEEP_ERROR, c)
			return
		} else {
			params.HOP_SLEEP = hopSleep
		}
	}
	if maxHopsStr, exist := c.GetQuery("max_hops"); exist {
		if maxHops, err := strconv.Atoi(maxHopsStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_MAX_HOPS_ERROR, c)
			return
		} else {
			params.MAX_HOPS = maxHops
		}
	}
	if maxUnknownHopsStr, exist := c.GetQuery("max_unknown_hops"); exist {
		if maxUnknownHops, err := strconv.Atoi(maxUnknownHopsStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_MAX_UNKNOWN_HOPS_ERROR, c)
			return
		} else {
			params.MAX_UNKNOWN_HOPS = maxUnknownHops
		}
	}
	if ringBufferSizeStr, exist := c.GetQuery("ring_buffer_size"); exist {
		if ringBufferSize, err := strconv.Atoi(ringBufferSizeStr); err != nil {
			APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_RING_BUFFER_SIZE_ERROR, c)
			return
		} else {
			params.RING_BUFFER_SIZE = ringBufferSize
		}
	}
	if _, exist := c.GetQuery("ptr_lookup"); exist {
		params.PTR_LOOKUP = true
	}
	if c.Request.Header.Get("Accept") == "application/json" {
		params.JsonFmt = true
	}
	if srcAddr, exist := c.GetQuery("src_addr"); exist {
		params.SrcAddr = net.ParseIP(srcAddr).String()
	}

	dstAddrs, err := net.LookupIP(c.Param("dest_addr"))
	if err != nil {
		APIResponse.ResponseError(err, http.StatusBadRequest, APIResponse.API_RESPONSE_PARSE_DESTINATION_ERROR, c)
		return
	}

	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	if !re.Match([]byte(dstAddrs[0].String())) {
		APIResponse.ResponseError(errors.New("destination is not ipv4 address"), http.StatusBadRequest, APIResponse.API_RESPONSE_DESTINATION_IS_NON_IPV4, c)
		return
	}

	m, ch, err := mtr.NewMTR(dstAddrs[0].String(), params.SrcAddr, params.TIMEOUT, params.INTERVAL, params.HOP_SLEEP, params.MAX_HOPS, params.MAX_UNKNOWN_HOPS, params.RING_BUFFER_SIZE, params.PTR_LOOKUP)
	if err != nil {
		APIResponse.ResponseError(err, http.StatusInternalServerError, APIResponse.API_RESPONSE_CREATING_MTR_PROJECT_ERROR, c)
		return
	}
	startTime := time.Now().Format("2006-01-02 15:04:05 -0700")
	go func(ch chan struct{}) {
		for {
			<-ch
		}
	}(ch)
	m.Run(ch, params.COUNT)
	if params.JsonFmt {
		APIResponse.ResponseOKWithData(m, c)
		return
	}

	result := &strings.Builder{}
	table := tablewriter.NewWriter(result)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetHeader([]string{"HOP", "Address", "Loss(%)", "Lost", "Sent", "Last", "Avg", "Best", "Worst"})
	for i := 0; i < len(m.Statistic); i++ {
		table.Append([]string{
			strconv.Itoa(i + 1),
			m.Statistic[i+1].Target,
			fmt.Sprintf("%.1f", m.Statistic[i+1].Loss()),
			strconv.Itoa(m.Statistic[i+1].Lost),
			strconv.Itoa(m.Statistic[i+1].Sent),
			fmt.Sprintf("%.1f", float64(m.Statistic[i+1].Last.Elapsed)/float64(time.Millisecond)),
			fmt.Sprintf("%.1f", m.Statistic[i+1].Avg()),
			fmt.Sprintf("%.1f", float64(m.Statistic[i+1].Best.Elapsed)/float64(time.Millisecond)),
			fmt.Sprintf("%.1f", float64(m.Statistic[i+1].Worst.Elapsed)/float64(time.Millisecond)),
		})
	}
	table.Render()
	outContent := fmt.Sprintf(
		"Start Time: %s\nTarget: %s\nSource IP: %s\nDestination IP: %s\n\n%s",
		startTime,
		c.Param("dest_addr"),
		params.SrcAddr,
		dstAddrs[0].String(),
		result,
	)
	APIResponse.ResponseOKWithData(outContent, c)
	return
}
