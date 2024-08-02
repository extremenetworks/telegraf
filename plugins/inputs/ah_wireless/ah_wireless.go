package ah_wireless

import (
	"log"
	"os"
	"bytes"
	"strconv"
	"syscall"
	"strings"
	"sync"
	"time"
	"os/exec"
	"fmt"
	"unsafe"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"golang.org/x/sys/unix"
)

var (
	offsetsMutex = new(sync.Mutex)
	newLineByte  = []byte("\n")
)


type Ah_wireless struct {
	fd     			int
	fe_fd			uintptr
	intf_m  		map[string]map[string]string
	arp_m			map[string]string
	Ifname			[]string	`toml:"ifname"`
	closed 			chan 		struct{}
	numclient 		int
	timer_count		uint8
	entity			map[string]map[string]unsafe.Pointer
	last_rf_stat	[4]awestats
	last_clt_stat	[4][50]ah_ieee80211_sta_stats_item
}


func ah_ioctl(fd uintptr, op, argp uintptr) error {
	        _, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(op), argp)
        if errno != 0 {
                return errno
        }
        return nil
}


const sampleConfig = `
[[inputs.ah_wireless]]
  interval = "5s"
  ifname = ["wifi0","wifi1"]
`
func NewAh_wireless(id int) *Ah_wireless {
	var err error
	// Create RAW  Socket.
        fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
        if err != nil {
                return nil
        }


	if id != -1 {
		open(fd, id)
	}

	return &Ah_wireless{
                fd:       fd,
				timer_count: 0,
        }

}

func getHDDStat(fd int, ifname string) *ah_ieee80211_hdd_stats {

        var cfg *ieee80211req_cfg_hdd
        cfg = new(ieee80211req_cfg_hdd)


        /* first 4 bytes is subcmd */
        cfg.cmd = AH_IEEE80211_GET_HDD_STATS;

        iwp := iw_point{pointer: unsafe.Pointer(cfg)}

        request := iwreq{data: iwp}


     //   request.data.length = uint16(unsafe.Sizeof(cfgss));
		request.data.length = VAP_BUFF_SIZE

        copy(request.ifrn_name[:], ah_ifname_radio2vap(ifname))

		offsetsMutex.Lock()

        if err := ah_ioctl(uintptr(fd), IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getHDDStat ioctl data error %s",err)
				offsetsMutex.Unlock()
                return nil
        }
		offsetsMutex.Unlock()

        return &cfg.hdd_stats

}

func getAtrTbl(fd int, ifname string) *ah_ieee80211_atr_user {
	var cfg *ieee80211req_cfg_atr
        cfg = new(ieee80211req_cfg_atr)

	/* first 4 bytes is subcmd */
        cfg.cmd = AH_IEEE80211_GET_ATR_TBL;

	iwp := iw_point{pointer: unsafe.Pointer(cfg)}

	request := iwreq{data: iwp}

	request.data.length = VAP_BUFF_SIZE

	copy(request.ifrn_name[:], ah_ifname_radio2vap(ifname))

	offsetsMutex.Lock()

	if err := ah_ioctl(uintptr(fd), IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getAtrTbl ioctl data error %s",err)
				offsetsMutex.Unlock()
                return nil
        }

	offsetsMutex.Unlock()

	return &cfg.atr

}


func getRFStat(fd int, ifname string) *awestats {

	var p *awestats
	p = new(awestats)

	request := IFReqData{Data: uintptr(unsafe.Pointer(p))}
	copy(request.Name[:], ifname)

	offsetsMutex.Lock()

        if err := ah_ioctl(uintptr(fd), SIOCGRADIOSTATS, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getRFStat ioctl data error %s",err)
				offsetsMutex.Unlock()
                return nil
        }

	offsetsMutex.Unlock()

	return p
}

func getStaStat(fd int, ifname string, buf unsafe.Pointer,count int) *ah_ieee80211_get_wifi_sta_stats {

        var cfg *ieee80211req_cfg_sta
        cfg = new(ieee80211req_cfg_sta)

        /* first 4 bytes is subcmd */
        cfg.cmd = IEEE80211_GET_WIFI_STA_STATS
		cfg.wifi_sta_stats.count = uint16(count)
		cfg.wifi_sta_stats.pointer = buf

        iwp := iw_point{pointer: unsafe.Pointer(cfg)}

        request := iwreq{data: iwp}

        request.data.length = VAP_BUFF_SIZE

        copy(request.ifrn_name[:], ah_ifname_radio2vap(ifname))

        offsetsMutex.Lock()

        if err := ah_ioctl(uintptr(fd), IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getStaStat ioctl data error %s",err)
				offsetsMutex.Unlock()
                return nil
        }

        offsetsMutex.Unlock()

        return &cfg.wifi_sta_stats
}

func getNumAssocs (fd int, ifname string) uint32 {

	ird := iwreq_data{}

        /* first 4 bytes is subcmd */
        ird.data = IEEE80211_PARAM_NUM_ASSOCS

        request := iwreq_clt{u: ird}

        copy(request.ifr_name[:], ah_ifname_radio2vap(ifname))


        offsetsMutex.Lock()

        if err := ah_ioctl(uintptr(fd), IEEE80211_IOCTL_GETPARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getNumAssocs ioctl data error %s",err)
				offsetsMutex.Unlock()
                return 0
        }


        offsetsMutex.Unlock()

	return uint32(request.u.data)
}

func getOneStaInfo(fd int, ifname string, mac_ad [MACADDR_LEN]uint8) *ah_ieee80211_sta_info {
        var cfg *ieee80211req_cfg_one_sta
        cfg = new(ieee80211req_cfg_one_sta)

        /* first 4 bytes is subcmd */
        cfg.cmd = AH_IEEE80211_GET_ONE_STA_INFO;
	cfg.sta_info.mac = mac_ad
        iwp := iw_point{pointer: unsafe.Pointer(cfg)}
        request := iwreq{data: iwp}
        request.data.length = VAP_BUFF_SIZE
        copy(request.ifrn_name[:], ifname)

        offsetsMutex.Lock()
        if err := ah_ioctl(uintptr(fd), IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getOneStaInfo ioctl data error %s",err)
				offsetsMutex.Unlock()
                return nil
        }

        offsetsMutex.Unlock()

        return &cfg.sta_info

}

