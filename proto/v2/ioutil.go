package v2

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func GenerateCodeBytesArray(s string) ([6]byte, error) {
	if len(s) != 6 {
		return [6]byte{}, fmt.Errorf("GenerateCodeBytesArray error")
	}
	var b [6]byte
	copy(b[:], s)
	return b, nil
}

func ReadCode(b []byte, cursor *int) (string, error) {
	if len(b) < *cursor+6 {
		return "", errors.New("read code overflow")
	}
	defer func() {
		*cursor += 6
	}()
	return string(b[*cursor : *cursor+6]), nil
}

func ReadTDXInt(b []byte, cursor *int) (int64, error) {
	value := int64(0)

	shift := 6

	cur := b[*cursor]
	*cursor += 1

	value += int64(cur & 0b00111111)
	negative := false
	if (cur & 0b01000000) > 0 {
		negative = true
	}

	for (cur & 0b10000000) > 0 {
		cur = b[*cursor]
		*cursor += 1

		value += (int64(cur&0b01111111) << shift)

		shift += 7
	}

	if negative {
		value = -value
	}
	return value, nil
}

func ReadTDXFloat(data []byte, cursor *int) (float64, error) {
	if len(data) < *cursor+4 {
		return 0, errors.New("read float overflow")
	}
	lleax := int(data[*cursor]) //[0]
	*cursor += 1
	lheax := int(data[*cursor]) //[1]
	*cursor += 1
	hleax := int(data[*cursor]) // [2]
	*cursor += 1
	logpoint := int(data[*cursor])
	*cursor += 1

	dwEcx := logpoint*2 - 0x7f // 0b0111 1111
	dwEdx := logpoint*2 - 0x86 // 0b1000 0110
	dwEsi := logpoint*2 - 0x8e // 0b1000 1110
	dwEax := logpoint*2 - 0x96 // 0b1001 0110
	tmpEax := dwEcx
	if dwEcx < 0 {
		tmpEax = -dwEcx
	} else {
		tmpEax = dwEcx
	}

	dbl_xmm6 := 0.0
	dbl_xmm6 = math.Pow(2.0, float64(tmpEax))
	if dwEcx < 0 {
		dbl_xmm6 = 1.0 / dbl_xmm6
	}

	dbl_xmm4 := 0.0
	dbl_xmm0 := 0.0

	if hleax > 0x80 {
		tmpdbl_xmm3 := 0.0
		dwtmpeax := dwEdx + 1
		tmpdbl_xmm3 = math.Pow(2.0, float64(dwtmpeax))
		dbl_xmm0 = math.Pow(2.0, float64(dwEdx)) * 128.0
		dbl_xmm0 += float64(hleax&0x7f) * tmpdbl_xmm3
		dbl_xmm4 = dbl_xmm0
	} else {
		if dwEdx >= 0 {
			dbl_xmm0 = math.Pow(2.0, float64(dwEdx)) * float64(hleax)
		} else {
			dbl_xmm0 = (1 / math.Pow(2.0, float64(dwEdx))) * float64(hleax)
		}
		dbl_xmm4 = dbl_xmm0
	}

	dbl_xmm3 := math.Pow(2.0, float64(dwEsi)) * float64(lheax)
	dbl_xmm1 := math.Pow(2.0, float64(dwEax)) * float64(lleax)
	if (hleax & 0x80) > 0 {
		dbl_xmm3 *= 2.0
		dbl_xmm1 *= 2.0
	}
	return dbl_xmm6 + dbl_xmm4 + dbl_xmm3 + dbl_xmm1, nil
}

func ReadInt[T any](b []byte, cursor *int, _type T) (T, error) {
	var length int = 0
	switch any(_type).(type) {
	case int8, uint8:
		length = 1
	case int16, uint16:
		length = 2
	case int32, uint32, int, uint:
		length = 4
	case int64, uint64:
		length = 8
	default:
		return _type, fmt.Errorf("unsupported type:%T", _type)
	}
	err := binary.Read(bytes.NewBuffer(b[*cursor:*cursor+length]), binary.LittleEndian, &_type)
	if err != nil {
		return _type, fmt.Errorf("ReadAsInt error:%s", err)
	}
	*cursor += length
	return _type, nil
}

