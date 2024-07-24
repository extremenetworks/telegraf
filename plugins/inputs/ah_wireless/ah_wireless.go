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
	fd     int
	fe_fd	uintptr
	intf_m  map[string]map[string]string
	arp_m	map[string]string
	Ifname	[]string	`toml:"ifname"`
	closed chan struct{}
	numclient int
}


func wl_ioctl(fd int, op, argp uintptr) error {

	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(op), argp)
	if errno != 0 {
		return errno
	}
	return nil
}

func fe_ioctl(fd uintptr, op, argp uintptr) error {
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
	log.Printf("Chayan NewAh_wireless")
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

        if err := wl_ioctl(fd, IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
                return nil
        }

	offsetsMutex.Unlock()

        return &cfg.hdd_stats

}

func getAtrTbl(fd int, ifname string) *ah_ieee80211_atr_user {
	log.Printf("Chayan getAtrTbl\n")
	var cfg *ieee80211req_cfg_atr
        cfg = new(ieee80211req_cfg_atr)

//	cfgss := ieee80211req_cfg_atr{}

	/* first 4 bytes is subcmd */
        cfg.cmd = AH_IEEE80211_GET_ATR_TBL;

	iwp := iw_point{pointer: unsafe.Pointer(cfg)}

	request := iwreq{data: iwp}

//        request.data.pointer = &cfg;
//        request.data.length = uint16(unsafe.Sizeof(cfgss));
	request.data.length = VAP_BUFF_SIZE

///	request := iwreq{data: uintptr(unsafe.Pointer(request.data.pointer))}

	copy(request.ifrn_name[:], ah_ifname_radio2vap(ifname))

//	var cfgs *ieee80211req_cfg_atr
//	var cfgs = (*ieee80211req_cfg_atr)(request.data.pointer)


//	log.Printf("Chayan getAtrTbl before ioctl cmd = %d\n", cfgs.cmd)
	offsetsMutex.Lock()

	if err := wl_ioctl(fd, IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
                return nil
        }

	offsetsMutex.Unlock()
//	log.Printf("Chayan getAtrTbl after ioctl\n")

	return &cfg.atr

}