func getOneSta(fd int, ifname string, mac_ad [MACADDR_LEN]uint8) unsafe.Pointer {
        var cfg *ieee80211req_cfg_one_sta_info
        cfg = new(ieee80211req_cfg_one_sta_info)

        /* first 4 bytes is subcmd */
        cfg.cmd = AH_IEEE80211_GET_ONE_STA
        cfg.mac = mac_ad
        iwp := iw_point{pointer: unsafe.Pointer(cfg)}
        request := iwreq{data: iwp}
        request.data.length = VAP_BUFF_SIZE
        copy(request.ifrn_name[:], ah_ifname_radio2vap(ifname))

        offsetsMutex.Lock()
	if err := ah_ioctl(uintptr(fd), IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("getOneSta ioctl data error %s",err)
		offsetsMutex.Unlock()
		return nil
        }
	offsetsMutex.Unlock()

        return request.data.pointer

}

func getProcNetDev(ifname string) *ah_dcd_dev_stats {
	log.Printf("getProcNetDev called: %s",ifname)
	table, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil;
	}

	lines := bytes.Split([]byte(table), newLineByte)

	var intfname string
	var stats = new(ah_dcd_dev_stats)
	comp := fmt.Sprintf(" %s:",ifname)

  for  _, curLine := range lines {
    if strings.Contains(string(curLine), comp) {
        n, _ := fmt.Sscanf(string(curLine),
        "%s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
                            &intfname,
                            &stats.rx_bytes,
                           &stats.rx_packets,
                           &stats.rx_errors,
                           &stats.rx_dropped,
                           &stats.rx_fifo_errors,
                           &stats.rx_frame_errors,
                           &stats.rx_compressed,  /* missing for <= 1 */
                           &stats.rx_multicast, /* missing for <= 1 */
                           &stats.tx_bytes, /* missing for 0 */
                           &stats.tx_packets,
                           &stats.tx_errors,
                           &stats.tx_dropped,
                           &stats.tx_fifo_errors,
                           &stats.collisions,
                           &stats.tx_carrier_errors,
                           &stats.tx_compressed,
                           &stats.tx_multicast,
                           &stats.rx_unicast,
                           &stats.rx_broadcast,
                           &stats.tx_unicast,
                           &stats.tx_broadcast)

						   log.Printf("%d %s",n,intfname)
    }
  }

	return stats
}

func getIfIndex(fd int, ifname string) int {
        ifr, err := unix.NewIfreq(ifname)

        if err != nil {
                log.Printf("failed to create ifreq  %v", err)
        }

	offsetsMutex.Lock()

        if err := unix.IoctlIfreq(fd, unix.SIOCGIFINDEX, ifr); err != nil {
                log.Printf("getIfIndex ioctl error %s",err)
				offsetsMutex.Unlock()
                return -1
        }

	offsetsMutex.Unlock()
	return int(ifr.Uint32())

}

func load_ssid(t *Ah_wireless, ifname string) {

	for i := 1; i < 1024; i++ {
		vifname := ifname + "." + strconv.Itoa(i)
		log.Printf(vifname)

		app := "wl"

		arg0 := "-i"
		arg1 := vifname
		arg2 := "status"
		//  arg3 := "| grep \"SSID: \"\\\""
		//log.Printf(app + " " + arg0 + " " + arg1 + " " + arg2)

		cmd := exec.Command(app, arg0, arg1, arg2)
		output, err := cmd.Output()

		if err != nil {
			//log.Printf(err.Error())
			return
		}

		lines := strings.Split(string(output),"\n")

		temp  := strings.Split(lines[0]," ")

		ssid := strings.Trim(temp[1], "\"")
		log.Printf("Adding: t.intf_m[%s][%s] = %s",ifname,ssid,vifname)
		t.intf_m[ifname][ssid] = vifname
	}
}

func load_arp_table(t *Ah_wireless) {

	log.Printf("load_arp")
	app := "arp"
	arg := "-v"

	cmd := exec.Command(app, arg)

	arp_str, err := cmd.Output()

	if err != nil {
		log.Printf(err.Error())
		return
	}

	arp_lines := strings.Split(string(arp_str),"\n")

	for i :=0; i<len(arp_lines); i++ {
		log.Printf("Line[%d]:[%s]" ,i,arp_lines[i])
		if len(arp_lines[i]) > 1 {
			arp_eliments := strings.Split(arp_lines[i]," ")

			log.Printf("Adding: t.arp_m[%s] = %s",arp_eliments[3], arp_eliments[1])
			t.arp_m[arp_eliments[3]] = arp_eliments[1]
		}
	}

}

func getFeIpnetScore(fd uintptr, clmac [MACADDR_LEN]uint8) int {

	msg := ah_flow_get_sta_net_health_msg{
					mac: clmac,
					net_health_score: 0,
				}
	ihdr := ah_fe_ioctl_hdr{
					retval: -1,
					msg_type: AH_GET_STATION_NETWORK_HEALTH,
					msg_size: uint16(unsafe.Sizeof(msg)),
				}
        dev_msg := ah_fw_dev_msg{
					hdr: ihdr,
					data: msg,
				}

        offsetsMutex.Lock()

        if err := ah_ioctl(fd, AH_FE_IOCTL_FLOW, uintptr(unsafe.Pointer(&dev_msg))); err != nil {
                log.Printf("getFeIpnetScore ioctl data error %s",err)
		offsetsMutex.Unlock()
                return -1
	}

	offsetsMutex.Unlock()

	if dev_msg.hdr.retval < 0 {
		log.Printf("Open ioctl data erro")
		return -1
	}

	return dev_msg.data.net_health_score
}

func getFeServerIp(fd uintptr, clmac [MACADDR_LEN]uint8) *ah_flow_get_sta_server_ip_msg {

        msg := ah_flow_get_sta_server_ip_msg{
                                        mac: clmac,
                                }
        ihdr := ah_fe_ioctl_hdr{
                                        retval: -1,
                                        msg_type: AH_FLOW_GET_STATION_SERVER_IP,
                                        msg_size: uint16(unsafe.Sizeof(msg)),
                                }
        dev_msg := ah_fw_dev_ip_msg{
                                        hdr: ihdr,
                                        data: msg,
                                }

        offsetsMutex.Lock()

        if err := ah_ioctl(fd, AH_FE_IOCTL_FLOW, uintptr(unsafe.Pointer(&dev_msg))); err != nil {
                log.Printf("getFeServerIp ioctl data error %s",err)
                offsetsMutex.Unlock()
                return nil
        }

        offsetsMutex.Unlock()


        if dev_msg.hdr.retval < 0 {
                log.Printf("Open ioctl data erro")
                return nil
        }

        return &dev_msg.data
}

func open(fd, id int) *Ah_wireless {

	//getProcNetDev("wifi1")
	//defer unix.Close(fd)

	return &Ah_wireless{fd: fd, closed: make(chan struct{})}
}

func (t *Ah_wireless) SampleConfig() string {
	return sampleConfig
}

func (t *Ah_wireless) Description() string {
	return "Hive OS wireless stat"
}

func (t *Ah_wireless) Init() error {
	return nil
}

