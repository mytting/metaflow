package policy

import (
	"reflect"
	"testing"
	"time"

	. "github.com/google/gopacket/layers"

	. "gitlab.x.lan/yunshan/droplet-libs/datatype"
)

var (
	forward  = AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddDirections(FORWARD).AddTagTemplates(TEMPLATE_EDGE_PORT)
	backward = AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddDirections(BACKWARD).AddTagTemplates(TEMPLATE_EDGE_PORT)
)

func getBackwardAcl(acl AclAction) AclAction {
	return acl.SetDirections(BACKWARD)
}

func CheckPolicyResult(basicPolicy *PolicyData, targetPolicy *PolicyData) bool {
	if reflect.DeepEqual(basicPolicy, targetPolicy) {
		return true
	}

	return false
}

func TestGetPlatformData(t *testing.T) {

	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)

	srcIp := NewIPFromString("192.168.2.12")
	dstIp := NewIPFromString("192.168.0.11")
	key := &LookupKey{
		SrcIp:  srcIp.Int(),
		SrcMac: 0x80027a42bfc,
		DstMac: 0x80027a42bfa,
		DstIp:  dstIp.Int(),
		Tap:    TAP_TOR,
	}
	ip := NewIPFromString("192.168.0.11")
	ipInfo := IpNet{
		Ip:       ip.Int(),
		SubnetId: 121,
		Netmask:  24,
	}

	ip1 := NewIPFromString("192.168.0.12")
	ipInfo1 := IpNet{
		Ip:       ip1.Int(),
		SubnetId: 122,
		Netmask:  25,
	}

	mac := NewMACAddrFromString("08:00:27:a4:2b:fc")
	launchServer := NewIPFromString("10.10.10.10")
	vifData := PlatformData{
		EpcId:      11,
		DeviceType: 2,
		DeviceId:   3,
		IfType:     3,
		IfIndex:    5,
		Mac:        mac.Int(),
		HostIp:     launchServer.Int(),
	}

	vifData.Ips = append(vifData.Ips, &ipInfo)
	vifData.Ips = append(vifData.Ips, &ipInfo1)

	ip2 := NewIPFromString("192.168.2.0")
	ipInfo2 := IpNet{
		Ip:       ip2.Int(),
		SubnetId: 125,
		Netmask:  24,
	}

	ip3 := NewIPFromString("192.168.2.12")

	ipInfo3 := IpNet{
		Ip:       ip3.Int(),
		SubnetId: 126,
		Netmask:  32,
	}

	mac1 := NewMACAddrFromString("08:00:27:a4:2b:fa")
	launchserver1 := NewIPFromString("10.10.10.10")

	vifData1 := PlatformData{
		EpcId:      0,
		DeviceType: 1,
		DeviceId:   100,
		IfType:     3,
		IfIndex:    5,
		Mac:        mac1.Int(),
		HostIp:     launchserver1.Int(),
	}

	vifData1.Ips = append(vifData1.Ips, &ipInfo2)
	vifData1.Ips = append(vifData1.Ips, &ipInfo3)

	var datas []*PlatformData
	datas = append(datas, &vifData)
	datas = append(datas, &vifData1)
	policy.UpdateInterfaceData(datas)
	result, _ := policy.LookupAllByKey(key)
	if result != nil {
		t.Log(result.SrcInfo, "\n")
		t.Log(result.DstInfo, "\n")
	}
}

func TestGetPlatformDataAboutArp(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)

	srcIp := NewIPFromString("192.168.2.12")
	dstIp := NewIPFromString("192.168.0.11")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	ip := NewIPFromString("192.168.0.11")
	ipInfo := IpNet{
		Ip:       ip.Int(),
		SubnetId: 121,
		Netmask:  24,
	}

	ip1 := NewIPFromString("192.168.0.12")
	ipInfo1 := IpNet{
		Ip:       ip1.Int(),
		SubnetId: 122,
		Netmask:  25,
	}

	mac := NewMACAddrFromString("08:00:27:a4:2b:fc")
	launchServer := NewIPFromString("10.10.10.10")
	vifData := PlatformData{
		EpcId:      11,
		DeviceType: 2,
		DeviceId:   3,
		IfType:     3,
		IfIndex:    5,
		Mac:        mac.Int(),
		HostIp:     launchServer.Int(),
	}

	vifData.Ips = append(vifData.Ips, &ipInfo)
	vifData.Ips = append(vifData.Ips, &ipInfo1)
	datas := make([]*PlatformData, 0, 2)
	datas = append(datas, &vifData)
	policy.UpdateInterfaceData(datas)
	now := time.Now()
	result, _ := policy.LookupAllByKey(key)
	t.Log(time.Now().Sub(now))
	if result != nil {
		t.Log(result.SrcInfo, "\n")
		t.Log(result.DstInfo, "\n")
	}
	now = time.Now()
	result, _ = policy.LookupAllByKey(key)
	t.Log(time.Now().Sub(now))
}

