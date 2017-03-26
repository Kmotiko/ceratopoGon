package message

import (
	"encoding/binary"
)

func UnMarshall(packet []byte) (msg MqttSnMessage) {
	h := MqttSnHeader{}
	h.UnMarshall(packet)
	switch h.MsgType {
	case MQTTSNT_ADVERTISE:
		msg = NewAdvertise()
		msg.UnMarshall(packet)
	case MQTTSNT_SEARCHGW:
		msg = NewSearchGw()
		msg.UnMarshall(packet)
	case MQTTSNT_GWINFO:
		msg = NewGwInfo()
		msg.UnMarshall(packet)
	case MQTTSNT_CONNECT:
		msg = NewConnect()
		msg.UnMarshall(packet)
	case MQTTSNT_CONNACK:
		msg = NewConnAck()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPICREQ:
		msg = NewWillTopicReq()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPIC:
		msg = NewWillTopic()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSGREQ:
		msg = NewWillMsgReq()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSG:
		msg = NewWillMsg()
		msg.UnMarshall(packet)
	case MQTTSNT_REGISTER:
		msg = NewRegister()
		msg.UnMarshall(packet)
	case MQTTSNT_REGACK:
		msg = NewRegAck()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBLISH:
		msg = NewPublish()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBACK:
		msg = NewPubAck()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBCOMP:
		msg = NewPubComp()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBREC:
		msg = NewPubRec()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBREL:
		msg = NewPubRel()
		msg.UnMarshall(packet)
	case MQTTSNT_SUBSCRIBE:
		msg = NewSubscribe()
		msg.UnMarshall(packet)
	case MQTTSNT_SUBACK:
		msg = NewSubAck()
		msg.UnMarshall(packet)
	case MQTTSNT_UNSUBSCRIBE:
		msg = NewUnSubscribe()
		msg.UnMarshall(packet)
	case MQTTSNT_UNSUBACK:
		msg = NewUnSubAck()
		msg.UnMarshall(packet)
	case MQTTSNT_PINGREQ:
		msg = NewPingReq()
		msg.UnMarshall(packet)
	case MQTTSNT_PINGRESP:
		msg = NewPingResp()
		msg.UnMarshall(packet)
	case MQTTSNT_DISCONNECT:
		msg = NewDisConnect()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPICUPD:
		msg = NewWillTopicUpd()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPICRESP:
		msg = NewWillTopicResp()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSGUPD:
		msg = NewWillMsgUpd()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSGRESP:
		msg = NewWillMsgResp()
		msg.UnMarshall(packet)
	}

	return msg
}

func mqttsnFlags(dup uint8, qos uint8, retain uint8, will uint8, csession uint8, tidType uint8) uint8 {
	return dup << 7 & qos << 5 & retain << 4 & will << 3 & csession << 2 & tidType
}

func dupFlag(flags uint8) bool {
	dup := false
	if (flags >> 7) > 0 {
		dup = true
	}
	return dup
}

func qosFlag(flags uint8) uint8 {
	return flags >> 5 & 0x03
}

func hasRetainFlag(flags uint8) bool {
	retain := false
	if (flags >> 4 & 0x01) > 0 {
		retain = true
	}
	return retain
}

func hasWillFlag(flags uint8) bool {
	will := false
	if (flags >> 3 & 0x01) > 0 {
		will = true
	}
	return will
}

func hasCleanSessionFlag(flags uint8) bool {
	csession := false
	if (flags >> 2 & 0x01) > 0 {
		csession = true
	}
	return csession
}

func topicIdType(flags uint8) uint8 {
	return flags & 0x03
}

/*********************************************/
/* MqttSnHeader                              */
/*********************************************/
// func NewMqttSnHeader(l uint16, t uint8) *MqttSnHeader {
// 	h := &MqttSnHeader{l, t}
// 	return h
// }

func (m *MqttSnHeader) Marshall() []byte {
	packet := make([]byte, m.Size())
	index := 0
	if m.Size() == 2 {
		packet[index] = (uint8)(m.Length)
		index += 1
	} else if m.Size() == 3 {
		binary.BigEndian.PutUint16(packet[index:], m.Length)
		index += 2
	}
	packet[index] = m.MsgType

	return packet
}

func (m *MqttSnHeader) UnMarshall(packet []byte) {
	index := 0
	length := packet[index]
	if length == 0x01 {
		m.Length = binary.BigEndian.Uint16(packet[index:])
		index += 2
	} else {
		m.Length = (uint16)(length)
		index++
	}
	m.MsgType = packet[index]
}