func Gather_Rf_Stat(t *Ah_wireless, acc telegraf.Accumulator) error {
	var ii int
	ii = 0
	for _, intfName := range t.Ifname {

		var rfstat *awestats
		var devstats *ah_dcd_dev_stats
		var ifindex int
		var atrStat *ah_ieee80211_atr_user
		var hddStat *ah_ieee80211_hdd_stats

		var idx			int 
		var tmp_count1	int64
		var tmp_count2	int64
		//var tx_ok		uint64
		//var rx_ok		uint64
		var tx_total	int64
		var rx_total	int64
		var tmp_count3	int32
		var tmp_count4	int32
		var tot_tx_bitrate_retries uint32
		var tot_rx_bitrate_retries uint32

		var rf_report	ah_dcd_stats_report_int_data

		rfstat  = getRFStat(t.fd, intfName)
		if (rfstat == nil) {
			continue
		}

		ifindex = getIfIndex(t.fd, intfName)
		if (ifindex <= 0) {
			continue
		}
		devstats = getProcNetDev(intfName)
		if (devstats == nil) {
			continue
		}
		atrStat = getAtrTbl(t.fd, intfName)
		if (atrStat == nil) {
			continue
		}
		hddStat = getHDDStat(t.fd, intfName)
		if (hddStat == nil) {
			continue
		}


		/* We need check and aggregation Tx/Rx bit rate distribution
 		* prcentage, if the bit rate equal in radio interface or client reporting.
 		*/

		for i := 0; i < NS_HW_RATE_SIZE; i++{

			if ((rfstat.ast_rx_rate_stats[i].ns_rateKbps == 0) && (rfstat.ast_tx_rate_stats[i].ns_rateKbps == 0)) {
				continue
			}

			for j := (i+1); j < NS_HW_RATE_SIZE; j++ {
				if ((rfstat.ast_rx_rate_stats[i].ns_rateKbps != 0) && (rfstat.ast_rx_rate_stats[i].ns_rateKbps == rfstat.ast_rx_rate_stats[j].ns_rateKbps)) {
					rfstat.ast_rx_rate_stats[i].ns_unicasts += rfstat.ast_rx_rate_stats[j].ns_unicasts;
					rfstat.ast_rx_rate_stats[i].ns_retries += rfstat.ast_rx_rate_stats[j].ns_retries;
					rfstat.ast_rx_rate_stats[j].ns_rateKbps = 0;
					rfstat.ast_rx_rate_stats[j].ns_unicasts = 0;
					rfstat.ast_rx_rate_stats[j].ns_retries = 0;
				}
				if ((rfstat.ast_tx_rate_stats[i].ns_rateKbps != 0) && (rfstat.ast_tx_rate_stats[i].ns_rateKbps == rfstat.ast_tx_rate_stats[j].ns_rateKbps)) {
					rfstat.ast_tx_rate_stats[i].ns_unicasts += rfstat.ast_tx_rate_stats[j].ns_unicasts;
					rfstat.ast_tx_rate_stats[i].ns_retries += rfstat.ast_tx_rate_stats[j].ns_retries;
					rfstat.ast_tx_rate_stats[j].ns_rateKbps = 0;
					rfstat.ast_tx_rate_stats[j].ns_unicasts = 0;
					rfstat.ast_tx_rate_stats[j].ns_retries = 0;
				}
			}
		}



/* Rate Calculation Copied from DCD code */

	/* Tx/Rx bit rate distribution */
	for idx = 0; idx < NS_HW_RATE_SIZE; idx++ {
		tx_total += int64(reportGetDiff(rfstat.ast_tx_rate_stats[idx].ns_unicasts,
			t.last_rf_stat[ii].ast_tx_rate_stats[idx].ns_unicasts))

		rx_total += int64(reportGetDiff(rfstat.ast_rx_rate_stats[idx].ns_unicasts,
			t.last_rf_stat[ii].ast_rx_rate_stats[idx].ns_unicasts))

		tot_tx_bitrate_retries += reportGetDiff(rfstat.ast_tx_rate_stats[idx].ns_retries,
			t.last_rf_stat[ii].ast_tx_rate_stats[idx].ns_retries)

		tot_rx_bitrate_retries += reportGetDiff(rfstat.ast_rx_rate_stats[idx].ns_retries,
			t.last_rf_stat[ii].ast_rx_rate_stats[idx].ns_retries)

	}
	//tx_ok = uint64(tx_total)
	//rx_ok = uint64(rx_total)



	for idx = 0; idx < NS_HW_RATE_SIZE; idx++ {

		tmp_count3 = int32(rfstat.ast_tx_rate_stats[idx].ns_unicasts - t.last_rf_stat[ii].ast_tx_rate_stats[idx].ns_unicasts)
		if (tx_total > 0 && tmp_count3 > 0) {
			rf_report.tx_bit_rate[idx].rate_dtn = uint8((int64(tmp_count3) * 100) / tx_total)
		} else {
			rf_report.tx_bit_rate[idx].rate_dtn = 0;
		}
		tmp_count4 = int32(rfstat.ast_rx_rate_stats[idx].ns_unicasts - t.last_rf_stat[ii].ast_rx_rate_stats[idx].ns_unicasts)
		if (rx_total > 0 && tmp_count4 > 0) {
			rf_report.rx_bit_rate[idx].rate_dtn = uint8((int64(tmp_count4) * 100) / rx_total)
		} else {
			rf_report.rx_bit_rate[idx].rate_dtn = 0;
		}

		/* Tx/Rx bit rate success distribution */
		tmp_count1 = int64(rfstat.ast_tx_rate_stats[idx].ns_retries - t.last_rf_stat[ii].ast_tx_rate_stats[idx].ns_retries)
		tmp_count2 = tmp_count1 + int64(tmp_count3)
		if (tmp_count2 > 0 && rf_report.tx_bit_rate[idx].rate_dtn > 0) {
			rf_report.tx_bit_rate[idx].rate_suc_dtn = uint8((int64(tmp_count3) * 100) / tmp_count2)
			if (rf_report.tx_bit_rate[idx].rate_suc_dtn > 100) {
				rf_report.tx_bit_rate[idx].rate_suc_dtn  = 100
				log.Printf("DCD stats report int data process: rate_suc_dtn1 is more than 100%\n")
			}
		} else {
			rf_report.tx_bit_rate[idx].rate_suc_dtn = 0;
		}

		tmp_count1 = int64(rfstat.ast_rx_rate_stats[idx].ns_retries - t.last_rf_stat[ii].ast_rx_rate_stats[idx].ns_retries)
		tmp_count2 = tmp_count1 + int64(tmp_count4)
		if (tmp_count2 > 0 && rf_report.rx_bit_rate[idx].rate_dtn > 0) {
			rf_report.rx_bit_rate[idx].rate_suc_dtn = uint8((int64(tmp_count4) * 100) / tmp_count2)
			if (rf_report.rx_bit_rate[idx].rate_suc_dtn > 100) {
				rf_report.rx_bit_rate[idx].rate_suc_dtn = 100;
				log.Printf("DCD stats report int data process: rate_suc_dtn2 is more than 100%\n");
			}

		} else {
			rf_report.rx_bit_rate[idx].rate_suc_dtn = 0;
		}
		rf_report.tx_bit_rate[idx].kbps = rfstat.ast_tx_rate_stats[idx].ns_rateKbps;
		rf_report.rx_bit_rate[idx].kbps = rfstat.ast_rx_rate_stats[idx].ns_rateKbps;
	}

/* Rate calculation copied from DCD code */

		fields := map[string]interface{}{

			"name_keys":					intfName,
			"ifindex_keys":					ifindex,

		}

			if atrStat.count > 0 {

				rx_util := atrStat.atr_info[atrStat.count - 1].rxf_pcnt
				tx_util := atrStat.atr_info[atrStat.count - 1].txf_pcnt

				total_util := atrStat.atr_info[atrStat.count - 1].rxc_pcnt

				var chan_util int
				if total_util > 100 {
					chan_util = 100
				} else {
					chan_util = int(total_util)
				}

				if (total_util > (rx_util + tx_util)) {
					fields["interferenceUtilization_min"]		= total_util - rx_util - tx_util
					fields["interferenceUtilization_max"]		= total_util - rx_util - tx_util
					fields["interferenceUtilization_avg"]		= total_util - rx_util - tx_util
				} else {
					fields["interferenceUtilization_min"]		= 0
					fields["interferenceUtilization_max"]		= 0
					fields["interferenceUtilization_avg"]		= 0
				}

				fields["channelUtilization_min"]			= chan_util
				fields["channelUtilization_max"]			= chan_util
				fields["channelUtilization_avg"]			= chan_util

				fields["txUtilization_min"]				= tx_util
				fields["txUtilization_max"]				= tx_util
				fields["txUtilization_avg"]				= tx_util

				fields["rxUtilization_min"]				= rx_util
				fields["rxUtilization_max"]				= rx_util
				fields["rxUtilization_avg"]				= rx_util

				fields["rxInbssUtilization_min"]			= atrStat.atr_info[atrStat.count - 1].rxf_inbss
				fields["rxInbssUtilization_max"]			= atrStat.atr_info[atrStat.count - 1].rxf_inbss
				fields["rxInbssUtilization_avg"]			= atrStat.atr_info[atrStat.count - 1].rxf_inbss

				fields["rxObssUtilization_min"]				= atrStat.atr_info[atrStat.count - 1].rxf_obss
				fields["rxObssUtilization_max"]				= atrStat.atr_info[atrStat.count - 1].rxf_obss
				fields["rxObssUtilization_avg"]				= atrStat.atr_info[atrStat.count - 1].rxf_obss
			} else {
				fields["channelUtilization_min"]			= 0
				fields["channelUtilization_max"]			= 0
				fields["channelUtilization_avg"]			= 0

				fields["txUtilization_min"]				= 0
				fields["txUtilization_max"]				= 0
				fields["txUtilization_avg"]				= 0

				fields["rxUtilization_min"]				= 0
				fields["rxUtilization_max"]				= 0
				fields["rxUtilization_avg"]				= 0

				fields["rxInbssUtilization_min"]			= 0
				fields["rxInbssUtilization_max"]			= 0
				fields["rxInbssUtilization_avg"]			= 0

				fields["rxObssUtilization_min"]				= 0
				fields["rxObssUtilization_max"]				= 0
				fields["rxObssUtilization_avg"]				= 0
			}

			fields["wifinterferenceUtilization_min"]			= 0
			fields["wifinterferenceUtilization_max"]			= 0
			fields["wifinterferenceUtilization_avg"]			= 0

			fields["noise_min"]						= rfstat.ast_noise_floor
			fields["noise_max"]						= rfstat.ast_noise_floor
			fields["noise_avg"]						= rfstat.ast_noise_floor

			fields["crcErrorRate_min"]					= rfstat.phy_stats.ast_rx_crcerr
			fields["crcErrorRate_max"]					= rfstat.phy_stats.ast_rx_crcerr
			fields["crcErrorRate_avg"]					= rfstat.phy_stats.ast_rx_crcerr


			fields["txPackets"]						= devstats.tx_packets
			fields["txErrors"]						= devstats.tx_errors
			fields["txDropped"]						= devstats.tx_dropped
			fields["txHwDropped"]						= rfstat.ast_as.ast_tx_shortpre + rfstat.ast_as.ast_tx_xretries + rfstat.ast_as.ast_tx_fifoerr
			fields["txSwDropped"]						= devstats.tx_dropped
			fields["txBytes"]						= devstats.tx_bytes
			fields["txRetryCount"]						= rfstat.phy_stats.ast_tx_shortretry + rfstat.phy_stats.ast_tx_longretry

			fields["txRate_min"]						= rfstat.ast_tx_rate_stats[0].ns_rateKbps
			fields["txRate_max"]						= rfstat.ast_tx_rate_stats[0].ns_rateKbps
			fields["txRate_avg"]						= rfstat.ast_tx_rate_stats[0].ns_rateKbps

			fields["txUnicastPackets"]					= rfstat.ast_tx_rate_stats[0].ns_unicasts
			fields["txMulticastPackets"]					= devstats.tx_multicast
			fields["txMulticastBytes"]					= rfstat.ast_as.ast_tx_mcast_bytes
			fields["txBcastBytes"]						= rfstat.ast_as.ast_tx_bcast_bytes
			fields["txBcastPackets"]					= devstats.tx_broadcast

			fields["rxPackets"]						= devstats.rx_packets
			fields["rxErrors"]						= devstats.rx_errors
			fields["rxDropped"]						= devstats.rx_dropped
			fields["rxBytes"]						= devstats.rx_bytes
			fields["rxRetryCount"]						= rfstat.ast_rx_retry

			fields["rxRate_min"]						= rfstat.ast_rx_rate_stats[0].ns_rateKbps
			fields["rxRate_max"]						= rfstat.ast_rx_rate_stats[0].ns_rateKbps
			fields["rxRate_avg"]						= rfstat.ast_rx_rate_stats[0].ns_rateKbps

			fields["rxMulticastBytes"]					= rfstat.ast_rx_mcast_bytes
			fields["rxMulticastPackets"]					= devstats.rx_multicast
			fields["rxBcastPackets"]					= devstats.rx_broadcast
			fields["rxBcastBytes"]						= rfstat.ast_rx_bcast_bytes

			fields["bsSpCnt"]						= hddStat.bs_sp_cnt
			fields["snrSpCnt"]						= hddStat.snr_sp_cnt
			fields["snAnswerCnt"]						= reportGetDiff(hddStat.sn_answer_cnt, hddStat.sn_answer_cnt)
			fields["rxPrbSpCnt"]						= rfstat.is_rx_hdd_probe_sup
			fields["rxAuthCnt"]						= rfstat.is_rx_hdd_auth_sup

			fields["txBitrateSuc"]						= rfstat.ast_tx_rix_invalids
			fields["rxBitrateSuc"]						= rfstat.ast_rx_rix_invalids

				for i := 0; i < NS_HW_RATE_SIZE; i++{
					kbps := fmt.Sprintf("kbps_%d_rxRateStats",i)
					rateDtn := fmt.Sprintf("rateDtn_%d_rxRateStats",i)
					rateSucDtn := fmt.Sprintf("rateSucDtn_%d_rxRateStats",i)
					fields[kbps]					= rf_report.rx_bit_rate[i].kbps
					fields[rateDtn]					= rf_report.rx_bit_rate[i].rate_dtn
					fields[rateSucDtn]				= rf_report.rx_bit_rate[i].rate_suc_dtn
				}


				for i := 0; i < NS_HW_RATE_SIZE; i++{
					kbps := fmt.Sprintf("kbps_%d_txRateStats",i)
					rateDtn := fmt.Sprintf("rateDtn_%d_txRateStats",i)
					rateSucDtn := fmt.Sprintf("rateSucDtn_%d_txRateStats",i)
					fields[kbps]					= rf_report.tx_bit_rate[i].kbps
					fields[rateDtn]					= rf_report.tx_bit_rate[i].rate_dtn
					fields[rateSucDtn]				= rf_report.tx_bit_rate[i].rate_suc_dtn
				}

			fields["clientCount"]						= t.numclient
			fields["lbSpCnt"]							= hddStat.lb_sp_cnt
			fields["rxProbeSup"]						= rfstat.is_rx_hdd_probe_sup
			fields["rxSwDropped"]						= devstats.rx_dropped
			fields["rxUnicastPackets"]					= rfstat.ast_rx_rate_stats[0].ns_unicasts


			acc.AddGauge("RfStats", fields, nil)
			t.last_rf_stat[ii] = *rfstat
			ii++
		}

		return nil
}