func TestGetGroupData(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	ip := NewIPFromString("192.168.0.11")
	ipInfo := IpNet{
		Ip:       ip.Int(),
		SubnetId: 121,
		Netmask:  32,
	}

	mac := NewMACAddrFromString("08:00:27:a4:2b:fc")
	launchServer := NewIPFromString("10.10.10.10")
	vifData := PlatformData{
		EpcId:      11,
		DeviceType: 1,
		DeviceId:   3,
		IfType:     4,
		IfIndex:    5,
		Mac:        mac.Int(),
		HostIp:     launchServer.Int(),
	}

	vifData.Ips = append(vifData.Ips, &ipInfo)
	var datas []*PlatformData
	datas = append(datas, &vifData)
	policy.UpdateInterfaceData(datas)
	ipGroup1 := &IpGroupData{
		Id:    2,
		EpcId: 11,
		Ips:   []string{"192.168.0.11/24"},
	}
	ipGroup2 := &IpGroupData{
		Id:    3,
		EpcId: 11,
		Ips:   []string{"192.168.0.11/24"},
	}
	ipGroup3 := &IpGroupData{
		Id:    4,
		EpcId: 12,
		Ips:   []string{"192.168.0.11/24"},
	}
	ipGroups := make([]*IpGroupData, 0, 2)
	ipGroups = append(ipGroups, ipGroup1)
	ipGroups = append(ipGroups, ipGroup2)
	ipGroups = append(ipGroups, ipGroup3)
	policy.UpdateIpGroupData(ipGroups)

	now := time.Now()
	result, _ := policy.LookupAllByKey(key)
	t.Log(time.Now().Sub(now))
	if result != nil {
		t.Log(result.SrcInfo, "\n")
		t.Log(result.DstInfo, "\n")
	}
	now = time.Now()
	result, _ = policy.LookupAllByKey(key)
	t.Log(time.Now().Sub(now))
}

func generatePlatformData(policy *PolicyTable) {
	ip := NewIPFromString("192.168.0.11")
	ipInfo := IpNet{
		Ip:       ip.Int(),
		SubnetId: 121,
		Netmask:  32,
	}

	mac := NewMACAddrFromString("08:00:27:a4:2b:fc")
	launchServer := NewIPFromString("10.10.10.10")
	vifData := PlatformData{
		EpcId:      11,
		DeviceType: 1,
		DeviceId:   3,
		IfType:     4,
		IfIndex:    5,
		Mac:        mac.Int(),
		HostIp:     launchServer.Int(),
	}

	vifData.Ips = append(vifData.Ips, &ipInfo)
	var datas []*PlatformData
	datas = append(datas, &vifData)
	policy.UpdateInterfaceData(datas)
}

func generateIpgroupData(policy *PolicyTable) {
	ipGroup1 := &IpGroupData{
		Id:    2,
		EpcId: 11,
		Ips:   []string{"192.168.0.11/24"},
	}
	ipGroup2 := &IpGroupData{
		Id:    3,
		EpcId: 11,
		Ips:   []string{"192.168.0.11/24"},
	}
	ipGroup3 := &IpGroupData{
		Id:    4,
		EpcId: 12,
		Ips:   []string{"192.168.0.11/24"},
	}
	ipGroups := make([]*IpGroupData, 0, 2)
	ipGroups = append(ipGroups, ipGroup1)
	ipGroups = append(ipGroups, ipGroup2)
	ipGroups = append(ipGroups, ipGroup3)
	policy.UpdateIpGroupData(ipGroups)
}