func (m *MqttSnHeader) Size() int {
	size := 3
	// if length is smaller than 256
	if m.Length < 256 {
		size = 2
	}

	return (int)(size)
}

/*********************************************/
/* Advertise                                 */
/*********************************************/
func NewAdvertise() *Advertise {
	m := &Advertise{
		Header: MqttSnHeader{5, MQTTSNT_ADVERTISE}}
	return m
}

func (m *Advertise) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.GwId
	binary.BigEndian.PutUint16(packet[index:], m.Duration)

	return packet
}

func (m *Advertise) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Size()

	m.GwId = packet[index]
	index++

	m.Duration = binary.BigEndian.Uint16(packet[index:])
}

func (m *Advertise) Size() int {
	return 5
}

/*********************************************/
/* SearchGw                                  */
/*********************************************/
func NewSearchGw() *SearchGw {
	m := &SearchGw{
		Header: MqttSnHeader{3, MQTTSNT_SEARCHGW}}
	return m
}

func (m *SearchGw) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.Radius

	return packet
}

func (m *SearchGw) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.Radius = packet[index]
}

func (m *SearchGw) Size() int {
	return 3
}

/*********************************************/
/* GwInfo                                    */
/*********************************************/
func NewGwInfo() *GwInfo {
	m := &GwInfo{
		Header: MqttSnHeader{3, MQTTSNT_GWINFO}}
	return m
}

func (m *GwInfo) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Size()

	packet[index] = m.GwId
	index++

	copy(packet[index:], m.GwAdd)

	return packet
}

func (m *GwInfo) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.GwId = packet[index]
	index++

	addrLen := m.Header.Length - 3
	m.GwAdd = make([]uint8, addrLen)
	copy(m.GwAdd, packet[index:])
}

func (m *GwInfo) Size() int {
	size := 3
	size += len(m.GwAdd)
	return size
}

/*********************************************/
/* Connect                                   */
/*********************************************/
func NewConnect() *Connect {
	m := &Connect{
		Header: MqttSnHeader{6, MQTTSNT_CONNECT}}
	return m
}

func (m *Connect) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.Flags
	index++

	packet[index] = m.ProtocolId
	index++

	binary.BigEndian.PutUint16(packet[index:], m.Duration)
	index += 2

	copy(packet[index:], []byte(m.ClientId))

	return packet
}

func (m *Connect) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.Flags = packet[index]
	index++

	m.ProtocolId = packet[index]
	index++

	m.Duration = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.ClientId = string(packet[index:m.Header.Length])
}

func (m *Connect) Size() int {
	size := 6
	size += len(m.ClientId)
	return size
}

/*********************************************/
/* ConnAck                                   */
/*********************************************/
func NewConnAck() *ConnAck {
	m := &ConnAck{
		Header: MqttSnHeader{3, MQTTSNT_CONNACK}}
	return m
}

func (m *ConnAck) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.ReturnCode

	return packet
}

func (m *ConnAck) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Size()

	m.ReturnCode = packet[index]
}

func (m *ConnAck) Size() int {
	return 3
}

/*********************************************/
/* WillTopicReq                              */
/*********************************************/
func NewWillTopicReq() *MqttSnHeader {
	m := &MqttSnHeader{2, MQTTSNT_WILLTOPICREQ}
	return m
}

/*********************************************/
/* WillTopic                                 */
/*********************************************/
func NewWillTopic() *WillTopic {
	return nil
}

func (m *WillTopic) Marshall() []byte {
	return nil
}

func (m *WillTopic) UnMarshall(packet []byte) {
}

func (m *WillTopic) Size() int {
	return 0
}

/*********************************************/
/* WillMsgReq                                */
/*********************************************/
func NewWillMsgReq() *MqttSnHeader {
	return nil
}

/*********************************************/
/* WillMsg                                   */
/*********************************************/
func NewWillMsg() *WillMsg {
	return nil
}

func (m *WillMsg) Marshall() []byte {
	return nil
}

func (m *WillMsg) UnMarshall(packet []byte) {
}

func (m *WillMsg) Size() int {
	return 0
}

/*********************************************/
/* Register                                  */
/*********************************************/
func NewRegister() *Register {
	return nil
}

func (m *Register) Marshall() []byte {
	return nil
}

func (m *Register) UnMarshall(packet []byte) {
}

func (m *Register) Size() int {
	return 0
}

/*********************************************/
/* RegAck                                    */
/*********************************************/
func NewRegAck() *RegAck {
	return nil
}

func (m *RegAck) Marshall() []byte {
	return nil
}

func (m *RegAck) UnMarshall(packet []byte) {
}