func Gather_Client_Stat(t *Ah_wireless, acc telegraf.Accumulator) error {
	tags := map[string]string{
	}


	fields2 := map[string]interface{}{
	}

	var total_client_count int
	var ii int
	total_client_count = 0
	ii = 0


        for _, intfName2 := range t.Ifname {

		var cltstat *ah_ieee80211_get_wifi_sta_stats
		var ifindex2 int
		var numassoc int
		var stainfo *ah_ieee80211_sta_info
		var client_mac string

		var tx_total		int64
		var rx_total		int64
		//var tx_ok			uint64
		//var rx_ok			uint64
		//var tot_rate_frame	uint32
		var tx_retries		uint32
		//var tx_retry_rate	uchar
		var idx				int
		//var conn_score		uint32
		var tmp_count1		int32
		var tmp_count2		int32
		var tmp_count3		uint32
		var tmp_count4		uint32
		var tmp_count5		uint64
		var tmp_count6		uint64
		var rf_report		ah_dcd_stats_report_int_data
		var tot_tx_bitrate_retries uint32
		var tot_rx_bitrate_retries	uint32

		var client_ssid string

		numassoc = int(getNumAssocs(t.fd, intfName2))

		if(numassoc == 0) {
			continue
		}

		total_client_count = total_client_count + numassoc

		clt_item := make([]ah_ieee80211_sta_stats_item, numassoc)


		ifindex2 = getIfIndex(t.fd, intfName2)
		if(ifindex2 <= 0 ) {
			continue
		}

		cltstat = getStaStat(t.fd, intfName2, unsafe.Pointer(&clt_item[0]),  numassoc)

		for cn := 0; cn < numassoc; cn++ {
			//if ( clt_item[cn] == nil) {
			//	continue
			//}
			client_ssid = string(bytes.Trim(clt_item[cn].ns_ssid[:], "\x00"))

			if(clt_item[cn].ns_mac[0] !=0 || clt_item[cn].ns_mac[1] !=0 || clt_item[cn].ns_mac[2] !=0 || clt_item[cn].ns_mac[3] !=0 || clt_item[cn].ns_mac[4] != 0 || clt_item[cn].ns_mac[5]!=0) {
				cintfName := t.intf_m[intfName2][client_ssid]
				client_mac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",clt_item[cn].ns_mac[0],clt_item[cn].ns_mac[1],clt_item[cn].ns_mac[2],clt_item[cn].ns_mac[3],clt_item[cn].ns_mac[4],clt_item[cn].ns_mac[5])

				//log.Printf("Going to call getOneStaInfo cintfName = %s mac = %s",cintfName,client_mac)

				stainfo = getOneStaInfo(t.fd, cintfName, clt_item[cn].ns_mac)

				if(stainfo==nil) {
					log.Printf("Error in getOneStaInfo")
					continue
				}

				if stainfo.rssi == 0 {
					continue
				}
			} else {
				stainfo = nil
				continue
			}

			f := init_fe()
			ipnet_score := getFeIpnetScore(f.Fd(), clt_item[cn].ns_mac)
			sta_ip := getFeServerIp(f.Fd(), clt_item[cn].ns_mac)
			f.Close()

			//client_mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",clt_item[cn].ns_mac[0],clt_item[cn].ns_mac[1],clt_item[cn].ns_mac[2],clt_item[cn].ns_mac[3],clt_item[cn].ns_mac[4],clt_item[cn].ns_mac[5])
			cfgptr := getOneSta(t.fd, intfName2, clt_item[cn].ns_mac)

			if(cfgptr == nil) {
				continue
			}

			var onesta *ieee80211req_sta_info = (*ieee80211req_sta_info)(cfgptr)

			var clt_last_stats *saved_stats = (*saved_stats)(t.entity[intfName2][client_mac])


			if (clt_last_stats != nil) {
				clt_last_stats.tx_airtime_min	= ((clt_last_stats.tx_airtime_min /10) / (60 *1000))
				clt_last_stats.tx_airtime_max	= ((clt_last_stats.tx_airtime_max /10) / (60 *1000))
				clt_last_stats.tx_airtime_average	= ((clt_last_stats.tx_airtime_average /10) / (60 *1000))

				clt_last_stats.rx_airtime_min	= ((clt_last_stats.rx_airtime_min /10) / (60 *1000))
				clt_last_stats.rx_airtime_max	= ((clt_last_stats.rx_airtime_max /10) / (60 *1000))
				clt_last_stats.rx_airtime_average	= ((clt_last_stats.rx_airtime_average /10) / (60 *1000))

				if (clt_last_stats.tx_airtime_min > 100) {
					clt_last_stats.tx_airtime_min = 100
				}

				if (clt_last_stats.tx_airtime_max > 100) {
					clt_last_stats.tx_airtime_max = 100
				}

				if (clt_last_stats.tx_airtime_average > 100) {
					clt_last_stats.tx_airtime_average = 100
				}

				if (clt_last_stats.rx_airtime_min > 100) {
					clt_last_stats.rx_airtime_min = 100
				}

				if (clt_last_stats.rx_airtime_max > 100) {
					clt_last_stats.rx_airtime_max = 100
				}

				if (clt_last_stats.rx_airtime_average > 100) {
					clt_last_stats.rx_airtime_average = 100
				}
			}


			/* We need check and aggregation Tx/Rx bit rate distribution
 			* prcentage, if the bit rate equal in radio interface or client reporting.
 			*/

			for i := 0; i < NS_HW_RATE_SIZE; i++{

				if ((clt_item[cn].ns_rx_rate_stats[i].ns_rateKbps == 0) && (clt_item[cn].ns_tx_rate_stats[i].ns_rateKbps == 0)) {
					continue
				}

				for j := (i+1); j < NS_HW_RATE_SIZE; j++ {
					if ((clt_item[cn].ns_rx_rate_stats[i].ns_rateKbps != 0) && (clt_item[cn].ns_rx_rate_stats[i].ns_rateKbps == clt_item[cn].ns_rx_rate_stats[j].ns_rateKbps)) {
						clt_item[cn].ns_rx_rate_stats[i].ns_unicasts += clt_item[cn].ns_rx_rate_stats[j].ns_unicasts;
						clt_item[cn].ns_rx_rate_stats[i].ns_retries += clt_item[cn].ns_rx_rate_stats[j].ns_retries;
						clt_item[cn].ns_rx_rate_stats[j].ns_rateKbps = 0;
						clt_item[cn].ns_rx_rate_stats[j].ns_unicasts = 0;
						clt_item[cn].ns_rx_rate_stats[j].ns_retries = 0;
					}
					if ((clt_item[cn].ns_tx_rate_stats[i].ns_rateKbps != 0) && (clt_item[cn].ns_tx_rate_stats[i].ns_rateKbps == clt_item[cn].ns_tx_rate_stats[j].ns_rateKbps)) {
						clt_item[cn].ns_tx_rate_stats[i].ns_unicasts += clt_item[cn].ns_tx_rate_stats[j].ns_unicasts;
						clt_item[cn].ns_tx_rate_stats[i].ns_retries += clt_item[cn].ns_tx_rate_stats[j].ns_retries;
						clt_item[cn].ns_tx_rate_stats[j].ns_rateKbps = 0;
						clt_item[cn].ns_tx_rate_stats[j].ns_unicasts = 0;
						clt_item[cn].ns_tx_rate_stats[j].ns_retries = 0;
					}
				}
			}

			/* Rate stat from DCD */
			for idx = 0; idx < NS_HW_RATE_SIZE; idx++ {
				if ( clt_item[cn].ns_tx_rate_stats[idx].ns_unicasts > t.last_clt_stat[ii][cn].ns_tx_rate_stats[idx].ns_unicasts) {
					tx_total += int64(clt_item[cn].ns_tx_rate_stats[idx].ns_unicasts - t.last_clt_stat[ii][cn].ns_tx_rate_stats[idx].ns_unicasts)
				}

				if (clt_item[cn].ns_rx_rate_stats[idx].ns_unicasts > t.last_clt_stat[ii][cn].ns_rx_rate_stats[idx].ns_unicasts) {
					rx_total += int64(clt_item[cn].ns_rx_rate_stats[idx].ns_unicasts - t.last_clt_stat[ii][cn].ns_rx_rate_stats[idx].ns_unicasts)
				}

				if (clt_item[cn].ns_tx_rate_stats[idx].ns_retries > t.last_clt_stat[ii][cn].ns_tx_rate_stats[idx].ns_retries) {
					tot_tx_bitrate_retries += clt_item[cn].ns_tx_rate_stats[idx].ns_retries - t.last_clt_stat[ii][cn].ns_tx_rate_stats[idx].ns_retries;
				}

				if (clt_item[cn].ns_rx_rate_stats[idx].ns_retries > t.last_clt_stat[ii][cn].ns_rx_rate_stats[idx].ns_retries) {
					tot_rx_bitrate_retries += clt_item[cn].ns_rx_rate_stats[idx].ns_retries - t.last_clt_stat[ii][cn].ns_rx_rate_stats[idx].ns_retries;
				}
			}

			//tx_ok = tx_total;
			//rx_ok = rx_total;
			/* Tx/Rx bit rate distribution */
			for idx = 0; idx < NS_HW_RATE_SIZE; idx++ {
				tmp_count1 = int32(clt_item[cn].ns_tx_rate_stats[idx].ns_unicasts - t.last_clt_stat[ii][cn].ns_tx_rate_stats[idx].ns_unicasts)
				if (tx_total > 0 && tmp_count1 > 0) {
					rf_report.tx_bit_rate[idx].rate_dtn = uint8((int64(tmp_count1) * 100) / tx_total)
				} else {
					rf_report.tx_bit_rate[idx].rate_dtn = 0;
				}

				tmp_count2 = int32(clt_item[cn].ns_rx_rate_stats[idx].ns_unicasts - t.last_clt_stat[ii][cn].ns_rx_rate_stats[idx].ns_unicasts)
				if (rx_total > 0 && tmp_count2 > 0) {
					rf_report.rx_bit_rate[idx].rate_dtn = uint8((int64(tmp_count2) * 100) / rx_total)
				} else {
					rf_report.rx_bit_rate[idx].rate_dtn = 0;
				}

				/* Tx/Rx bit rate success distribution */
				tmp_count3 = uint32(clt_item[cn].ns_tx_rate_stats[idx].ns_retries - t.last_clt_stat[ii][cn].ns_tx_rate_stats[idx].ns_retries)
				tmp_count5 = uint64(tmp_count1) + uint64(tmp_count3)
				if (tmp_count5 > 0 && rf_report.tx_bit_rate[idx].rate_dtn > 0) {
					rf_report.tx_bit_rate[idx].rate_suc_dtn = uint8((uint64(tmp_count1) * 100) / tmp_count5)
					if (rf_report.tx_bit_rate[idx].rate_suc_dtn > 100) {
						rf_report.tx_bit_rate[idx].rate_suc_dtn = 100
						log.Printf("DCD stats report client data process: rate_suc_dtn1 is more than 100%\n")
					}
				} else {
					rf_report.tx_bit_rate[idx].rate_suc_dtn = 0;
				}
				tx_retries += tmp_count3;

				tmp_count4 = clt_item[cn].ns_rx_rate_stats[idx].ns_retries - t.last_clt_stat[ii][cn].ns_rx_rate_stats[idx].ns_retries
				tmp_count6 = uint64(tmp_count2) + uint64(tmp_count4)
				if (tmp_count6 > 0 && rf_report.rx_bit_rate[idx].rate_dtn > 0) {
					rf_report.rx_bit_rate[idx].rate_suc_dtn = uint8((uint64(tmp_count2) * 100) / tmp_count6)
					if (rf_report.rx_bit_rate[idx].rate_suc_dtn > 100) {
						rf_report.rx_bit_rate[idx].rate_suc_dtn = 100
						log.Printf("DCD stats report client data process: rate_suc_dtn2 is more than 100%\n")
					}
				} else {
					rf_report.rx_bit_rate[idx].rate_suc_dtn = 0
				}

				rf_report.tx_bit_rate[idx].kbps = clt_item[cn].ns_tx_rate_stats[idx].ns_rateKbps
				rf_report.rx_bit_rate[idx].kbps = clt_item[cn].ns_rx_rate_stats[idx].ns_rateKbps
			}
			/* Rate stat from DCD */

			t.last_clt_stat[ii][cn] = clt_item[cn]
			ii++


			fields2["ifname"]               = intfName2
			fields2["ifIndex"]              = ifindex2

			fields2["mac_keys"]		= client_mac

			fields2["number"]		= cltstat.count
			fields2["ssid"]			= client_ssid
                        fields2["txPackets"]		= stainfo.tx_pkts
                        fields2["txBytes"]		= stainfo.tx_bytes
                        fields2["txDrop"]		= clt_item[cn].ns_tx_drops
                        fields2["slaDrop"]		= clt_item[cn].ns_sla_traps
                        fields2["rxPackets"]		= stainfo.rx_pkts
                        fields2["rxBytes"]		= stainfo.rx_bytes
                        fields2["rxDrop"]		= clt_item[cn].ns_tx_drops
                        fields2["avgSnr"]		= clt_item[cn].ns_snr
                        fields2["psTimes"]		= clt_item[cn].ns_ps_times
                        fields2["radioScore"]		= 0					/* TBD (Needs to implent calculations)    */
                        fields2["ipNetScore"]		= ipnet_score
			if ipnet_score == 0 {
				fields2["appScore"]	= ipnet_score
			} else {
				fields2["appScore"]	= clt_item[cn].ns_app_health_score
			}
                        fields2["phyMode"]		= getMacProtoMode(onesta.isi_phymode)
			if(stainfo != nil) {
				fields2["rssi"]		= int(stainfo.rssi) + int(stainfo.noise_floor)
			} else {
				fields2["rssi"]		= 0
			}
                        fields2["os"]			= ""					/* Not implemented in DCD also */
			fields2["name"]			= string(onesta.isi_name[:])
                        fields2["host"]			= ""
                        fields2["profName"]		= "default-profile"			/* TBD (Needs shared memory of dcd)	*/
                        fields2["dhcpIp"]		= sta_ip.dhcp_server
			fields2["gwIp"]			= sta_ip.gateway
                        fields2["dnsIp"]		= sta_ip.dns[0].dns_ip
			fields2["clientIp"]		= sta_ip.client_static_ip
                        fields2["dhcpTime"]		= sta_ip.dhcp_time
                        fields2["gwTime"]		= 0					/* TBD (Needs shared memory of auth2)     */
                        fields2["dnsTime"]		= sta_ip.dns[0].dns_response_time
                        fields2["clientTime"]		= onesta.isi_assoc_time


			for i := 0; i < AH_TX_NSS_MAX; i++{
				txNssUsage := fmt.Sprintf("txNssUsage_%d",i)
				fields2[txNssUsage]           = clt_item[cn].ns_tx_nss[i]
			}



			for i := 0; i < NS_HW_RATE_SIZE; i++{
				kbps := fmt.Sprintf("kbps_%d_rxRateStats",i)
				rateDtn := fmt.Sprintf("rateDtn_%d_rxRateStats",i)
				rateSucDtn := fmt.Sprintf("rateSucDtn_%d_rxRateStats",i)
				fields2[kbps]			= rf_report.rx_bit_rate[i].kbps
				fields2[rateDtn]		= rf_report.rx_bit_rate[i].rate_dtn
				fields2[rateSucDtn]		= rf_report.rx_bit_rate[i].rate_suc_dtn
			}


            for i := 0; i < NS_HW_RATE_SIZE; i++{
                kbps := fmt.Sprintf("kbps_%d_txRateStats",i)
		rateDtn := fmt.Sprintf("rateDtn_%d_txRateStats",i)
                rateSucDtn := fmt.Sprintf("rateSucDtn_%d_txRateStats",i)
                fields2[kbps]			= rf_report.tx_bit_rate[i].kbps
                fields2[rateDtn]		= rf_report.tx_bit_rate[i].rate_dtn
                fields2[rateSucDtn]		= rf_report.tx_bit_rate[i].rate_suc_dtn
            }

			if (clt_last_stats != nil) {
				fields2["txAirtime_min"]				= clt_last_stats.tx_airtime_min 
				fields2["txAirtime_max"]				= clt_last_stats.tx_airtime_max 
				fields2["txAirtime_avg"]				= clt_last_stats.tx_airtime_average

				fields2["rxAirtime_min"]				= clt_last_stats.rx_airtime_min
				fields2["rxAirtime_max"]				= clt_last_stats.rx_airtime_max
				fields2["rxAirtime_avg"]				= clt_last_stats.rx_airtime_average

				fields2["bwUsage_min"]					= clt_last_stats.bw_usage_min
				fields2["bwUsage_max"]					= clt_last_stats.bw_usage_max
				fields2["bwUsage_avg"]					= clt_last_stats.bw_usage_average
			}
		}
		acc.AddFields("ClientStats", fields2, tags, time.Now())

	}
	t.numclient =  total_client_count
	return nil
}