func getRFStat(fd int, ifname string) *awestats {
//	bs := awestats{}

//	bs := awestats.new()
//          var p = &bs
    //   p = C.malloc(C.sizeof(C.struct_awestats))
//	p := unsafe.Pointer(&bs)


//      ifrd := ifr.withData(p)


	var p *awestats
	p = new(awestats)

//ethtool - https://github.com/weaveworks/weave/blob/master/net/ethtool.go#L51
        // Request current value/
	request := IFReqData{Data: uintptr(unsafe.Pointer(p))}
//	request := IFReqData{Data: uintptr(unsafe.Pointer(&bs))}
//	request := IFReqData{Data: uintptr(p)}
	copy(request.Name[:], ifname)

//	p.ast_noise_floor = 1

//        log.Printf("Chayan open before calling ioctl data - %d %d",int(bs.ast_as.ast_rx_bytes), int(bs.ast_noise_floor))

//ethtool
	offsetsMutex.Lock()

        if err := wl_ioctl(fd, SIOCGRADIOSTATS, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
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

        if err := wl_ioctl(fd, IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
                return nil
        }

        offsetsMutex.Unlock()

        return &cfg.wifi_sta_stats
}

func getNumAssocs (fd int, ifname string) uint32 {

        //var ird *iwreq_data
//	ird := new(iwreq_data)

//	ird := iwreq_data{data: IEEE80211_PARAM_NUM_ASSOCS}

	ird := iwreq_data{}

        /* first 4 bytes is subcmd */
        ird.data = IEEE80211_PARAM_NUM_ASSOCS

    //    iwp := iw_point{pointer: unsafe.Pointer(cfg)}


        request := iwreq_clt{u: ird}

    //    request.data.length = VAP_BUFF_SIZE

        copy(request.ifr_name[:], ah_ifname_radio2vap(ifname))

//		log.Printf("Chayan getNumAssocs %s",request.ifr_name[:])

        offsetsMutex.Lock()

        if err := wl_ioctl(fd, IEEE80211_IOCTL_GETPARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
                return 0
        }

//	log.Printf("Chayan getNumAssocs %d", uint32(request.u.data))

        offsetsMutex.Unlock()

//        return int(ird.data)
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
	log.Printf("Chayan ifname = %s",request.ifrn_name)
	log.Printf("Chayan One STA MAC %x:%x:%x:%x:%x:%x", cfg.sta_info.mac[0],cfg.sta_info.mac[1],cfg.sta_info.mac[2],cfg.sta_info.mac[3],cfg.sta_info.mac[4],cfg.sta_info.mac[5])

        offsetsMutex.Lock()
        if err := wl_ioctl(fd, IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
                return nil
        }

	log.Printf("Chayan One STA rssi = %d, noise_floor = %d",cfg.sta_info.rssi, cfg.sta_info.noise_floor)
        offsetsMutex.Unlock()

        return &cfg.sta_info

}

func getOneSta(fd int, ifname string, mac_ad [MACADDR_LEN]uint8) unsafe.Pointer {
//func getOneSta(fd int, ifname string, mac_ad [MACADDR_LEN]uint8) uintptr {
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
	if err := wl_ioctl(fd, IEEE80211_IOCTL_GENERIC_PARAM, uintptr(unsafe.Pointer(&request))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
		offsetsMutex.Unlock()
		return nil
        }
	offsetsMutex.Unlock()

        return request.data.pointer
//	return uintptr(request.data.pointer)

}

func getProcNetDev(ifname string) *ah_dcd_dev_stats {
	table, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil;
	}

	lines := bytes.Split([]byte(table), newLineByte)

	var intfname string
	var stats = new(ah_dcd_dev_stats)

  for  _, curLine := range lines {
    if strings.Contains(string(curLine), " wifi1:") {
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
  //      fmt.Println(stats.rx_packets)
        fmt.Println(n)
    }
  }
//	log.Printf(string(table))

	return stats
}

func getIfIndex(fd int, ifname string) int {
        ifr, err := unix.NewIfreq(ifname)

        if err != nil {
                log.Printf("failed to create ifreq  %v", err)
        }

//        log.Printf("Chayan open before calling ioctl index = %d",int(ifr.Uint32()))
	offsetsMutex.Lock()

        if err := unix.IoctlIfreq(fd, unix.SIOCGIFINDEX, ifr); err != nil {
                log.Printf("Chayan ioctl error %s",err)
                return -1
        }

	offsetsMutex.Unlock()
	return int(ifr.Uint32())
//        log.Printf("Chayan open after calling ioctl index = %d",int(ifr.Uint32()))

}

