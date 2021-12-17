package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	fmt.Printf("Part 1 solution: %d\n", part1())

	fmt.Printf("Part 2 solution: %d\n", part2())
}

const (
	input = "420D50000B318100415919B24E72D6509AE67F87195A3CCC518CC01197D538C3E00BC9A349A09802D258CC16FC016100660DC4283200087C6485F1C8C015A00A5A5FB19C363F2FD8CE1B1B99DE81D00C9D3002100B58002AB5400D50038008DA2020A9C00F300248065A4016B4C00810028003D9600CA4C0084007B8400A0002AA6F68440274080331D20C4300004323CC32830200D42A85D1BE4F1C1440072E4630F2CCD624206008CC5B3E3AB00580010E8710862F0803D06E10C65000946442A631EC2EC30926A600D2A583653BE2D98BFE3820975787C600A680252AC9354FFE8CD23BE1E180253548D057002429794BD4759794BD4709AEDAFF0530043003511006E24C4685A00087C428811EE7FD8BBC1805D28C73C93262526CB36AC600DCB9649334A23900AA9257963FEF17D8028200DC608A71B80010A8D50C23E9802B37AA40EA801CD96EDA25B39593BB002A33F72D9AD959802525BCD6D36CC00D580010A86D1761F080311AE32C73500224E3BCD6D0AE5600024F92F654E5F6132B49979802129DC6593401591389CA62A4840101C9064A34499E4A1B180276008CDEFA0D37BE834F6F11B13900923E008CF6611BC65BCB2CB46B3A779D4C998A848DED30F0014288010A8451062B980311C21BC7C20042A2846782A400834916CFA5B8013374F6A33973C532F071000B565F47F15A526273BB129B6D9985680680111C728FD339BDBD8F03980230A6C0119774999A09001093E34600A60052B2B1D7EF60C958EBF7B074D7AF4928CD6BA5A40208E002F935E855AE68EE56F3ED271E6B44460084AB55002572F3289B78600A6647D1E5F6871BE5E598099006512207600BCDCBCFD23CE463678100467680D27BAE920804119DBFA96E05F00431269D255DDA528D83A577285B91BCCB4802AB95A5C9B001299793FCD24C5D600BC652523D82D3FCB56EF737F045008E0FCDC7DAE40B64F7F799F3981F2490"
)

func parseBinary(bits []byte) int64 {
	if len(bits) > 64 {
		log.Panicf("can't parse %s", string(bits))
	}
	result := int64(0)
	for i := 0; i < len(bits); i++ {
		if bits[len(bits)-i-1] == '1' {
			result += 1 << i
		}
	}
	//	log.Printf("parsed %s as %d", string(bits), result)
	return result
}

var (
	hexCharToBits = map[byte]string{
		'0': "0000",
		'1': "0001",
		'2': "0010",
		'3': "0011",
		'4': "0100",
		'5': "0101",
		'6': "0110",
		'7': "0111",
		'8': "1000",
		'9': "1001",
		'A': "1010",
		'B': "1011",
		'C': "1100",
		'D': "1101",
		'E': "1110",
		'F': "1111",
	}
)

func hexToBits(chars []byte) []byte {
	res := make([]byte, len(chars)*4)
	for i, c := range chars {
		copy(res[i*4:(i+1)*4], hexCharToBits[c])
	}
	return res
}

const (
	PacketLiteralType = 4
)

type Packet struct {
	Version int64
	Type    int64

	Literal    int64
	Subpackets []*Packet
}

func getNextBitsOrDie(cursor int, arr []byte, length int) []byte {
	if cursor+length > len(arr) {
		log.Fatalf("exhausted byte stream; wanted %d bits starting from %d, but have only %d", length, cursor, len(arr))
	}
	return arr[cursor : cursor+length]
}