func ReadByteArray(b []byte, cursor *int, length int) ([]byte, error) {
	if *cursor+length > len(b) {
		return nil, errors.New("ReadAsByteArray overflow")
	}
	defer func() {
		*cursor += length
	}()
	return b[*cursor : *cursor+length], nil
}

func ReadTDXString(b []byte, cursor *int, fixedLength int) (string, error) {
	nameGBKBuf, err := ReadByteArray(b, cursor, fixedLength)
	if err != nil {
		return "", err
	}
	nameUtf8Buf, err := simplifiedchinese.GBK.NewDecoder().Bytes(nameGBKBuf)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(nameUtf8Buf), "\x00"), nil
}

func ReadTDXTime(b []byte, pos *int, t CandleStickPeriodType) (year int, month int, day int, hour int, minute int, err error) {
	switch t {
	case CandleStickPeriodType_5Min,
		CandleStickPeriodType_15Min,
		CandleStickPeriodType_30Min,
		CandleStickPeriodType_1Hour,
		CandleStickPeriodType_1Min:
		var zipday, tminutes uint16
		zipday, err = ReadInt(b, pos, zipday)
		if err != nil {
			return
		}
		tminutes, err = ReadInt(b, pos, tminutes)
		if err != nil {
			return
		}
		year = int((zipday >> 11) + 2004)
		month = int((zipday % 2048) / 100)
		day = int((zipday % 2048) % 100)
		hour = int(tminutes / 60)
		minute = int(tminutes % 60)
		return
	case CandleStickPeriodType_Day,
		CandleStickPeriodType_Week,
		CandleStickPeriodType_Month:
		var zipday uint32
		zipday, err = ReadInt(b, pos, zipday)
		year = int(zipday / 10000)
		month = int((zipday % 10000) / 100)
		day = int(zipday % 100)
		hour = 15
	default:
		err = errors.New("unsupported CandleStickPeriodType")
		return
	}
	return
}

///////////////////////////////////////////
// codec
///////////////////////////////////////////

type tdxCodec struct {
	book []uint32
}

func (c *tdxCodec) shift(b0, b1 uint32) (uint32, uint32, error) {
	v := make([]uint32, 19)
	var (
		v3 uint32 = 0
		v4 uint32 = 0
		v5 uint32 = 0
	)
	v3 += c.book[tdxCodecBookOffset1+((b1^c.book[1]^c.shiftUint32(c.book[0]^b0))>>24)&0xFF] +
		c.book[tdxCodecBookOffset2+((b1^c.book[1]^c.shiftUint32(c.book[0]^b0))>>16)&0xFF]

	v4 = c.book[2] ^ c.book[0] ^ b0 ^
		(c.book[0xFF&(b1^c.book[1]^(0xFF&c.shiftUint32(c.book[0]^b0)))+tdxCodecBookOffset4] +
			(c.book[0xFF&((b1^c.book[1]^c.shiftUint32(c.book[0]^b0))>>8)+tdxCodecBookOffset3] ^ v3))
	v5 = c.book[3] ^ b1 ^ c.book[1] ^ c.shiftUint32(v4) ^ c.shiftUint32(c.book[0]^b0)

	v[4] = v4
	v[5] = v5
	for i := 6; i <= 18; i += 1 {
		v[i] = c.book[i-2] ^ c.shiftUint32(v[i-1]) ^ v[i-2]
	}

	return c.book[17] ^ v[17], v[18], nil
}

const (
	tdxCodecBookOffset0 = 0
	tdxCodecBookOffset1 = tdxCodecBookOffset0 + 0x12
	tdxCodecBookOffset2 = tdxCodecBookOffset1 + 0x100
	tdxCodecBookOffset3 = tdxCodecBookOffset2 + 0x100
	tdxCodecBookOffset4 = tdxCodecBookOffset3 + 0x100
)

func (c *tdxCodec) shiftUint32(v uint32) uint32 {
	v0 := c.book[tdxCodecBookOffset1+(v>>24)&0xFF]
	v1 := c.book[tdxCodecBookOffset2+(v>>16)&0xFF]
	v2 := c.book[tdxCodecBookOffset3+(v>>8)&0xFF]
	v3 := c.book[tdxCodecBookOffset4+(v>>0)&0xFF]
	return v3 + (v2 ^ (v1 + v0))
}

func (c *tdxCodec) dumpBook() string {
	b := bytes.NewBuffer(nil)
	for _, v := range c.book {
		binary.Write(b, binary.LittleEndian, v)
	}
	return hex.Dump(b.Bytes())
}