//测试全局Pass策略匹配direction==3
func TestAllPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	acl1 := &Acl{
		Id:     10,
		Type:   TAP_TOR,
		TapId:  11,
		Proto:  0,
		Vlan:   0,
		Action: []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试资源组forward策略匹配 direction==1
func TestGroupForwardPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	srcGroups := []uint32{3}
	acl1 := &Acl{
		Id:        10,
		Type:      TAP_TOR,
		TapId:     11,
		SrcGroups: srcGroups,
		Proto:     0,
		Vlan:      0,
		Action:    []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试资源组backward策略匹配 direction==2
func TestGroupBackwardPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstGroups := []uint32{3}
	acl1 := &Acl{
		Id:        10,
		Type:      TAP_TOR,
		TapId:     11,
		DstGroups: dstGroups,
		Proto:     0,
		Vlan:      0,
		Action:    []AclAction{backward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Port策略匹配 acl配置port=0，查询SrcPort=30，DstPort=30，查询到ACl
func TestAllPortPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := make([]uint16, 1)
	dstPorts = append(dstPorts, 30)
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    0,
		Vlan:     0,
		Action:   []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		SrcPort: 30,
		DstPort: 30,
		EthType: EthernetTypeARP,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Port策略匹配 acl配置port=30，查询Srcport=30，查到acl的direction=2
func TestSrcPortPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{30}
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    0,
		Vlan:     0,
		Action:   []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		SrcPort: 30,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Port策略匹配 acl配置port=30，查询Dstport=30，查到acl的direction=1
func TestDstPortPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{30}
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    0,
		Vlan:     0,
		Action:   []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 30,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Port策略匹配 acl配置port=30，查询SrcPort=30, Dstport=30，查到acl的direction=3
func TestSrcDstPortPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{30}
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    0,
		Vlan:     0,
		Action:   []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 30,
		SrcPort: 30,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Vlan策略匹配 acl配置Vlan=30，查询Vlan=30, 查询到Acl
func TestVlanPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	acl1 := &Acl{
		Id:     10,
		Type:   TAP_TOR,
		TapId:  11,
		Proto:  0,
		Vlan:   30,
		Action: []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})

	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 30,
		SrcPort: 30,
		Vlan:    30,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Vlan策略匹配 acl配置Vlan=0，Port=8000,查询Vlan=30,Port=8000 查询到Acl
func TestVlanPortPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    0,
		Vlan:     0,
		Action:   []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 30,
		SrcPort: 8000,
		Vlan:    30,
		Ttl:     64,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试Vlan策略匹配 acl配置Proto=6，Port=8000,查询Proto=6,Port=8000 查询到Acl
func TestPortProtoPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     0,
		Action:   []AclAction{forward},
	}
	policy.UpdateAclData([]*Acl{acl1})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 8000,
		SrcPort: 8000,
		Vlan:    30,
		Ttl:     64,
		Proto:   6,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试两条acl proto为6和17 查询proto=6的acl,proto为6的匹配成功
func TestAclsPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	aclAction1 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     0,
		Action:   []AclAction{aclAction1},
	}
	aclAction2 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl2 := &Acl{
		Id:       20,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    17,
		Vlan:     0,
		Action:   []AclAction{aclAction2},
	}
	policy.UpdateAclData([]*Acl{acl1, acl2})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 8000,
		SrcPort: 8000,
		Vlan:    30,
		Ttl:     64,
		Proto:   6,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试两条acl vlan为10和0  查询vlan=10的策略，结果两条都能匹配
func TestVlanAclsPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	aclAction1 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     0,
		Action:   []AclAction{aclAction1},
	}
	aclAction2 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl2 := &Acl{
		Id:       20,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     10,
		Action:   []AclAction{aclAction2},
	}
	policy.UpdateAclData([]*Acl{acl1, acl2})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 8000,
		SrcPort: 8000,
		Vlan:    10,
		Ttl:     64,
		Proto:   6,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)

	aclAction2 = aclAction2.SetDirections(FORWARD)
	aclAction2Backward := aclAction2
	aclAction2Backward.SetDirections(BACKWARD)

	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{aclAction2, aclAction2Backward, forward, backward}, 20)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试两条acl vlan=10和port=8000  查询vlan=10,port=1000，匹配到vlan=10的策略
