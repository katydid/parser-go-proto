package prototests

func NewPopulatedBigMsg(r randyMsg) *BigMsg {
	this := &BigMsg{}
	if r.Intn(5) != 0 {
		v1 := int64(r.Int63())
		if r.Intn(2) == 0 {
			v1 *= -1
		}
		this.Field = &v1
	}
	if r.Intn(5) != 0 {
		this.Msg = NewPopulatedSmallMsg(r)
	}
	return this
}

func NewPopulatedSmallMsg(r randyMsg) *SmallMsg {
	this := &SmallMsg{}
	if r.Intn(5) != 0 {
		v2 := string(randStringMsg(r))
		this.ScarBusStop = &v2
	}
	if r.Intn(5) != 0 {
		v3 := r.Intn(10)
		this.FlightParachute = make([]uint32, v3)
		for i := 0; i < v3; i++ {
			this.FlightParachute[i] = uint32(r.Uint32())
		}
	}
	if r.Intn(5) != 0 {
		v4 := string(randStringMsg(r))
		this.MapShark = &v4
	}
	return this
}

func NewPopulatedPacked(r randyMsg) *Packed {
	this := &Packed{}
	if r.Intn(5) != 0 {
		v5 := r.Intn(10)
		this.Ints = make([]int64, v5)
		for i := 0; i < v5; i++ {
			this.Ints[i] = int64(r.Int63())
			if r.Intn(2) == 0 {
				this.Ints[i] *= -1
			}
		}
	}
	if r.Intn(5) != 0 {
		v6 := r.Intn(10)
		this.Floats = make([]float64, v6)
		for i := 0; i < v6; i++ {
			this.Floats[i] = float64(r.Float64())
			if r.Intn(2) == 0 {
				this.Floats[i] *= -1
			}
		}
	}
	if r.Intn(5) != 0 {
		v7 := r.Intn(10)
		this.Uints = make([]uint32, v7)
		for i := 0; i < v7; i++ {
			this.Uints[i] = uint32(r.Uint32())
		}
	}
	return this
}

type randyMsg interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneMsg(r randyMsg) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringMsg(r randyMsg) string {
	v8 := r.Intn(100)
	tmps := make([]rune, v8)
	for i := 0; i < v8; i++ {
		tmps[i] = randUTF8RuneMsg(r)
	}
	return string(tmps)
}
func randUnrecognizedMsg(r randyMsg, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldMsg(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldMsg(dAtA []byte, r randyMsg, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateMsg(dAtA, uint64(key))
		v9 := r.Int63()
		if r.Intn(2) == 0 {
			v9 *= -1
		}
		dAtA = encodeVarintPopulateMsg(dAtA, uint64(v9))
	case 1:
		dAtA = encodeVarintPopulateMsg(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateMsg(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateMsg(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateMsg(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateMsg(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