func parsePacket(bits []byte) (*Packet, int, error) {
	packetCursor := 0

	version := parseBinary(getNextBitsOrDie(packetCursor, bits, 3))
	log.Printf("see version %d (%s)", version, bits[packetCursor:packetCursor+3])
	packetCursor += 3

	packetType := parseBinary(getNextBitsOrDie(packetCursor, bits, 3))
	log.Printf("see type %d (%s)", packetType, bits[packetCursor:packetCursor+3])
	packetCursor += 3

	if packetType == PacketLiteralType {
		log.Printf("see literal")
		startOfLiteral := packetCursor
		literalValue := int64(0)
		for {
			nibble := bits[packetCursor : packetCursor+5]
			packetCursor += 5

			literalValue <<= 4
			literalValue += parseBinary(nibble[1:])

			if nibble[0] == '0' {
				break
			}
		}
		log.Printf("parsed literal %d (%s)", literalValue, bits[startOfLiteral:packetCursor])
		return &Packet{
			Version: version,
			Type:    PacketLiteralType,
			Literal: literalValue,
		}, packetCursor, nil

	} else {
		lengthType := parseBinary(getNextBitsOrDie(packetCursor, bits, 1))
		log.Printf("Got length type %d (%s)", lengthType, bits[packetCursor:packetCursor+1])
		packetCursor += 1

		totalSubPacketLength := int64(-1)
		subPacketCount := int64(-1)
		switch lengthType {
		case 0:
			totalSubPacketLength = parseBinary(getNextBitsOrDie(packetCursor, bits, 15))
			packetCursor += 15
			log.Printf("Looking for subpackets totaling length %d (%s)", totalSubPacketLength, string(bits[packetCursor:packetCursor+15]))
		case 1:
			subPacketCount = parseBinary(getNextBitsOrDie(packetCursor, bits, 11))
			log.Printf("Looking for %d subpackets (%s)", subPacketCount, string(bits[packetCursor:packetCursor+11]))
			packetCursor += 11
		default:
			log.Fatalf("unknown subpacket length type %d", lengthType)

		}

		log.Printf("Descending...")
		subPackets := []*Packet{}
		for subPacketBitsConsumed := int64(0); (lengthType == 0 && subPacketBitsConsumed < totalSubPacketLength) ||
			(lengthType == 1 && int64(len(subPackets)) < subPacketCount); {

			subPacket, subPacketLength, err := parsePacket(bits[packetCursor:])
			if err != nil {
				log.Fatal(err)
			}

			packetCursor += subPacketLength
			subPacketBitsConsumed += int64(subPacketLength)
			subPackets = append(subPackets, subPacket)
			log.Printf("got a subpacket of length %d", subPacketLength)
		}
		log.Printf("Ascending...")

		return &Packet{
			Version:    version,
			Type:       packetType,
			Subpackets: subPackets,
		}, packetCursor, nil
	}
}

func sumVersions(p *Packet) int64 {
	total := int64(0)
	total += p.Version
	for _, subP := range p.Subpackets {
		total += sumVersions(subP)
	}
	return total
}

func part1() int64 {
	packets := []*Packet{}
	bitsRemaining := hexToBits([]byte(input))
	log.Println(string(bitsRemaining))
	for len(bitsRemaining) > 6 {
		nextPacket, bitsConsumed, err := parsePacket(bitsRemaining)
		if err != nil {
			log.Fatalf("failed to parse packet %v", err)
		}
		packets = append(packets, nextPacket)
		bitsRemaining = bitsRemaining[bitsConsumed:]

		log.Printf("got packet: %+v", nextPacket)
	}
	log.Printf("trailing bits: %s", bitsRemaining)

	versionTotal := int64(0)
	for _, p := range packets {
		versionTotal += sumVersions(p)
	}

	return versionTotal
}

func evalPacket(p *Packet) int64 {
	switch p.Type {
	case 0: // sum
		result := int64(0)
		for _, s := range p.Subpackets {
			result += evalPacket(s)
		}
		return result

	case 1: // product
		result := int64(1)
		for _, s := range p.Subpackets {
			result *= evalPacket(s)
		}
		return result

	case 2: // minimum
		result := int64(math.MaxInt64)
		for _, s := range p.Subpackets {
			subResult := evalPacket(s)
			if subResult < result {
				result = subResult
			}
		}
		return result

	case 3: // maximum
		result := int64(math.MinInt64)
		for _, s := range p.Subpackets {
			subResult := evalPacket(s)
			if subResult > result {
				result = subResult
			}
		}
		return result

	case 4: // literal
		return p.Literal

	case 5: // greater-than
		if len(p.Subpackets) != 2 {
			log.Fatalf("gt packet with %d subpackets", len(p.Subpackets))
		}
		if evalPacket(p.Subpackets[0]) > evalPacket(p.Subpackets[1]) {
			return 1
		} else {
			return 0
		}

	case 6: // less-than
		if len(p.Subpackets) != 2 {
			log.Fatalf("lt packet with %d subpackets", len(p.Subpackets))
		}
		if evalPacket(p.Subpackets[0]) < evalPacket(p.Subpackets[1]) {
			return 1
		} else {
			return 0
		}

	case 7: // equal
		if len(p.Subpackets) != 2 {
			log.Fatalf("eq packet with %d subpackets", len(p.Subpackets))
		}
		if evalPacket(p.Subpackets[0]) == evalPacket(p.Subpackets[1]) {
			return 1
		} else {
			return 0
		}

	default:
		log.Fatalf("unknown op type %d", p.Type)
	}

	return -1
}

func part2() int64 {
	packets := []*Packet{}
	bitsRemaining := hexToBits([]byte(input))
	log.Println(string(bitsRemaining))
	for len(bitsRemaining) > 6 {
		nextPacket, bitsConsumed, err := parsePacket(bitsRemaining)
		if err != nil {
			log.Fatalf("failed to parse packet %v", err)
		}
		packets = append(packets, nextPacket)
		bitsRemaining = bitsRemaining[bitsConsumed:]

		log.Printf("got packet: %+v", nextPacket)
	}
	log.Printf("trailing bits: %s", bitsRemaining)

	return evalPacket(packets[0])
}