func TestVlanPortAclsPassPolicy(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	aclAction1 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     0,
		Action:   []AclAction{aclAction1},
	}
	aclAction2 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddDirections(FORWARD).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl2 := &Acl{
		Id:     20,
		Type:   TAP_TOR,
		TapId:  11,
		Proto:  6,
		Vlan:   10,
		Action: []AclAction{aclAction2},
	}
	policy.UpdateAclData([]*Acl{acl1, acl2})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 1000,
		Vlan:    10,
		Ttl:     64,
		Proto:   6,
		Tap:     TAP_TOR,
	}
	backward := getBackwardAcl(aclAction2)
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{aclAction2, backward}, 20)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试两条acl vlan=10和port=8000  查询vlan=10,port=8000，两条策略都匹配到
func TestVlanPortAclsPassPolicy1(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	aclAction1 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     0,
		Action:   []AclAction{aclAction1},
	}
	aclAction2 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddDirections(FORWARD).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl2 := &Acl{
		Id:     20,
		Type:   TAP_TOR,
		TapId:  11,
		Proto:  6,
		Vlan:   10,
		Action: []AclAction{aclAction2},
	}
	policy.UpdateAclData([]*Acl{acl1, acl2})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 8000,
		Vlan:    10,
		Ttl:     64,
		Proto:   6,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	acl2Backward := getBackwardAcl(aclAction2)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{aclAction2, acl2Backward, forward}, 20)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

