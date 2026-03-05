package ah_trap

import (
	"net"
	"unsafe"
)

type AhTrapType uint32

const (
	AH_MAX_TRAP_OBJ_NAME     = 64
	AH_MAX_TRAP_IF_NAME      = 16
	AH_MAX_TRAP_SSID_NAME    = 32
	AH_MAX_TRAP_HOST_NAME    = 32
	AH_MAX_TRAP_USER_NAME    = 128
	AH_MAX_TRAP_PROF_NAME    = 32
	AH_MAX_NAME_LEN          = 32
	AH_UCHAR_MAX             = 255
	AH_MAX_NUM_STA_ADDRS6    = 5
	AH_TRAP_MSG_TYPE         = 1
	AH_FA_MVLAN_TRAP_TYPE    = 124
	TRAP_DCRPT_LEN           = 96
	AH_MSG_TRAP_DFS_BANG     = 12
	AH_MSG_TRAP_STA_LEAVE_STATS = 6
	MACADDR_LEN              = 6
	MAX_DESCRIBLE_LEN        = 128
	AH_CAPWAP_STAT_NAME_MAX_LEN = 32
	MAX_OBJ_NAME_LEN         = 4
	AH_MSG_TRAP_SSID_BIND_UNBIND = 5
	AH_MSG_TRAP_BSSID_SPOOFING = 7
	AH_MSG_TRAP_TB           = 2
	AH_TRAP_SIZE_300	  = 300
	AH_TRAP_SIZE_256         = 256
	AH_MSG_TRAP_DEV_IP_CHANGE = 17
	AH_MSG_TRAP_CLT_CAPS = 119
	AH_MGT0_ADDR6_NUM_MAX    = 2
	AH_SNMP_TRUE             = 1
	AH_SNMP_FALSE            = 2
	AH_MSG_TRAP_SET          = 0
	AH_MSG_TRAP_CLEAR        = 1
	AH_MSG_TRAP_VERIFY_OOB_SN  = 117
	AH_MSG_TRAP_PORTAL_CHANGE  = 120
	AH_TRAP_CLT_CAPS_MAX_STR_LEN  = 12
	AH_TRAP_CLT_CAPS_MIN_STR_LEN  = 8
	AH_MSG_TRAP_VPN          = 4
	AH_MSG_TRAP_CAPTURE_WARN = 127
	AH_MSG_TRAP_CAPWAP_DELAY = 10
	AH_CAPWAP_DELAY_TRAP     = 108
	MAX_CAPTURE_FILE_NAME_LEN = 16
	MAX_CAPTURE_FILES         = 16
	AH_MSG_TRAP_FA_ASSGN_MAP_CHANGE = 126
	AH_TELEGRAF_FA_MAP_MAX_ENTRIES = 94
	AH_TELEGRAF_SELF_REG_MAX_LEN  = 129
	AH_TELEGRAF_SELF_REG_USER_LEN = 257
	AH_MSG_TRAP_SELF_REG_INFO     = 11
	AH_MSG_TRAP_REPORT_CWP_INFO   = 15
	AH_TELEGRAF_CWP_FIELD_NAME_LEN = 50
	AH_TELEGRAF_CWP_FIELD_LEN      = 65
	AH_TELEGRAF_CWP_MAX_FIELDS     = 8
	AH_TELEGRAF_CWP_MAX_OPT_FIELDS = 8

)

const (
	AH_FAILURE_TRAP_TYPE AhTrapType = iota + 1
	AH_THRESHOLD_TRAP_TYPE
	AH_STATE_CHANGE_TRAP_TYPE
	AH_CONNECTION_CHANGE_TRAP_TYPE
	AH_IDP_AP_EVENT_TRAP_TYPE
	AH_CLIENT_INFO_TRAP_TYPE
	AH_POWER_INFO_TRAP_TYPE
	AH_CHANNEL_POWER_TRAP_TYPE
	AH_IDP_MITIGATE_TRAP_TYPE
	AH_INTERFERENCE_ALERT_TRAP_TYPE
	AH_BW_SENTINEL_TRAP_TYPE
	AH_ALARM_ALRT_TRAP_TYPE
	AH_MESH_MGT0_VLAN_CHANGE_TRAP_TYPE
	AH_KEY_FULL_ALARM_TRAP_TYPE
	AH_MESH_STABLE_STAGE_TRAP_TYPE
	AH_CHAIN_STREAM_TRAP_TYPE
)