func (m *RegAck) Size() int {
	return 0
}

/*********************************************/
/* Publish                                   */
/*********************************************/
func NewPublish() *Publish {
	return nil
}

func (m *Publish) Marshall() []byte {
	return nil
}

func (m *Publish) UnMarshall(packet []byte) {
}

func (m *Publish) Size() int {
	return 0
}

/*********************************************/
/* PubAck                                    */
/*********************************************/
func NewPubAck() *PubAck {
	return nil
}

func (m *PubAck) Marshall() []byte {
	return nil
}

func (m *PubAck) UnMarshall(packet []byte) {
}

func (m *PubAck) Size() int {
	return 0
}

/*********************************************/
/* PubRec                                    */
/*********************************************/
func NewPubRec() *PubRec {
	return nil
}

func (m *PubRec) Marshall() []byte {
	return nil
}

func (m *PubRec) UnMarshall(packet []byte) {
}

func (m *PubRec) Size() int {
	return 0
}

/*********************************************/
/* PubRel                                    */
/*********************************************/
func NewPubRel() *PubRec {
	return nil
}

/*********************************************/
/* PubComp                                   */
/*********************************************/
func NewPubComp() *PubRec {
	return nil
}

/*********************************************/
/* Subscribe                                 */
/*********************************************/
func NewSubscribe() *Subscribe {
	return nil
}

func (m *Subscribe) Marshall() []byte {
	return nil
}

func (m *Subscribe) UnMarshall(packet []byte) {
}

func (m *Subscribe) Size() int {
	return 0
}

/*********************************************/
/* SubAck                                 */
/*********************************************/
func NewSubAck() *SubAck {
	return nil
}

func (m *SubAck) Marshall() []byte {
	return nil
}

func (m *SubAck) UnMarshall(packet []byte) {
}

func (m *SubAck) Size() int {
	return 0
}

/*********************************************/
/* UnSubscribe                               */
/*********************************************/
func NewUnSubscribe() *UnSubscribe {
	return nil
}

func (m *UnSubscribe) Marshall() []byte {
	return nil
}

func (m *UnSubscribe) UnMarshall(packet []byte) {
}

func (m *UnSubscribe) Size() int {
	return 0
}

/*********************************************/
/* UnSubAck                                  */
/*********************************************/
func NewUnSubAck() *UnSubAck {
	return nil
}

func (m *UnSubAck) Marshall() []byte {
	return nil
}

func (m *UnSubAck) UnMarshall(packet []byte) {
}

func (m *UnSubAck) Size() int {
	return 0
}

/*********************************************/
/* PingReq                                  */
/*********************************************/
func NewPingReq() *PingReq {
	return nil
}

func (m *PingReq) Marshall() []byte {
	return nil
}

func (m *PingReq) UnMarshall(packet []byte) {
}

func (m *PingReq) Size() int {
	return 0
}

/*********************************************/
/* PingResp                                  */
/*********************************************/
func NewPingResp() *MqttSnHeader {
	return nil
}

/*********************************************/
/* DisConnect                                */
/*********************************************/
func NewDisConnect() *DisConnect {
	return nil
}

func (m *DisConnect) Marshall() []byte {
	return nil
}

func (m *DisConnect) UnMarshall(packet []byte) {
}

func (m *DisConnect) Size() int {
	return 0
}

/*********************************************/
/* WillTopicUpd                                */
/*********************************************/
func NewWillTopicUpd() *WillTopicUpd {
	return nil
}

func (m *WillTopicUpd) Marshall() []byte {
	return nil
}

func (m *WillTopicUpd) UnMarshall(packet []byte) {
}

func (m *WillTopicUpd) Size() int {
	return 0
}

/*********************************************/
/* WillMsgUpd                                */
/*********************************************/
func NewWillMsgUpd() *WillMsgUpd {
	return nil
}

func (m *WillMsgUpd) Marshall() []byte {
	return nil
}

func (m *WillMsgUpd) UnMarshall(packet []byte) {
}

func (m *WillMsgUpd) Size() int {
	return 0
}

/*********************************************/
/* WillTopicResp                             */
/*********************************************/
func NewWillTopicResp() *WillTopicResp {
	return nil
}

func (m *WillTopicResp) Marshall() []byte {
	return nil
}

func (m *WillTopicResp) UnMarshall(packet []byte) {
}

func (m *WillTopicResp) Size() int {
	return 0
}

/*********************************************/
/* WillMsgResp                             */
/*********************************************/
func NewWillMsgResp() *WillTopicResp {
	return nil
}