//测试两条acl vlan=10和port=8000  查询port=8000，匹配到port=8000的策略
func TestVlanPortAclsPassPolicy2(t *testing.T) {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)
	generatePlatformData(policy)
	generateIpgroupData(policy)
	dstPorts := []uint16{8000}
	aclAction1 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl1 := &Acl{
		Id:       10,
		Type:     TAP_TOR,
		TapId:    11,
		DstPorts: dstPorts,
		Proto:    6,
		Vlan:     0,
		Action:   []AclAction{aclAction1},
	}
	aclAction2 := AclAction(0).AddActionFlags(ACTION_PACKET_COUNTING).AddTagTemplates(TEMPLATE_EDGE_PORT)
	acl2 := &Acl{
		Id:     20,
		Type:   TAP_TOR,
		TapId:  11,
		Proto:  6,
		Vlan:   10,
		Action: []AclAction{aclAction2},
	}
	policy.UpdateAclData([]*Acl{acl1, acl2})
	srcIp := NewIPFromString("192.168.0.11")
	dstIp := NewIPFromString("192.168.0.12")
	key := &LookupKey{
		SrcIp:   srcIp.Int(),
		SrcMac:  0x80027a42bfc,
		DstMac:  0x80027a42bfa,
		DstIp:   dstIp.Int(),
		EthType: EthernetTypeARP,
		DstPort: 8000,
		Ttl:     64,
		Proto:   6,
		Tap:     TAP_TOR,
	}
	_, policyData := policy.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{forward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	_, policyData = policy.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto FastPath Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

// 以下是云平台信息和policy结合起来的测试
var (
	server = NewIPFromString("172.20.1.1").Int()

	group1ip1  = NewIPFromString("192.168.1.10").Int()
	group1mac  = NewMACAddrFromString("11:11:11:11:11:11").Int()
	group1ip2  = NewIPFromString("192.168.1.20").Int()
	group1mac2 = NewMACAddrFromString("11:11:11:11:11:12").Int()
	group1ip3  = NewIPFromString("102.168.33.22").Int()
	group1Id   = uint32(10)

	group2ip1 = NewIPFromString("10.30.1.10").Int()
	group2mac = NewMACAddrFromString("22:22:22:22:22:22").Int()
	group2ip2 = NewIPFromString("10.30.1.20").Int()
	group2Id  = uint32(20)
)

func generateIpNet(ip uint32, mask uint32) *IpNet {
	ipInfo := IpNet{
		Ip:       ip,
		SubnetId: 121,
		Netmask:  mask,
	}
	return &ipInfo
}

func generatePlatformDataWithGroupId(epcId int32, groupId uint32, mac uint64) *PlatformData {
	data := PlatformData{
		EpcId:      epcId,
		DeviceType: 2,
		DeviceId:   3,
		IfType:     3,
		IfIndex:    5,
		Mac:        mac,
		HostIp:     server,
	}
	data.GroupIds = append(data.GroupIds, groupId)
	return &data
}

func generatePolicyTable() *PolicyTable {
	policy := NewPolicyTable(ACTION_PACKET_COUNTING, 1, 1024, false)

	datas := make([]*PlatformData, 0, 2)

	data := generatePlatformDataWithGroupId(int32(group1Id), group1Id, group1mac)
	ip := generateIpNet(group1ip1, 24)
	data.Ips = append(data.Ips, ip)
	ip = generateIpNet(group1ip2, 25)
	data.Ips = append(data.Ips, ip)
	datas = append(datas, data)

	data = generatePlatformDataWithGroupId(int32(group1Id), 0, group1mac2)
	ip = generateIpNet(group1ip3, 18)
	data.Ips = append(data.Ips, ip)
	datas = append(datas, data)

	data = generatePlatformDataWithGroupId(int32(group2Id), group2Id, group2mac)
	ip = generateIpNet(group2ip1, 24)
	data.Ips = append(data.Ips, ip)
	ip = generateIpNet(group2ip2, 25)
	data.Ips = append(data.Ips, ip)
	datas = append(datas, data)

	policy.UpdateInterfaceData(datas)
	return policy
}

func generateAclAction(id ACLID, actionFlags ActionFlag) AclAction {
	return AclAction(0).AddActionFlags(actionFlags).AddDirections(FORWARD).AddTagTemplates(TEMPLATE_EDGE_PORT)
}

func generatePolicyAcl(table *PolicyTable, action AclAction, aclID ACLID, srcGroupId, dstGroupId uint32,
	proto uint8, port uint16, vlan uint32) *Acl {
	srcGroups := make([]uint32, 0, 1)
	dstGroups := make([]uint32, 0, 1)
	dstPorts := make([]uint16, 0, 1)

	srcGroups = append(srcGroups, srcGroupId)
	dstGroups = append(dstGroups, dstGroupId)
	if port != 0 {
		dstPorts = append(dstPorts, port)
	}
	acl := &Acl{
		Id:        aclID,
		Type:      TAP_TOR,
		TapId:     uint32(aclID + 1),
		SrcGroups: srcGroups,
		DstGroups: dstGroups,
		DstPorts:  dstPorts,
		Proto:     proto,
		Vlan:      vlan,
		Action:    []AclAction{action},
	}
	return acl
}

func generateLookupKey(srcMac, dstMac uint64, vlan uint16, srcIp, dstIp uint32,
	proto uint8, srcPort, dstPort uint16) *LookupKey {
	key := &LookupKey{
		SrcMac:  srcMac,
		DstMac:  dstMac,
		SrcIp:   srcIp,
		DstIp:   dstIp,
		Proto:   proto,
		SrcPort: srcPort,
		DstPort: dstPort,
		Vlan:    vlan,
		Tap:     TAP_TOR,
	}
	return key
}

func TestPolicySimple(t *testing.T) {
	acls := []*Acl{}
	// 创建 policyTable
	table := generatePolicyTable()
	// 构建acl action  1->2 tcp 8000
	action := generateAclAction(10, ACTION_PACKET_COUNTING)
	acl := generatePolicyAcl(table, action, 10, group1Id, group2Id, 6, 8000, 0)
	acls = append(acls, acl)
	table.UpdateAclData(acls)
	// 构建查询key  1:0->2:8000 tcp
	key := generateLookupKey(group1mac, group2mac, 0, group1ip1, group2ip1, 6, 0, 8000)

	// 获取查询first结果
	_, policyData := table.LookupAllByKey(key)
	// 构建预期结果
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{action}, 10)
	// 查询结果和预期结果比较
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	// 构建查询key  2:8000->1:0 tcp
	key = generateLookupKey(group2mac, group1mac, 0, group2ip1, group1ip1, 6, 8000, 0)
	// key和acl方向相反，构建反向的action
	backward := getBackwardAcl(action)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{backward}, 10)
	// 查询结果和预期结果比较
	_, policyData = table.LookupAllByKey(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	// 构建无效查询key  2:0->1:8000 tcp
	key = generateLookupKey(group2mac, group1mac, 0, group2ip1, group1ip1, 6, 0, 8000)
	_, policyData = table.LookupAllByKey(key)
	basicPolicyData = INVALID_POLICY_DATA
	// key不匹配，返回无效policy
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	// 测试同样的key, 匹配两条action
	action2 := generateAclAction(12, ACTION_PACKET_COUNTING)
	acl2 := generatePolicyAcl(table, action2, 12, group1Id, group2Id, 6, 8000, 0)
	acls = append(acls, acl2)
	table.UpdateAclData(acls)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{action, action2}, 10)
	key = generateLookupKey(group1mac, group2mac, 0, group1ip1, group2ip1, 6, 0, 8000)

	_, policyData = table.LookupAllByKey(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

func TestPolicyEpcPolicy(t *testing.T) {
	acls := []*Acl{}
	// 创建 policyTable
	table := generatePolicyTable()
	// 构建acl action  1->2 tcp 8000
	action := generateAclAction(10, ACTION_PACKET_COUNTING)
	acl := generatePolicyAcl(table, action, 10, group1Id, 0, 6, 8000, 0)
	acls = append(acls, acl)
	table.UpdateAclData(acls)
	// 构建查询key  1:0->2:8000 tcp
	key := generateLookupKey(group1mac, group1mac2, 0, group1ip1, group1ip3, 6, 0, 8000)

	// 获取查询first结果
	_, policyData := table.LookupAllByKey(key)
	// 构建预期结果
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{action}, 10)
	// 查询结果和预期结果比较
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("TestPolicyEpcPolicy Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	// 查询结果和预期结果比较
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("TestPolicyEpcPolicy Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	backward := getBackwardAcl(action)
	key = generateLookupKey(group1mac2, group1mac, 0, group1ip3, group1ip1, 6, 8000, 0)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{backward}, 10)
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	// 查询结果和预期结果比较
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("TestPolicyEpcPolicy Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	key = generateLookupKey(group1mac2, group1mac, 0, group1ip3, group1ip1, 6, 0, 8000)
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	basicPolicyData = nil
	// 查询结果和预期结果比较
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("TestPolicyEpcPolicy Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	_, policyData = table.LookupAllByKey(key)
	basicPolicyData = INVALID_POLICY_DATA
	// 查询结果和预期结果比较
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("TestPolicyEpcPolicy Check Failed")
		t.Log("Result:", policyData, "\n")
	}
}

func TestFlowVlanAcls(t *testing.T) {
	acls := []*Acl{}
	table := generatePolicyTable()
	action := generateAclAction(10, ACTION_FLOW_COUNTING)
	acl := generatePolicyAcl(table, action, 10, group1Id, group2Id, 6, 0, 10)
	acls = append(acls, acl)
	table.UpdateAclData(acls)
	// 构建查询key  1->2 tcp vlan:10
	key := generateLookupKey(group1mac, group2mac, 10, group1ip1, group2ip1, 6, 11, 10)
	_, policyData := table.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{action}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	// key和acl方向相反，构建反向的action
	backward := getBackwardAcl(action)
	basicPolicyData2 := NewPolicyData()
	basicPolicyData2.Merge([]AclAction{backward}, 10)
	key = generateLookupKey(group2mac, group1mac, 10, group2ip1, group1ip1, 6, 11, 10)
	_, policyData2 := table.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData2, policyData2) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData2, "\n")
		t.Log("Expect:", basicPolicyData2, "\n")
	}

	// key不匹配，返回无效policy
	key = generateLookupKey(group2mac, group1mac, 11, group2ip1, group1ip1, 6, 11, 10)
	_, policyData3 := table.LookupAllByKey(key)
	basicPolicyData3 := INVALID_POLICY_DATA
	if !CheckPolicyResult(basicPolicyData3, policyData3) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData3, "\n")
		t.Log("Expect:", basicPolicyData3, "\n")
	}
}