func load_ssid(t *Ah_wireless, ifname string) {
	log.Printf("load_ssid")

	for i := 1; i < 1024; i++ {
		vifname := ifname + "." + strconv.Itoa(i)
		log.Printf(vifname)

		app := "wl"

		arg0 := "-i"
		arg1 := vifname
		arg2 := "status"
		//  arg3 := "| grep \"SSID: \"\\\""
		log.Printf(app + " " + arg0 + " " + arg1 + " " + arg2)

		cmd := exec.Command(app, arg0, arg1, arg2)
		output, err := cmd.Output()

		if err != nil {
			log.Printf(err.Error())
			return
		}

//		log.Printf(string(output))

		lines := strings.Split(string(output),"\n")

//		log.Printf(lines[0])

		temp  := strings.Split(lines[0]," ")

//		log.Printf(temp[1])

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

//	log.Printf(string(arp_str))

	arp_lines := strings.Split(string(arp_str),"\n")

//	log.Printf("len = %d",len(arp_lines))

	for i :=0; i<len(arp_lines); i++ {
		log.Printf("Line[%d]:[%s]" ,i,arp_lines[i])
		if len(arp_lines[i]) > 1 {
			arp_eliments := strings.Split(arp_lines[i]," ")

//			for j :=0; j<len(arp_eliments); j++ {
//				log.Printf(arp_eliments[j])
			log.Printf("Adding: t.arp_m[%s] = %s",arp_eliments[3], arp_eliments[1])
			t.arp_m[arp_eliments[3]] = arp_eliments[1]
//			}
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

//	log.Printf("Chayan calling getFeIpnetScore fd = %d size = %d msg_type = %d ioctl = %d", fd, ihdr.msg_size, ihdr.msg_type, AH_FE_IOCTL_FLOW)

        offsetsMutex.Lock()

        if err := fe_ioctl(fd, AH_FE_IOCTL_FLOW, uintptr(unsafe.Pointer(&dev_msg))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
		offsetsMutex.Unlock()
                return -1
	}

	offsetsMutex.Unlock()

//	log.Printf("Chayan calling getFeIpnetScore  = %d",msg.net_health_score)
	log.Printf("Chayan calling getFeIpnetScore  = %d",dev_msg.data.net_health_score)

	if dev_msg.hdr.retval < 0 {
		log.Printf("Chayan open ioctl data erro")
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

        if err := fe_ioctl(fd, AH_FE_IOCTL_FLOW, uintptr(unsafe.Pointer(&dev_msg))); err != nil {
                log.Printf("Chayan open ioctl data error %s",err)
                offsetsMutex.Unlock()
                return nil
        }

        offsetsMutex.Unlock()


        if dev_msg.hdr.retval < 0 {
                log.Printf("Chayan open ioctl data erro")
                return nil
        }

        return &dev_msg.data
}

func open(fd, id int) *Ah_wireless {
	log.Printf("Chayan open")

	getProcNetDev("wifi1")
	//defer unix.Close(fd)

	return &Ah_wireless{fd: fd, closed: make(chan struct{})}
}

func (t *Ah_wireless) SampleConfig() string {
	log.Printf("Chayan SampleConfig")
	return sampleConfig
}

func (t *Ah_wireless) Description() string {
	return "Hive OS wireless stat"
}

func (t *Ah_wireless) Init() error {
	log.Printf("Chayan Init")

//	return init()
//	var err error
	return nil
}


func (t *Ah_wireless) Gather(acc telegraf.Accumulator) error {


//	log.Printf("Chayan Gather fd = %d interface = %s",t.fd,  t.Ifname[0])
//	t.arp_m = make(map[string]string)
	//load_arp_table(t)

	for _, intfName := range t.Ifname {

//	var rfstat1 awestats
//	var rfstat2 awestats
	var rfstat *awestats
	var devstats *ah_dcd_dev_stats
	var ifindex int
	var atrStat *ah_ieee80211_atr_user
	var hddStat *ah_ieee80211_hdd_stats


//	offsetsMutex.Lock()
//	getRFStat(sockfd, "wifi0")
//-----------------------------------------
	rfstat  = getRFStat(t.fd, intfName)
	ifindex = getIfIndex(t.fd, intfName)
	devstats = getProcNetDev(intfName)
	atrStat = getAtrTbl(t.fd, intfName)
	hddStat = getHDDStat(t.fd, intfName)
//-----------------------------------------
//	offsetsMutex.Unlock()

//	log.Printf("Chayan Gather wifi0 - %d %d",rfstat1.ast_as.ast_rx_bytes, rfstat1.ast_noise_floor)
//	log.Printf("Chayan Gather %s - %d : %d %d %d -- %d %d %d -- %d %d %d",intfName, ifindex, rfstat.ast_as.ast_rx_bytes, rfstat.ast_noise_floor, 
//		devstats.rx_packets, atrStat.atr_info[0].rxc_pcnt, atrStat.atr_info[0].rxf_pcnt, 
//		atrStat.atr_info[0].txf_pcnt, hddStat.bs_sp_cnt, hddStat.lb_sp_cnt, hddStat.snr_sp_cnt)
  //      log.Printf("Chayan Gather wifi0 - %d %d",int(rfstat1.ast_as.ast_rx_bytes), int(rfstat1.ast_noise_floor))
    //    log.Printf("Chayan Gather wifi1 - %d %d",int(rfstat2.ast_as.ast_rx_bytes), int(rfstat2.ast_noise_floor))

	fields := map[string]interface{}{

		"name_keys":					intfName,
		"ifindex_keys":					ifindex,

	}
		 fields["channelUtilization_min"] = 0
		 fields["channelUtilization_max"] = 0
		 fields["channelUtilization_avg"] = 0

		 fields["interferenceUtilization_min"]		= 0
		 fields["interferenceUtilization_max"]		= 0
		 fields["interferenceUtilization_avg"]		= 0

		 fields["txUtilization_min"]				= atrStat.atr_info[1].txf_pcnt
		fields["txUtilization_max"]				= atrStat.atr_info[1].txf_pcnt
		fields["txUtilization_avg"]				= atrStat.atr_info[1].txf_pcnt

		fields["rxUtilization_min"]				= atrStat.atr_info[1].rxf_pcnt
		fields["rxUtilization_max"]				= atrStat.atr_info[1].rxf_pcnt
		fields["rxUtilization_avg"]				= atrStat.atr_info[1].rxf_pcnt

		fields["rxInbssUtilization_min"]			= atrStat.atr_info[1].rxf_inbss
		fields["rxInbssUtilization_max"]			= atrStat.atr_info[1].rxf_inbss
		fields["rxInbssUtilization_avg"]			= atrStat.atr_info[1].rxf_inbss

		fields["rxObssUtilization_min"]			= atrStat.atr_info[1].rxf_obss
		fields["rxObssUtilization_max"]			= atrStat.atr_info[1].rxf_obss
		fields["rxObssUtilization_avg"]			= atrStat.atr_info[1].rxf_obss

		fields["wifinterferenceUtilization_min"]		= 0
		fields["wifinterferenceUtilization_max"]		= 0
		fields["wifinterferenceUtilization_avg"]		= 0

		fields["noise_min"]					= rfstat.ast_noise_floor
		fields["noise_max"]					= rfstat.ast_noise_floor
		fields["noise_avg"]					= rfstat.ast_noise_floor

		fields["crcErrorRate_min"]				= rfstat.phy_stats.ast_rx_crcerr
		fields["crcErrorRate_max"]				= rfstat.phy_stats.ast_rx_crcerr
		fields["crcErrorRate_avg"]				= rfstat.phy_stats.ast_rx_crcerr


		    fields["txPackets"]				= devstats.tx_packets
		fields["txErrors"]				= devstats.tx_errors
		fields["txDropped"]				= devstats.tx_dropped
		fields["txHwDropped"]				= rfstat.ast_as.ast_tx_shortpre + rfstat.ast_as.ast_tx_xretries + rfstat.ast_as.ast_tx_fifoerr
		fields["txSwDropped"]				= devstats.tx_dropped
		fields["txBytes"]				= devstats.tx_bytes
		fields["txRetryCount"]				= rfstat.phy_stats.ast_tx_shortretry + rfstat.phy_stats.ast_tx_longretry

		fields["txRate_min"]				= rfstat.ast_tx_rate_stats[0].ns_rateKbps
		fields["txRate_max"]				= rfstat.ast_tx_rate_stats[0].ns_rateKbps
		fields["txRate_avg"]				= rfstat.ast_tx_rate_stats[0].ns_rateKbps

		fields["txUnicastPackets"]			= rfstat.ast_tx_rate_stats[0].ns_unicasts
		fields["txMulticastPackets"]			= devstats.tx_multicast
		fields["txMulticastBytes"]			= rfstat.ast_as.ast_tx_mcast_bytes
		fields["txBcastBytes"]				= rfstat.ast_as.ast_tx_bcast_bytes
		fields["txBcastPackets"]			= devstats.tx_broadcast

		fields["rxPackets"]				= devstats.rx_packets
		fields["rxErrors"]				= devstats.rx_errors
		fields["rxDropped"]				= devstats.rx_dropped
		fields["rxBytes"]				= devstats.rx_bytes
		fields["rxRetryCount"]				= rfstat.ast_rx_retry

		fields["rxRate_min"]				= rfstat.ast_rx_rate_stats[0].ns_rateKbps
		fields["rxRate_max"]				= rfstat.ast_rx_rate_stats[0].ns_rateKbps
		fields["rxRate_avg"]				= rfstat.ast_rx_rate_stats[0].ns_rateKbps

		fields["rxMulticastBytes"]			= rfstat.ast_rx_mcast_bytes
		fields["rxMulticastPackets"]			= devstats.rx_multicast
		fields["rxBcastPackets"]			= devstats.rx_broadcast
		fields["rxBcastBytes"]				= rfstat.ast_rx_bcast_bytes

		fields["bsSpCnt"]				= hddStat.bs_sp_cnt
		fields["snrSpCnt"]				= hddStat.snr_sp_cnt
		fields["snAnswerCnt"]				= 0
		fields["rxPrbSpCnt"]				= 0
		fields["rxAuthCnt"]				= 0
		fields["txBitrateSuc"]				= rfstat.ast_tx_rix_invalids
		fields["rxBitrateSuc"]				= rfstat.ast_rx_rix_invalids

			for i := 0; i < NS_HW_RATE_SIZE; i++{
				kbps := fmt.Sprintf("kbps_%d_rxRateStats",i)
				rateDtn := fmt.Sprintf("rateDtn_%d_rxRateStats",i)
				rateSucDtn := fmt.Sprintf("rateSucDtn_%d_rxRateStats",i)
				fields[kbps]			= rfstat.ast_rx_rate_stats[i].ns_rateKbps
				fields[rateDtn]		= 0
				fields[rateSucDtn]		= 0
			}


			for i := 0; i < NS_HW_RATE_SIZE; i++{
				kbps := fmt.Sprintf("kbps_%d_txRateStats",i)
				rateDtn := fmt.Sprintf("rateDtn_%d_txRateStats",i)
				rateSucDtn := fmt.Sprintf("rateSucDtn_%d_txRateStats",i)
				fields[kbps]			= rfstat.ast_tx_rate_stats[i].ns_rateKbps
				fields[rateDtn]		= 0
				fields[rateSucDtn]		= 0
			}

		fields["clientCount"]				= t.numclient



		fields["lbSpCnt"]				= hddStat.lb_sp_cnt


		fields["rxProbeSup"]				= rfstat.is_rx_hdd_probe_sup

		fields["rxSwDropped"]				= devstats.rx_dropped
		fields["rxUnicastPackets"]			= rfstat.ast_rx_rate_stats[0].ns_unicasts


		acc.AddGauge("radioInterfaceRfStats", fields, nil)
	}

	tags := map[string]string{
	}


	fields2 := map[string]interface{}{
	}

	var total_client_count int
	total_client_count = 0

        for _, intfName2 := range t.Ifname {

		var cltstat *ah_ieee80211_get_wifi_sta_stats
		var ifindex2 int
		var numassoc int
		var stainfo *ah_ieee80211_sta_info

		//	var clt_item *ah_ieee80211_sta_stats_item
		//	clt_item = new(ah_ieee80211_sta_stats_item)

		numassoc = int(getNumAssocs(t.fd, intfName2))

		log.Printf("Chayan interface: %s numassoc: %d",intfName2, numassoc)

		if(numassoc == 0) {
			continue;
		}

		total_client_count = total_client_count + numassoc

		clt_item := make([]ah_ieee80211_sta_stats_item, numassoc)

		log.Printf("Chayan make clt_item")

		ifindex2 = getIfIndex(t.fd, intfName2)

		log.Printf("Chayan getIfIndex = %d",ifindex2)

		cltstat = getStaStat(t.fd, intfName2, unsafe.Pointer(&clt_item[0]),  numassoc)

		log.Printf("Chayan called getStaStat")

		for cn := 0; cn < numassoc; cn++ {
			client_ssid := string(bytes.Trim(clt_item[cn].ns_ssid[:], "\x00"))

			log.Printf("Chayan clt_item[%d] ssid = %s",cn,client_ssid)

			if(clt_item[cn].ns_mac[0] !=0 || clt_item[cn].ns_mac[1] !=0 || clt_item[cn].ns_mac[2] !=0 || clt_item[cn].ns_mac[3] !=0 || clt_item[cn].ns_mac[4] != 0 || clt_item[cn].ns_mac[5]!=0) {
				cintfName := t.intf_m[intfName2][client_ssid]
				stainfo = getOneStaInfo(t.fd, cintfName, clt_item[cn].ns_mac)

				if stainfo.rssi == 0 {
					continue
				}
			} else {
				stainfo = nil
				continue
			}

			f := init_fe()
			ipnet_score := getFeIpnetScore(f.Fd(), clt_item[cn].ns_mac)
			log.Printf("Chayan ipnet_score = %d",ipnet_score)
			sta_ip := getFeServerIp(f.Fd(), clt_item[cn].ns_mac)
			log.Printf("Chayan sta_ip = %d %d", sta_ip.client_static_ip, sta_ip.gateway)
			f.Close()

			client_mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",clt_item[cn].ns_mac[0],clt_item[cn].ns_mac[1],clt_item[cn].ns_mac[2],clt_item[cn].ns_mac[3],clt_item[cn].ns_mac[4],clt_item[cn].ns_mac[5])
			cfgptr := getOneSta(t.fd, intfName2, clt_item[cn].ns_mac)

			var onesta *ieee80211req_sta_info = (*ieee80211req_sta_info)(cfgptr)

			log.Printf("Chayan phymode = %d %s",onesta.isi_phymode, client_mac)

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
				fields2[kbps]			= clt_item[cn].ns_rx_rate_stats[i].ns_rateKbps
				fields2[rateDtn]		= 0
				fields2[rateSucDtn]		= 0
			}


                        for i := 0; i < NS_HW_RATE_SIZE; i++{
                                kbps := fmt.Sprintf("kbps_%d_txRateStats",i)
                                rateDtn := fmt.Sprintf("rateDtn_%d_txRateStats",i)
                                rateSucDtn := fmt.Sprintf("rateSucDtn_%d_txRateStats",i)
                                fields2[kbps]			= clt_item[cn].ns_tx_rate_stats[i].ns_rateKbps
                                fields2[rateDtn]		= 0
                                fields2[rateSucDtn]		= 0
                        }
		}
		acc.AddFields("radioInterfaceClientStats", fields2, tags, time.Now())

	}
	t.numclient =  total_client_count

	return nil
}



func (t *Ah_wireless) Start(acc telegraf.Accumulator) error {
	log.Printf("Chayan Start")
	t.intf_m = make(map[string]map[string]string)
//	t.arp_m = make(map[string]string)

	for _, intfName := range t.Ifname {
		t.intf_m[intfName] = make(map[string]string)
		load_ssid(t, intfName)
	}
//	load_arp_table(t)
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
	log.Printf("Chayan Stop")
	unix.Close(t.fd)
}


func init() {
	log.Printf("Chayan init")
	inputs.Add("ah_wireless", func() telegraf.Input {
		return NewAh_wireless(1)
	})
}