func Gather_AirTime(t *Ah_wireless, acc telegraf.Accumulator) error {

	for _, intfName2 := range t.Ifname {

//		var ifindex2 int
		var numassoc1 int
//		var stainfo *ah_ieee80211_sta_info
		var client_mac1 string
		var cintfName string
		var client_ssid string

		numassoc1 = int(getNumAssocs(t.fd, intfName2))

		if(numassoc1 == 0) {
			continue
		}

		//total_client_count = total_client_count + numassoc

		clt_item := make([]ah_ieee80211_sta_stats_item, numassoc1)
		var cltstat *ah_ieee80211_get_wifi_sta_stats
		cltstat = getStaStat(t.fd, intfName2, unsafe.Pointer(&clt_item[0]),  numassoc1)

		for cn := 0; cn < numassoc1; cn++ {
		//	if ( clt_item[cn] == nil) {
		//		continue
		//	}
			client_ssid = string(bytes.Trim(clt_item[cn].ns_ssid[:], "\x00"))

			if(clt_item[cn].ns_mac[0] !=0 || clt_item[cn].ns_mac[1] !=0 || clt_item[cn].ns_mac[2] !=0 || clt_item[cn].ns_mac[3] !=0 || clt_item[cn].ns_mac[4] != 0 || clt_item[cn].ns_mac[5]!=0) {
				cintfName = t.intf_m[intfName2][client_ssid]
				client_mac1 = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",clt_item[cn].ns_mac[0],clt_item[cn].ns_mac[1],clt_item[cn].ns_mac[2],clt_item[cn].ns_mac[3],clt_item[cn].ns_mac[4],clt_item[cn].ns_mac[5])
			} else {
				//stainfo = nil
				continue
			}

			log.Printf("AirTime %s %s %s %d",cintfName,client_mac1,client_ssid, cltstat.count)
			//TBD
			var clt_last_stats *saved_stats = (*saved_stats)(t.entity[intfName2][client_mac1])

			if(clt_last_stats == nil) {
				log.Printf("AirTime clt_last_stats is nil")
				clt_new_stats := saved_stats{
								tx_airtime_min:0,
								tx_airtime_max:0,
								tx_airtime_average:0,
								rx_airtime_min:0,
								rx_airtime_max:0,
								rx_airtime_average:0,
								bw_usage_min:0,
								bw_usage_max:0,
								bw_usage_average:0,
								tx_airtime:0,
								rx_airtime:0}
				t.entity[intfName2][client_mac1] = unsafe.Pointer(&clt_new_stats)
				return nil
			}

			clt_new_stats := saved_stats{}
			clt_new_stats.tx_airtime = clt_item[cn].ns_tx_airtime
			clt_new_stats.rx_airtime = clt_item[cn].ns_rx_airtime

			/* Calculate tx airtime min, max, average */

			if ((clt_last_stats.tx_airtime_min > clt_item[cn].ns_tx_airtime) || (clt_last_stats.tx_airtime_min == 0) ) {
				clt_new_stats.tx_airtime_min = clt_item[cn].ns_tx_airtime - clt_last_stats.tx_airtime
			}

			if (clt_last_stats.tx_airtime_max < clt_item[cn].ns_tx_airtime ) {
				clt_new_stats.tx_airtime_max = clt_item[cn].ns_tx_airtime - clt_last_stats.tx_airtime
			}

			clt_new_stats.tx_airtime_average = ((clt_last_stats.tx_airtime_average + clt_new_stats.tx_airtime_min + clt_new_stats.tx_airtime_max)/3)

			/* Calculate rx airtime min, max, average */

			if ((clt_last_stats.rx_airtime_min > clt_item[cn].ns_rx_airtime) || (clt_last_stats.rx_airtime_min == 0) ) {
				clt_new_stats.rx_airtime_min = clt_item[cn].ns_rx_airtime - clt_last_stats.rx_airtime
			}

			if (clt_last_stats.rx_airtime_max < clt_item[cn].ns_rx_airtime ) {
				clt_new_stats.rx_airtime_max = clt_item[cn].ns_rx_airtime - clt_last_stats.rx_airtime
			}

			clt_new_stats.rx_airtime_average = ((clt_last_stats.rx_airtime_average + clt_new_stats.rx_airtime_min + clt_new_stats.rx_airtime_max)/3)

			/* Calculate bandwidth usage min, max, average */

			bw_usage := (((clt_item[cn].ns_tx_bytes + clt_item[cn].ns_rx_bytes) * 8) / (60)) / 1000;

			if ((clt_last_stats.bw_usage_min > bw_usage) || (clt_last_stats.bw_usage_min == 0)) {
				clt_new_stats.bw_usage_min = bw_usage
			}

			if (clt_last_stats.bw_usage_max < bw_usage) {
				clt_new_stats.bw_usage_max = bw_usage
			}

			clt_new_stats.bw_usage_average = ((clt_last_stats.bw_usage_average + clt_new_stats.bw_usage_min + clt_new_stats.bw_usage_max)/3)

			//clt_last_stats := saved_stats{}

			//log.Printf("String entity[%s][%s]=%d %d %d",intfName2,client_mac1,clt_new_stats.tx_airtime_min,clt_new_stats.tx_airtime_max,clt_new_stats.tx_airtime_average)
			t.entity[intfName2][client_mac1] = unsafe.Pointer(&clt_new_stats)
		}

	}




	return nil
}