func (c *tdxCodec) Encode(src []byte) ([]byte, error) {
	var err error
	if len(src)%8 != 0 {
		return src, nil
	}
	dest := bytes.NewBuffer(nil)
	cursor := 0
	for cursor < len(src) {
		b0 := readLittleEndianUint32([4]byte{src[cursor], src[cursor+1], src[cursor+2], src[cursor+3]})
		b1 := readLittleEndianUint32([4]byte{src[cursor+4], src[cursor+5], src[cursor+6], src[cursor+7]})
		b0, b1, err = c.shift(b0, b1)
		if err != nil {
			return nil, err
		}
		dest.Write(writeLittleEndianUint32(b0))
		dest.Write(writeLittleEndianUint32(b1))
		cursor += 8
	}
	return dest.Bytes(), nil
}

func readLittleEndianUint32(b [4]byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
func writeLittleEndianUint32(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b
}

func NewTDXCodec() (*tdxCodec, error) {
	var err error
	e := &tdxCodec{}
	e.book = make([]uint32, 0)
	buf := bytes.NewBuffer(nil)
	for _, s := range []string{encryptSeedPart1HexDump, encryptSeedPart2HexDump} {
		d, err := hex.DecodeString(s)
		if err != nil {
			return nil, err
		}
		buf.Write(d)
	}
	for buf.Len() > 0 {
		v := uint32(0)
		err = binary.Read(buf, binary.LittleEndian, &v)
		if err != nil {
			return nil, err
		}
		e.book = append(e.book, v)
	}

	cursor := 0
	keyCurssor := 0
	for i := 0; i < 18; i += 1 {
		mask := uint32(0)
		for j := 0; j < 4; j += 1 {
			mask = mask << 8
			mask += uint32(encryptKey[keyCurssor])
			keyCurssor = (keyCurssor + 1) % len(encryptKey)
		}
		e.book[cursor] ^= mask
		cursor += 1
	}

	cursor = 0
	var b0 uint32 = 0
	var b1 uint32 = 0
	for i := 0; i < 9; i += 1 {
		b0, b1, err = e.shift(b0, b1)
		if err != nil {
			return nil, err
		}
		e.book[cursor], e.book[cursor+1] = b0, b1
		cursor += 2
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 128; j++ {
			b0, b1, err = e.shift(b0, b1)
			if err != nil {
				return nil, err
			}
			e.book[cursor], e.book[cursor+1] = b0, b1
			cursor += 2
		}
	}

	return e, nil
}

const encryptKey = "SECURE20031107_TDXAB"
const encryptSeedPart1HexDump = "886A3F24D308A3852E8A191344737003223809A4D0319F2998FA2E08896C4EECE62128457713D038CF6654BE6C0CE934B729ACC0DD507CC9B5D5843F170947B5D9D516921BFB7989"
const encryptSeedPart2HexDump = "A60B31D1ACB5DF98DB72FD2FB7DF1AD0EDAFE1B8967E266A45907CBA997F2CF14799A124F76C91B3E2F2010816FC8E85D8206963694E5771A3FE58A47E3D93F48F74950D58B68E7258CD8B71EE4A15821DA4547BB5595AC239D5309C1360F22A23B0D1C5F0856028187941CAEF38DBB8B0DC798E0E183A608B0E9E6C3E8A1EB0C17715D7274B31BDDA2FAF78605C6055F32555E694AB55AA629848574014E8636A39CA55B610AB2A345CCCB4CEE84111AF8654A193E9727C1114EEB32ABC6F635DC5A92BF6311874163E5CCE1E93879B33BAD6AF5CCF246C8153327A7786952898488F3BAFB94B6B1BE8BFC493212866CC09D86191A921FB60AC7C483280EC5D5D5D84EFB17585E9022326DC881B65EB813E8923C5AC96D3F36F6D0F3942F48382440B2E042084A44AF0C8695E9B1F9E4268C6219A6CE9F6619C0C67F088D3ABD2A0516A682F54D828A70F96A33351AB6C0BEF6EE43B7A1350F03BBA982AFB7E1D65F1A17601AF393E59CA66880E43821986EE8CB49F6F45C3A5847DBE5E8B3BD8756FE07320C1859F441A40A66AC15662AAD34E06773F3672DFFE1B3D029B4224D7D03748120AD0D3EA0FDB9BC0F149C97253077B1B9980D879D425F7DEE8F61A50FEE33B4C79B6BDE06C97BA06C004B64FA9C1C4609F40C29E5C5E63246A19AF6FFB68B5536C3EEBB239136FEC523B1F51FC6D2C95309B444581CC09BD5EAF04D0E3BEFD4A33DE07280F66B34B2E1957A8CBC00F74C845395F0BD2DBFBD3B9BDC079550A32601AC600A1D679722C40FE259F67CCA31FFBF8E9A58EF82232DBDF16753C156B61FDC81E502FAB5205ADFAB53D32608723FD487B315382DF003EBB575C9EA08C6FCA2E56871ADB6917DFF6A842D5C3FF7E28C63267AC73554F8CB0275B69C858CABB5DA3FFE1A011F0B8983DFA10B88321FD6CB5FC4A5BD3D12D79E4539A6545F8B6BC498ED29097FB4BDAF2DDE1337ECBA44113FB62E8C6E4CEDACA20EF014C7736FE9E7ED0B41FF12B4DDADB95989190AE718EADEAA0D5936BD0D18ED0E025C7AF2F5B3C8EB794758EFBE2F68F642B12F212B888881CF00D90A05EAD4F1CC38F6891F1CFD1ADC1A8B318222F2F77170EBEFE2D75EAA11F028B0FCCA0E5E8746FB5D6F3AC1899E289CEE04FA8B4B7E013FD813BC47CD9A8ADD266A25F16057795801473CC9377141A216520ADE686FAB577F54254C7CF359DFB0CAFCDEBA0893E7BD31B41D6497E1EAE2D0E25005EB37120BB006822AFE0B8579B3664241EB909F01D916355AAA6DF598943C1787F535AD9A25B7D20C5B9E50276032683A9CF95626819C811414A734ECA2D47B34AA9147B5200511B1529539A3F570FD6E4C69BBC76A4602B0074E681B56FBA081FE91B576BEC96F215D90D2A216563B6B6F9B9E72E0534FF645685C55D2DB053A18F9FA99947BA086A07856EE9707A4B4429B3B52E0975DB232619C4B0A66EAD7DDFA749B860EE9C66B2ED8F718CAAECFF179A696C526456E19EB1C2A5023619294C0975401359A03E3A18E49A98543F659D425BD6E48F6BD63FF799079CD2A1F530E8EFE6382D4DC15D25F08620DD4C26EB7084C6E982635ECC1E023F6B6809C9EFBA3E1418973CA1706A6B84357F6886E2A05205539CB7370750AA1C84073E5CAEDE7FEC447D8EB8F2165737DA3AB00D0C50F0041F1CF0FFB300021AF50CAEB274B53C587A8325BD2109DCF91391D1F62FA97C734732940147F52281E5E53ADCDAC2373476B5C8A7DDF39A466144A90E03D00F3EC7C8EC411E75A499CD38E22F0EEA3BA1BB803231B33E18388B544E08B96D4F030D426FBF040AF69012B82C797C972472B07956AF89AFBC1F779ADE100893D912AE8BB32E3FCFDC1F72125524716B2EE6DD1A5087CD849F1847587A17DA0874BC9A9FBC8C7D4BE93AEC7AECFA1D85DB66430963D2C364C447181CEF08D91532373B43DD16BAC224434DA11251C4652A02009450DDE43A139EF8DF71554E3110D677AC819B19115FF15635046BC7A3D73B18113C09A52459EDE68FF2FAFBF1972CBFBA9E6E3C151E7045E386B16FE9EA0A5E0E86B32A3E5A1CE71F77FA063D4EB9DC65290F1DE799D6893E8025C8665278C94C2E6AB3109CBA0E15C678EAE294533CFCA5F42D0A1EA74EF7F23D2B1D360F2639196079C21908A72352B61213F76EFEADEB661FC3EA9545BCE383C87BA6D1377FB128FF8C01EFDD32C3A55A6CBE852158650298AB680FA5CEEE3B952FDBAD7DEF2A842F6E5B28B62115706107297547DDEC10159F6130A8CC1396BD61EB1EFE3403CF6303AA905C73B539A2704C0B9E9ED514DEAACBBC86CCEEA72C6260AB5CAB9C6E84F3B2AF1E8B64CAF0BD19B96923A050BB5A65325A6840B3B42A3CD5E99E31F7B821C0190B549B99A05F877E99F795A87D3D629A8837F8772DE3975F93ED11811268162988350ED61FE6C7A1DFDE9699BA5878A584F5576372221BFFC3839B9646C21AEB0AB3CD54302E53E448D98F2831BC6DEFF2EB58EAFFC63461ED28FE733C7CEED9144A5DE3B764E8145D1042E0133E20B6E2EE45EAABAAA3154F6CDBD04FCBFA42F442C7B5BB6AEF1D3B4F650521CD419E791ED8C74D85866A474BE45062813DF2A162CF46268D5BA08388FCA3B6C7C1C324157F9274CB690B8A844785B2925600BF5B099D4819AD74B16214000E82232A8D4258EAF5550C3EF4AD1D61703F2392F07233417E938DF1EC5FD6DB3B226C5937DE7C6074EECBA7F285406E3277CE848007A69E50F81955D8EFE83597D961AAA769A9C2060CC5FCAB045ADCCA0B802E7A449E843445C30567D5FDC99E1E0ED3DB73DBCD88551079DA5F67404367E36534C4C5D8383E719EF8283D20FF6DF1E7213E154A3DB08F2B9FE3E6F7AD83DB685A3DE9F74081941C264CF634296994F7201541F7D402762E6BF4BC6800A2D4712408D46AF42033B7D4B743AF6100502EF6391E46452497744F211440888BBF1DFC954DAF91B596D3DDF470452FA066EC09BCBF8597BD03D06DAC7F0485CB31B327EB964139FD55E64725DA9A0ACAAB25785028F4290453DA862C0AFB6DB6E96214DC68006948D7A4C00E68EE8DA127A2FE3F4F8CAD87E806E08CB5B6D6F47A7C1ECEAAEC5F37D399A378CE422A6B40359EFE20B985F3D9ABD739EE8B4E123BF7FAC91D56186D4B3166A326B297E3EA74FA6E3A32435BDDF7E74168FB2078CA4EF50AFB97B3FED8AC564045279548BA3A3A5355878D8320B7A96BFE4B9596D0BC67A855589A15A16329A9CC33DBE199564A2AA6F925313F1C7EF45E7C31299002E8F8FD702F27045C15BB80E32C28054815C195226DC6E43F13C148DC860FC7EEC9F9070F1F0441A4794740176E885DEB515F32D1C09BD58FC1BCF26435114134787B25609C2A60A3E8F8DF1B6C631FC2B4120E9E32E102D14F66AF1581D1CAE095236BE1923E33620B243B22B9BEEE0EA2B285990DBAE68C0C72DE28F7A22D457812D0FD94B79562087D64F0F5CCE76FA34954FA487D8727FD9DC31E8D3EF34163470A74FF2E99AB6E6F3A37FDF8F460DC12A8F8DDEBA14CE11B990D6B6EDB10557BC6372C676D3BD4652704E8D0DCC70D29F1A3FF00CC920F39B50BED0F69FB9F7B669C7DDBCE0BCF91A0A35E15D9882F13BB24AD5B51BF79947BEBD63B76B32E3937795911CC97E226802D312EF4A7AD42683B2B6AC6CC4C75121CF12E783742126AE75192B7E6BBA1065063FB4B18106B1AFAEDCA11D8BD253DC9C3E1E2591642448613120A6EEC0CD92AEAABD54E67AF645FA886DA88E9BFBEFEC3E4645780BC9D86C0F7F0F87B78604D6003604683FDD1B01F38F604AE4577CCFC36D7336B428371AB1EF0874180B05F5E003CBE57A07724AEE8BD99424655612E58BF8FF4584EA2FDDDF238EF74F4C2BD8987C3F96653748EB3C855F275B4B9D9FC466126EB7A84DF1D8B790E6A84E2955F918E596E467057B4209155D58C4CDE02C9E1AC0BB9D00582BB4862A8119EA97475B6197FB709DCA9E0A1092D66334632C4021F5AE88CBEF00925A0994A10FE6E1D1D3DB91ADFA4A50B0FF286A169F1682883DAB7DCFE0639579BCEE2A1527FCD4F015E1150FA8306A7C4B502A027D0E60D278CF89A41863F77064C60C3B506A861287A17F0E086F5C0AA586000627DDC30D79EE61163EA382394DDC2533416C2C256EECBBBDEB6BC90A17DFCEB761D59CE09E4056F88017C4B3D0A7239247C927C5F72E386B99D4D72B45BC11AFCB89ED3785554EDB5A5FC08D37C3DD8C40FAD4D5EEF501EF8E661B1D91485A23C13516CE7C7D56FC44EE156CEBF2A3637C8C6DD34329AD7128263928EFA0E67E000604037CE393ACFF5FAD33777C2AB1B2DC55A9E67B05C4237A34F402782D3BE9BBC999D8E11D515730FBF7E1C2DD67BC400C76B1B8CB74590A121BEB16EB2B46E366A2FAB4857796E94BCD276A3C6C8C24965EEF80F537DDE8D461D0A73D5C64DD04CDBBB39295046BAA9E82695AC04E35EBEF0D5FAA19A512D6AE28CEF6322EE869AB8C289C0F62E2443AA031EA5A4D0F29CBA61C0834D6AE99B5015E58FD65B64BAF9A22628E13A3AA78695A94BE96255EFD3EF2FC7DAF752F7696F043F590AFA7715A9E4800186B087ADE6099B93E53E3B5AFD90E997D7349ED9B7F02C518B2B023AACD5967DA67D01D63ECFD1282D7D7CCF259F1F9BB8F2AD72B4D65A4CF5885A71AC29E0E6A519E0FDACB0479BFA93ED8DC4D3E8CC573B282966D5F8282E137991015F78556075ED440E96F78C5ED3E3D46D0515BA6DF4882561A103BDF06405159EEBC3A257903CEC1A27972A073AA99B6D3F1BF521631EFB669CF519F3DC2628D93375F5FD55B182345603BB3CBA8A11775128F8D90AC26751CCAB5F92ADCC5117E84D8EDC303862589D3791F92093C2907AEACE7B3EFB64CE215132BE4F777EE3B6A8463D29C36953DE4880E613641008AEA224B26DDDFD2D8569662107090A469AB3DDC04564CFDE6C58AEC8201CDDF7BE5B408D581B7F01D2CCBBE3B46B7E6AA2DD45FF593A440A353ED5CDB4BCA8CEEA72BB8464FAAE12668D476F3CBF63E49BD29E5D2F541B77C2AE70634EF68D0D0E7457135BE7711672F85D7D53AF08CB4040CCE2B44E6A46D23484AF15012804B0E11D3A9895B49FB80648A06ECE823B3F6F82AB20354B1D1A01F8277227B1601561DC3F93E72B793ABBBD254534E13988A04B79CE51B7C9322FC9BA1FA07EC81CE0F6D1C7BCC31101CFC7AAE8A14987901A9ABD4FD4CBDEDAD038DA0AD52AC33903673691C67C31F98D4F2BB1E0B7599EF73ABBF543FF19D5F29C45D9272C2297BF2AFCE61571FC910F2515949B6193E5FAEB9CB6CE5964A8C2D1A8BA125E07C1B60C6A05E36550D21042A403CB0E6EECE03BDB9816BEA0984C64E9783232951F9FDF92D3E02B34A0D31EF2718941740A1B8C34A34B2071BEC5D83276C38D9F35DF2E2F999B476F0BE61DF1E30F54DA4CE591D8DA1ECF7962CE6F7E3ECD66B11816051D2CFDC5D28F849922FBF657F323F5237632A63135A89302CDCC566281F0ACB5EB755A9736166ECC73D288926296DED049B9811B90504C1456C671BDC7C6E60A147A3206D0E1459A7BF2C3FD53AAC9000FA862E2BF25BBF6D2BD3505691271220204B27CCFCBB62B9C76CDC03E1153D3E3401660BDAB38F0AD47259C2038BA76CE46F7C5A1AF77606075204EFECB85D88DE88AB0F9AA7A7EAAF94C5CC248198C8AFB02E46AC301F9E1EBD669F8D490A0DE5CA62D25093F9FE608C232614EB75BE277CEE3DF8F57E672C33A"

const (
	keySimpleXOR0547 = 0x93
)

func decryptSimpleXOR(src []byte, key byte) []byte {
	for i := 0; i < len(src); i += 1 {
		src[i] ^= key
	}
	return src
}