func TestVlanPortAcl(t *testing.T) {
	acls := []*Acl{}
	table := generatePolicyTable()
	// group1->group2,tcp,vlan:10,dstport:20
	action := generateAclAction(10, ACTION_FLOW_COUNTING)
	acl := generatePolicyAcl(table, action, 10, group1Id, group2Id, 6, 20, 10)
	// group2->group1,tcp,vlan:10,dstport:21
	action2 := generateAclAction(12, ACTION_FLOW_COUNTING)
	acl2 := generatePolicyAcl(table, action2, 12, group2Id, group1Id, 6, 21, 10)
	acls = append(acls, acl)
	acls = append(acls, acl2)
	table.UpdateAclData(acls)
	// 构建查询key  1:21->2:20 tcp vlan:10 ,匹配两条acl
	key := generateLookupKey(group1mac, group2mac, 10, group1ip1, group2ip1, 6, 21, 20)
	_, policyData := table.LookupAllByKey(key)
	backward := getBackwardAcl(action2)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{action, backward}, 10)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}

func TestVlanPortAcl2(t *testing.T) {
	acls := []*Acl{}
	table := generatePolicyTable()
	// group1->group2, vlan:10
	action1 := generateAclAction(11, ACTION_FLOW_COUNTING)
	acl1 := generatePolicyAcl(table, action1, 11, group1Id, group2Id, 0, 0, 10)
	// group1->group2, proto:6
	action2 := generateAclAction(12, ACTION_FLOW_COUNTING)
	acl2 := generatePolicyAcl(table, action2, 12, group1Id, group2Id, 6, 0, 0)
	// group1->group2, port:80
	action3 := generateAclAction(13, ACTION_FLOW_COUNTING)
	acl3 := generatePolicyAcl(table, action3, 13, group1Id, group2Id, 0, 80, 0)
	acls = append(acls, acl1)
	acls = append(acls, acl2)
	acls = append(acls, acl3)
	table.UpdateAclData(acls)
	// 构建查询1-key  1:10->2:10 proto:6 vlan:10
	key := generateLookupKey(group1mac, group2mac, 10, group1ip1, group2ip1, 6, 10, 10)
	// 获取first查询结果
	_, policyData := table.LookupAllByKey(key)
	basicPolicyData := NewPolicyData()
	basicPolicyData.Merge([]AclAction{action1, action2}, 11)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 获取fastpath查询结果
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 2-key: 1:10 -> 2:80 proto:1 vlan:10
	key = generateLookupKey(group1mac, group2mac, 10, group1ip1, group2ip1, 1, 10, 80)
	// 获取first查询结果
	_, policyData = table.LookupAllByKey(key)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{action1, action3}, 11)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 获取fastpath查询结果
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 3-key: 1:10 -> 2:80 proto:6 vlan:0
	key = generateLookupKey(group1mac, group2mac, 0, group1ip1, group2ip1, 6, 10, 80)
	// 获取first查询结果
	_, policyData = table.LookupAllByKey(key)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{action2, action3}, 12)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 获取fastpath查询结果
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}

	acls = []*Acl{}
	table = generatePolicyTable()
	// port:80
	action4 := generateAclAction(14, ACTION_FLOW_COUNTING)
	acl4 := generatePolicyAcl(table, action4, 14, 0, 0, 0, 80, 0)
	// group1->group2, proto:6
	action5 := generateAclAction(15, ACTION_FLOW_COUNTING)
	acl5 := generatePolicyAcl(table, action5, 15, group1Id, group2Id, 6, 0, 0)
	acls = append(acls, acl4)
	acls = append(acls, acl5)
	table.UpdateAclData(acls)
	// 4-key  1:10->2:80 proto:6
	key = generateLookupKey(group1mac, group2mac, 0, group1ip1, group2ip1, 6, 10, 80)
	// 获取first查询结果
	_, policyData = table.LookupAllByKey(key)
	backward1 := getBackwardAcl(action4)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{action5, action4}, 15)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 获取fastpath查询结果
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 5-key 2:80->1:10 proto:6
	key = generateLookupKey(group2mac, group1mac, 0, group2ip1, group1ip1, 6, 80, 10)
	// 获取first查询结果
	_, policyData = table.LookupAllByKey(key)
	backward2 := getBackwardAcl(action5)
	basicPolicyData = NewPolicyData()
	basicPolicyData.Merge([]AclAction{backward2, backward1}, 15)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
	// 获取fastpath查询结果
	_, policyData = table.policyLabeler.GetPolicyByFastPath(key)
	if !CheckPolicyResult(basicPolicyData, policyData) {
		t.Error("PortProto Check Failed")
		t.Log("Result:", policyData, "\n")
		t.Log("Expect:", basicPolicyData, "\n")
	}
}