func vpnPhaseToString(phase int32) string {
	switch phase {
	case 1:
		return "PHASE1"
	case 2:
		return "PHASE 2"
	default:
		return "UNKNOWN"
	}
}

func vpnStatusToString(status int32) string {
	switch status {
	case 1:
		return "UP"
	case 2:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

func severityToString(level int32) string {
	switch level {
	case 1:
		return "CLEAR"
	case 2:
		return "INFO"
	case 3:
		return "MINOR"
	case 4:
		return "MAJOR"
	case 5:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func clientAuthMethodToString(method int32) string {
	switch method {
	case 0:
		return "CWP"
	case 1:
		return "OPEN"
	case 2:
		return "WEP_OPEN"
	case 3:
		return "WEP_SHARED"
	case 4:
		return "WPA_PSK"
	case 5:
		return "WPA2_PSK"
	case 6:
		return "WPA_802.1X"
	case 7:
		return "WPA2_802.1X"
	case 8:
		return "WPA_AUTO_PSK"
	case 9:
		return "WPA_AUTO_802.1X"
	case 10:
		return "DYNAMIC_WEP"
	case 11:
		return "802.1X"
	case 12:
		return "WPA3_SAE"
	case 13:
		return "WPA3_802.1X"
	case 14:
		return "WPA3_SAE_EXTENDED"
	case 15:
		return "OWE"
	case 16:
		return "WPA3_802.1X_SUITE_B_192"
	default:
		return "UNKNOWN"
	}
}

func clientEncryptMethodToString(method int32) string {
	switch method {
	case 0:
		return "AES"
	case 1:
		return "TKIP"
	case 2:
		return "WEP"
	case 3:
		return "NONE"
	case 4:
		return "AES_GCMP"
	default:
		return "UNKNOWN"
	}
}

func clientMacProtoToString(proto int32) string {
	switch proto {
	case 0:
		return "802.11A"
	case 1:
		return "802.11B"
	case 2:
		return "802.11G"
	case 3:
		return "802.11N_5GHZ"
	case 4:
		return "802.11N_2.4GHZ"
	case 5:
		return "802.11AC_WIFI5"
	case 6:
		return "802.11AX_2.4GHZ_WIFI6"
	case 7:
		return "802.11AX_5GHZ_WIFI6"
	case 8:
		return "ETHERNET"
	case 9:
		return "802.11AX_6GHZ_WIFI6E"
	case 10:
		return "802.11BE_2.4GHZ_WIFI7"
	case 11:
		return "802.11BE_5GHZ_WIFI7"
	case 12:
		return "802.11BE_6GHZ_WIFI7"
	default:
		return "UNKNOWN"
	}
}
func MapStateToString(state uint8) string {
	switch state {
	case 1:
		return "PENDING"
	case 2:
		return "ACCEPT"
	case 3:
		return "REJECT"
	default:
		return "UNKNOWN"
	}
}

var eventReasonCodeNames = []string{
	"IDLE_TIMEOUT",
	"SESSION_TIMEOUT",
	"ADMIN_DEAUTH",
	"ASSOC_FAILED",
	"AUTH_FAILED",
	"AUTHZ_FAILED",
	"ROAMING",
	"SESSION_TIMEOUT_2",
	"AP_INITIATED_DISCONNECT",
	"CLIENT_INITIATED_DISCONNECT",
}

func eventReasonCodeToString(code uint32) string {
	if int(code) < len(eventReasonCodeNames) {
		return eventReasonCodeNames[code]
	}
	return "UNKNOWN"
}

var eventTypeNames = []string{
	"CLIENT_CONNECT",
	"CLIENT_DISCONNECT",
	"CLIENT_ROAM",
	"CLIENT_FAILURE",
}

func eventTypeToString(eventType uint32) string {
	if int(eventType) < len(eventTypeNames) {
		return eventTypeNames[eventType]
	}
	return "UNKNOWN"
}

var reasonCodeNames = []string{
	"RESERVED",
	"UNSPECIFIED",
	"PREV_AUTH_NOT_VALID",
	"STA_LEAVING_IBSS_ESS",
	"INACTIVITY",
	"AP_UNABLE_TO_HANDLE",
	"CLASS2_FRAME_FROM_NONAUTH_STA",
	"CLASS3_FRAME_FROM_NONASSOC_STA",
	"STA_LEAVING_BSS",
	"STA_NOT_AUTHENTICATED",
	"POWER_CAPABILITY_UNACCEPTABLE",
	"SUPPORTED_CHANNELS_UNACCEPTABLE",
	"RESERVED_12",
	"INVALID_IE",
	"MIC_FAILURE",
	"FOURWAY_HANDSHAKE_TIMEOUT",
	"GROUP_KEY_HANDSHAKE_TIMEOUT",
	"IE_IN_FOURWAY_DIFFERENT",
	"INVALID_GROUP_CIPHER",
	"INVALID_PAIRWISE_CIPHER",
	"INVALID_AKMP",
	"UNSUPPORTED_RSN_IE_VERSION",
	"INVALID_RSN_IE_CAPABILITIES",
	"IEEE_802_1X_AUTH_FAILED",
	"CIPHER_SUITE_REJECTED",
}

func reasonCodeToString(code uint32) string {
	if int(code) < len(reasonCodeNames) {
		return reasonCodeNames[code]
	}
	return "UNKNOWN"
}

type AhFaMvlanChangeTrap struct {
	TrapType      uint8
	SystemID      [10]uint8
	NativeTagged uint8
	MgmtVlan      uint16
	NativeVlan    uint16
}

type AhTgrafDfsTrap struct {
	TrapType  uint8
	TrapId    uint8
	IfName    [AH_MAX_TRAP_IF_NAME + 1]byte
	Desc      [TRAP_DCRPT_LEN]byte
}

type AhTgrafTbTrap struct {
	TrapId       uint8
	WarningLevel uint8
	Validity     uint8
	Desc         [MAX_DESCRIBLE_LEN]byte
}

type AhTgrafVpnTrap struct {
	TrapId    uint8
	Phase     uint8
	Status    uint8
	Objn      [AH_MAX_TRAP_IF_NAME]byte
	LocalIp   [AH_MAX_TRAP_HOST_NAME + 1]byte
	RemoteIp  [AH_MAX_TRAP_HOST_NAME + 1]byte
	Desc      [MAX_DESCRIBLE_LEN]byte
}

type AhCaptureFileInfo struct {
    FileName [MAX_CAPTURE_FILE_NAME_LEN]byte
    FileSize uint64
}

type AhTgrafCaptureWarnTrap struct {
    TrapId    uint8
    Desc      [MAX_DESCRIBLE_LEN]byte
    _         [7]byte
    TotalSize uint64
    FileCount uint8
    _        [7]byte
    Files     [MAX_CAPTURE_FILES]AhCaptureFileInfo
}

type AhTgrafFaAssignMapData struct {
       Isid  uint32
       Vlan  uint16
       State uint8
       _     [1]byte
}

type AhTgrafFaAssignMapChangeTrap struct {
       Ifindex int32
       TrapId  uint8
       Count   uint8
       _       [2]byte
       Data    [AH_TELEGRAF_FA_MAP_MAX_ENTRIES]AhTgrafFaAssignMapData
}

type AhTgrafCapwapDelayTrap struct {
	AvgDelay       uint64
	CurDelay       uint64
	MinorThreshold uint64
	MajorThreshold uint64
	Severity       [AH_MAX_TRAP_IF_NAME]byte
	Desc           [MAX_DESCRIBLE_LEN]byte
	TrapId         uint8
	Clear          uint8
	_             [4]byte
}

type AhTgrafSsidBindUnbindTrap struct {
	TrapType    uint8
	TrapID      uint8
	IfName      [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex     int32
	Description [TRAP_DCRPT_LEN]byte
	BssidMAC    [MACADDR_LEN]byte
	SSID        [AH_MAX_TRAP_SSID_NAME + 1]byte
	State       uint8
}

type AhTgrafCwpSelfRegInfoTrap struct {
        TrapType    uint8
        StaMAC      [MACADDR_LEN]byte
        ExpireTime  uint32
        UserName    [AH_TELEGRAF_SELF_REG_USER_LEN]byte
        Email       [AH_TELEGRAF_SELF_REG_MAX_LEN]byte
        CompanyName [AH_TELEGRAF_SELF_REG_MAX_LEN]byte
        CwpSsid     [AH_TELEGRAF_SELF_REG_MAX_LEN]byte
}

type AhTgrafCwpField struct {
        Name  [AH_TELEGRAF_CWP_FIELD_NAME_LEN]byte
        Value [AH_TELEGRAF_CWP_FIELD_LEN]byte
}

type AhTgrafCwpInfoTrap struct {
        TrapType      uint8
        MacAddr       [MACADDR_LEN]byte
        ObjName       [6]byte
        FieldCount    uint8
        OptFieldCount uint8
        Fields        [AH_TELEGRAF_CWP_MAX_FIELDS]AhTgrafCwpField
        OptFields     [AH_TELEGRAF_CWP_MAX_OPT_FIELDS]AhTgrafCwpField
}

type AhTgrafBSSIDSpoofingTrap struct {
    TrapID       uint8
    IfName      [AH_MAX_TRAP_OBJ_NAME + 1]byte
    Description  [MAX_DESCRIBLE_LEN]byte
    IfIndex      uint32
    BssidMAC    [MACADDR_LEN]byte
    AttackMAC    [MACADDR_LEN]byte
    AttackCount  uint32
    Protocol     uint16
    Severity     uint8
    SourceIP     uint32
    TargetIP     uint32
}

type AhTgrafDevIpChangeIpv6Data struct {
	Ipv6AddrType       uint8
	_                  [3]byte
	Ipv6Addr           [16]byte
	Ipv6Prefix         uint32
	Ipv6DefaultGateway [16]byte
}

type AhTgrafDevIpChangeTrap struct {
	TrapType           uint8
	_                  [3]byte
	Ipv4Addr           uint32
	Ipv4Netmask        uint32
	Ipv4DefaultGateway uint32
	Ipv6AddrNum        uint8
	_                  [3]byte
	Ipv6Data           [AH_MGT0_ADDR6_NUM_MAX]AhTgrafDevIpChangeIpv6Data
}

type AhVerifyOobSnTrap struct {
	TrapType	uint8;
	SerialNumber [AH_MAX_NAME_LEN]byte;
}

type AhPortalChangeTrap struct {
	TrapType	uint8;
	Macaddr [MACADDR_LEN]byte;
}
type AhTelegrafCltCapsTrap struct {
	TrapType    uint8
	CltMac      [MACADDR_LEN]byte
	BssidMac    [MACADDR_LEN]byte
	Channel     uint8
	Type        [AH_TRAP_CLT_CAPS_MAX_STR_LEN]byte
	Bw          [AH_TRAP_CLT_CAPS_MAX_STR_LEN]byte
	Nss         uint8
	Mode        [AH_TRAP_CLT_CAPS_MAX_STR_LEN]byte
	MinTxPower  int8
	MaxTxPower  int8
	MuMimo      [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
	Wmm         [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
	Cipher      [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
	Akm         [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
	Mfp         [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
	Mobile      [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
	Uapsd       [AH_TRAP_CLT_CAPS_MIN_STR_LEN]byte
}

type AhFailureTrap struct {
	Name  [AH_MAX_TRAP_OBJ_NAME+1]byte
	Cause int32
	Set   int32
}

type AhThresholdTrap struct {
	Name           [AH_MAX_TRAP_OBJ_NAME+1]byte
	CurVal         int32
	ThresholdHigh  int32
	ThresholdLow   int32
}

type AhStateChangeTrap struct {
	Name          [AH_MAX_TRAP_OBJ_NAME+1]byte
	PreState      int32
	CurState      int32
	OperationMode int32
}

type AhIdpApEventTrap struct {
	Name           [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex        int32
	RemoteID       [6]byte
	IdpType        int32
	IdpChannel     int32
	IdpRSSI        int32
	IdpCompliance  int32
	SSID           [AH_MAX_TRAP_SSID_NAME + 1]byte
	StationType    int32
	StationData    int32
	IdpRemoved     int32
	IdpInNet       int32
}

type AhClientInfoTrap struct {
	Name         [AH_MAX_TRAP_OBJ_NAME + 1]byte
	Ssid         [AH_MAX_TRAP_SSID_NAME + 1]byte
	ClientMac    [6]byte
	HostName     [AH_MAX_TRAP_HOST_NAME + 1]byte
	UserName     [AH_MAX_TRAP_USER_NAME + 1]byte
	ClientIP     uint32
	MgtStus      uint16
	StaAddr6Num  uint8
	StaAddr6     [AH_MAX_NUM_STA_ADDRS6][16]byte
}

type AhPowerInfoTrap struct {
	Name         [AH_MAX_TRAP_OBJ_NAME + 1]byte
	PowerSrc     int
	Eth0On       int
	Eth1On       int
	Eth0Pwr      int
	Eth1Pwr      int
	Eth0Speed    int
	Eth1Speed    int
	Wifi0Setting int
	Wifi1Setting int
	Wifi2Setting int
}

type AhChannelPowerTrap struct {
	Name           [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex        int32
	RadioChannel   int32
	RadioTxPower   int32
	BeaconInterval uint32
	ChnlStrfmt     uint16
	PwrStrfmt      uint16
	RadioEirp      [8]byte
	Reason         int32
}

type AhIdpMitigateTrap struct {
	Name         [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex      int32
	RemoteID     [6]byte
	BSSID        [6]byte
	Removed      int32
	DiscoverAge  uint32
	UpdateAge    uint32
}

type AhInterferenceAlertTrap struct {
	Name                [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex             int32
	InterferenceThres   int32
	AveInterference     int32
	ShortInterference   int32
	SnapInterference    int32
	CRCErrRateThres     int32
	CRCErrRate          int32
	Set                 int32
}

type AhBwSentinelTrap struct {
	Name              [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex           int32
	ClientMac         [6]byte
	BwSentinelStatus  int32
	GBW               int32
	ActualBW          int32
	ActionTaken       uint32
	ChnlUtil          uint8
	InterferenceUtil  uint8
	TxUtil            uint8
	RxUtil            uint8
}

type AhAlarmAlrtTrap struct {
	Name               [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex            int32
	ClientMac          [6]byte
	Level              int32
	SSID               [AH_MAX_TRAP_SSID_NAME + 1]byte
	AlertType          int32
	ThresInterference  int32
	ShortInterference  int32
	SnapInterference   int32
	Set                int32
}

type AhMeshMgt0VlanChangeTrap struct {
	Name           [AH_MAX_TRAP_OBJ_NAME + 1]byte
	OldVlan        uint16
	NewVlan        uint16
	OldNativeVlan  uint16
	NewNativeVlan  uint16
}

type AhMeshStableStageTrap struct {
	Name            [AH_MAX_TRAP_OBJ_NAME + 1]byte
	MeshStableStage int32
	MeshDataRate    int32
}

type AhKeyFullAlarmTrap struct {
	Name      [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex   int32
	BSSID     [MACADDR_LEN]byte
	ClientMAC [MACADDR_LEN]byte
	GtkVLAN   uint32
}

type AhChainStreamTrap struct {
	Name         [AH_MAX_TRAP_OBJ_NAME + 1]byte
	IfIndex      uint32
	NewTxChain   uint8
	NewTxStream  uint8
	NewRxChain   uint8
	NewRxStream  uint8
	Reason       uint8
}

// Helper functions
func intToIPv4(num uint32) string {
	ip := net.IPv4(
		byte(num>>24),
		byte(num>>16),
		byte(num>>8),
		byte(num),
	)
	return ip.String()
}

func intToIPv6(addrs [][16]byte, count int) string {
	if count > 0 && count <= len(addrs) {
		ip := net.IP(addrs[0][:])
		ipStr := ip.String()
		return ipStr
	}
	return ""
}

func IntToIPv6_1(addrs []AhStaAddr6, count int) []string {
	var result []string
	for i := 0; i < count && i < len(addrs); i++ {
		ip := net.IP(addrs[i].StaAddr6[:])
		result = append(result, ip.String())
	}
	return result
}

type AhStaAddr6 struct {
	AddrType byte      // char addr_type
	StaAddr6 [16]byte  // struct in6_addr (IPv6 address is 16 bytes)
}

// GetTrapClearStatus extracts clear status from trap based on trap type and set field
// Mimics ah_get_msg_trap_type_clear() from syslogd.c
// Returns: false = SET (alarm raised), true = CLEAR (alarm cleared)
func GetTrapClearStatus(trapType uint32, unionData []byte) bool {
	isClear := false // Default: SET

	switch AhTrapType(trapType) {
	case AH_FAILURE_TRAP_TYPE:
		var failure AhFailureTrap
		if len(unionData) >= int(unsafe.Sizeof(failure)) {
			copy((*[1 << 10]byte)(unsafe.Pointer(&failure))[:unsafe.Sizeof(failure)], unionData)
			if failure.Set == AH_SNMP_TRUE {
				isClear = false // SET
			} else {
				isClear = true // CLEAR
			}
		}

	case AH_INTERFERENCE_ALERT_TRAP_TYPE:
		var interference AhInterferenceAlertTrap
		if len(unionData) >= int(unsafe.Sizeof(interference)) {
			copy((*[1 << 10]byte)(unsafe.Pointer(&interference))[:unsafe.Sizeof(interference)], unionData)
			if interference.Set == AH_SNMP_TRUE {
				isClear = false // SET
			} else {
				isClear = true // CLEAR
			}
		}

	case AH_ALARM_ALRT_TRAP_TYPE:
		var alarmAlert AhAlarmAlrtTrap
		if len(unionData) >= int(unsafe.Sizeof(alarmAlert)) {
			copy((*[1 << 10]byte)(unsafe.Pointer(&alarmAlert))[:unsafe.Sizeof(alarmAlert)], unionData)
			if alarmAlert.Set == AH_SNMP_TRUE {
				isClear = false // SET
			} else {
				isClear = true // CLEAR
			}
		}

	case AH_CAPWAP_DELAY_TRAP:
		var capwapDelay AhTgrafCapwapDelayTrap
		if len(unionData) >= int(unsafe.Sizeof(capwapDelay)) {
			copy((*[1 << 10]byte)(unsafe.Pointer(&capwapDelay))[:unsafe.Sizeof(capwapDelay)], unionData)
			if capwapDelay.Clear != 0 {
				isClear = true // CLEAR
			} else {
				isClear = false // SET
			}
		}

	default:
		// Event traps always return SET (false)
		isClear = false
	}

	return isClear
}