func (t *Ah_wireless) Gather(acc telegraf.Accumulator) error {	
	if t.timer_count == 9 {
		for _, intfName := range t.Ifname {
			t.intf_m[intfName] = make(map[string]string)
			load_ssid(t, intfName)
		}
		Gather_Rf_Stat(t, acc)
		Gather_Client_Stat(t, acc)
		t.timer_count = 0
	} else {
		Gather_AirTime(t,acc)
		t.timer_count++
	}

	return nil
}



func (t *Ah_wireless) Start(acc telegraf.Accumulator) error {
	t.intf_m = make(map[string]map[string]string)
	t.entity = make(map[string]map[string]unsafe.Pointer)

	for _, intfName := range t.Ifname {
		t.entity[intfName] = make(map[string]unsafe.Pointer)
//		load_ssid(t, intfName)
	}
	return nil
}

func init_fe() *os.File {
        file, err := os.Open(AH_FE_DEV_NAME)
        if err != nil {
                log.Printf("Error opening file:", err)
                return nil
        }
        //defer file.Close()


        // Get the current flags
        flags, err := unix.FcntlInt(file.Fd(), syscall.F_GETFD, 0)
        if err != nil {
                log.Printf("Error getting flags:", err)
                return nil
        }

        flags |= unix.FD_CLOEXEC

        // Set the close-on-exec flag
        _, err = unix.FcntlInt(file.Fd(), syscall.F_SETFD, flags)
        if err != nil {
                log.Printf("Error setting flags:", err)
                return nil
        }
	return file
}


func (t *Ah_wireless) Stop() {
	unix.Close(t.fd)
}


func init() {
	inputs.Add("ah_wireless", func() telegraf.Input {
		return NewAh_wireless(1)
	})
}
