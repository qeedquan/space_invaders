package main

import (
	"bytes"
	"fmt"
)

const (
	B = iota
	C
	D
	E
	H
	L
	F
	A
)

var opcyc = [256]uint64{
	//  0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, // 0
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, // 1
	4, 10, 16, 5, 5, 5, 7, 4, 4, 10, 16, 5, 5, 5, 7, 4, // 2
	4, 10, 13, 5, 10, 10, 10, 4, 4, 10, 13, 5, 5, 5, 7, 4, // 3
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 4
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 5
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 6
	7, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 7, 5, // 7
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // 8
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // 9
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // A
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // B
	5, 10, 10, 10, 11, 11, 7, 11, 5, 10, 10, 10, 11, 11, 7, 11, // C
	5, 10, 10, 10, 11, 11, 7, 11, 5, 10, 10, 10, 11, 11, 7, 11, // D
	5, 10, 10, 18, 11, 11, 7, 11, 5, 5, 10, 5, 11, 11, 7, 11, // E
	5, 10, 10, 4, 11, 11, 7, 11, 5, 5, 10, 4, 11, 11, 7, 11, // F
}

var opstr = [256]string{
	"nop", "lxi b,#", "stax b", "inx b", "inr b", "dcr b", "mvi b,#", "rlc",
	"ill", "dad b", "ldax b", "dcx b", "inr c", "dcr c", "mvi c,#", "rrc",
	"ill", "lxi d,#", "stax d", "inx d", "inr d", "dcr d", "mvi d,#", "ral",
	"ill", "dad d", "ldax d", "dcx d", "inr e", "dcr e", "mvi e,#", "rar",
	"ill", "lxi h,#", "shld", "inx h", "inr h", "dcr h", "mvi h,#", "daa",
	"ill", "dad h", "lhld", "dcx h", "inr l", "dcr l", "mvi l,#", "cma",
	"ill", "lxi sp,#", "sta $", "inx sp", "inr M", "dcr M", "mvi M,#", "stc",
	"ill", "dad sp", "lda $", "dcx sp", "inr a", "dcr a", "mvi a,#", "cmc",
	"mov b,b", "mov b,c", "mov b,d", "mov b,e", "mov b,h", "mov b,l",
	"mov b,M", "mov b,a", "mov c,b", "mov c,c", "mov c,d", "mov c,e",
	"mov c,h", "mov c,l", "mov c,M", "mov c,a", "mov d,b", "mov d,c",
	"mov d,d", "mov d,e", "mov d,h", "mov d,l", "mov d,M", "mov d,a",
	"mov e,b", "mov e,c", "mov e,d", "mov e,e", "mov e,h", "mov e,l",
	"mov e,M", "mov e,a", "mov h,b", "mov h,c", "mov h,d", "mov h,e",
	"mov h,h", "mov h,l", "mov h,M", "mov h,a", "mov l,b", "mov l,c",
	"mov l,d", "mov l,e", "mov l,h", "mov l,l", "mov l,M", "mov l,a",
	"mov M,b", "mov M,c", "mov M,d", "mov M,e", "mov M,h", "mov M,l", "hlt",
	"mov M,a", "mov a,b", "mov a,c", "mov a,d", "mov a,e", "mov a,h",
	"mov a,l", "mov a,M", "mov a,a", "add b", "add c", "add d", "add e",
	"add h", "add l", "add M", "add a", "adc b", "adc c", "adc d", "adc e",
	"adc h", "adc l", "adc M", "adc a", "sub b", "sub c", "sub d", "sub e",
	"sub h", "sub l", "sub M", "sub a", "sbb b", "sbb c", "sbb d", "sbb e",
	"sbb h", "sbb l", "sbb M", "sbb a", "ana b", "ana c", "ana d", "ana e",
	"ana h", "ana l", "ana M", "ana a", "xra b", "xra c", "xra d", "xra e",
	"xra h", "xra l", "xra M", "xra a", "ora b", "ora c", "ora d", "ora e",
	"ora h", "ora l", "ora M", "ora a", "cmp b", "cmp c", "cmp d", "cmp e",
	"cmp h", "cmp l", "cmp M", "cmp a", "rnz", "pop b", "jnz $", "jmp $",
	"cnz $", "push b", "adi #", "rst 0", "rz", "ret", "jz $", "ill", "cz $",
	"call $", "aci #", "rst 1", "rnc", "pop d", "jnc $", "out p", "cnc $",
	"push d", "sui #", "rst 2", "rc", "ill", "jc $", "in p", "cc $", "ill",
	"sbi #", "rst 3", "rpo", "pop h", "jpo $", "xthl", "cpo $", "push h",
	"ani #", "rst 4", "rpe", "pchl", "jpe $", "xchg", "cpe $", "ill", "xri #",
	"rst 5", "rp", "pop psw", "jp $", "di", "cp $", "push psw", "ori #",
	"rst 6", "rm", "sphl", "jm $", "ei", "cm $", "ill", "cpi #", "rst 7",
}

type IO interface {
	In(uint8)
	Out(uint8)
}

type CPU struct {
	IO
	R               [8]uint8
	M               [0x10000]uint8
	PC              uint16
	SP              uint16
	CY, HC, Z, S, P bool
	IF, HLT         bool
	Cycles          uint64
}

func NewCPU() *CPU {
	return &CPU{}
}

func (c *CPU) Reset() {
	c.R = [8]uint8{}
	c.PC = 0
	c.SP = 0
	c.CY = false
	c.HC = false
	c.Z = false
	c.S = false
	c.P = false
	c.IF = false
	c.HLT = false
	c.Cycles = 0
}

func (c *CPU) Step() error {
	op := c.fetch8()
	c.Cycles += opcyc[op]

	switch op {
	// nop
	case 0x00, 0x10, 0x20, 0x30,
		0x08, 0x18, 0x28, 0x38:

	// mov r, r
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x47,
		0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4f,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x57,
		0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5f,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x67,
		0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6f,
		0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7f:
		d := (op >> 3) & 0x7
		s := op & 0x7
		c.R[d] = c.R[s]

	// mov r, m
	case 0x46, 0x4e, 0x56, 0x5e, 0x66, 0x6e, 0x7e:
		d := (op >> 3) & 0x7
		a := c.hl()
		c.R[d] = c.M[a]
	// mov m, r
	case 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77:
		s := op & 0x7
		a := c.hl()
		c.M[a] = c.R[s]

	// mvi r
	case 0x3e, 0x06, 0x0e, 0x16, 0x1e, 0x26, 0x2e:
		d := (op >> 3) & 0x7
		v := c.fetch8()
		c.R[d] = v
	// mvi m
	case 0x36:
		a := c.hl()
		v := c.fetch8()
		c.M[a] = v

	// ldax b
	case 0x0a:
		c.R[A] = c.Read8(c.bc())
	// ldax d
	case 0x1a:
		c.R[A] = c.Read8(c.de())
	// lda w
	case 0x3a:
		c.R[A] = c.Read8(c.fetch16())

	// stax b
	case 0x02:
		c.Write8(c.bc(), c.R[A])
	// stax d
	case 0x12:
		c.Write8(c.de(), c.R[A])
	// sta w
	case 0x32:
		c.Write8(c.fetch16(), c.R[A])

	// lxi b, w
	case 0x01:
		c.setbc(c.fetch16())
	// lxi d, w
	case 0x11:
		c.setde(c.fetch16())
	// lxi h, w
	case 0x21:
		c.sethl(c.fetch16())
	// lxi sp, w
	case 0x31:
		c.SP = c.fetch16()
	// lhld
	case 0x2a:
		a := c.fetch16()
		v := c.Read16(a)
		c.sethl(v)
	// shld
	case 0x22:
		a := c.fetch16()
		v := c.hl()
		c.Write16(a, v)
	// sphl
	case 0xf9:
		c.SP = c.hl()

	// xchg
	case 0xeb:
		c.xchg()
	// xthl
	case 0xe3:
		c.xthl()

	// add r
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x87:
		d := op & 0x7
		c.add(A, c.R[d], 0)
	// add m
	case 0x86:
		c.add(A, c.Read8(c.hl()), 0)
	// adi
	case 0xc6:
		c.add(A, c.fetch8(), 0)

	// adc r
	case 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8f:
		d := op & 0x7
		c.add(A, c.R[d], truth(c.CY))
	// adc m
	case 0x8e:
		c.add(A, c.Read8(c.hl()), truth(c.CY))
	// aci
	case 0xce:
		c.add(A, c.fetch8(), truth(c.CY))

	// sub r
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x97:
		d := op & 0x7
		c.sub(A, c.R[d], 0)
	// sub m
	case 0x96:
		c.sub(A, c.Read8(c.hl()), 0)
	// sui
	case 0xd6:
		c.sub(A, c.fetch8(), 0)

	// sbb r
	case 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9f:
		d := op & 0x7
		c.sub(A, c.R[d], truth(c.CY))
	// sbb m
	case 0x9e:
		c.sub(A, c.Read8(c.hl()), truth(c.CY))
	// sbi
	case 0xde:
		c.sub(A, c.fetch8(), truth(c.CY))

	// dad bc
	case 0x09:
		c.dad(c.bc())
	// dad de
	case 0x19:
		c.dad(c.de())
	// dad hl
	case 0x29:
		c.dad(c.hl())
	// dad sp
	case 0x39:
		c.dad(c.SP)

	// di
	case 0xf3:
		c.IF = false
	// ei
	case 0xfb:
		c.IF = true
	// hlt
	case 0x76:
		c.PC--
		c.HLT = true

	// inr r
	case 0x4, 0xc, 0x14, 0x1c, 0x24, 0x2c, 0x3c:
		d := (op >> 3) & 0x7
		c.R[d] = c.inr(c.R[d])
	// inr m
	case 0x34:
		a := c.hl()
		v := c.Read8(a)
		c.Write8(a, c.inr(v))

	// dcr r
	case 0x05, 0x0d, 0x15, 0x1d, 0x25, 0x2d, 0x3d:
		d := (op >> 3) & 0x7
		c.R[d] = c.dcr(c.R[d])
	// dcr m
	case 0x35:
		a := c.hl()
		v := c.Read8(a)
		c.Write8(a, c.dcr(v))

	// inx b
	case 0x03:
		c.setbc(c.bc() + 1)
	// inx d
	case 0x13:
		c.setde(c.de() + 1)
	// inx h
	case 0x23:
		c.sethl(c.hl() + 1)
	// inx sp
	case 0x33:
		c.SP++

	// dcx b
	case 0x0b:
		c.setbc(c.bc() - 1)
	// dcx d
	case 0x1b:
		c.setde(c.de() - 1)
	// dcx h
	case 0x2b:
		c.sethl(c.hl() - 1)
	// dcx sp
	case 0x3b:
		c.SP--

	// daa
	case 0x27:
		c.daa()
	// cma
	case 0x2f:
		c.R[A] ^= 0xff
	// stc
	case 0x37:
		c.CY = true
	// cmc
	case 0x3f:
		c.CY = !c.CY

	// rlc
	case 0x07:
		c.rlc()
	// rrc
	case 0x0f:
		c.rrc()
	// ral
	case 0x17:
		c.ral()
	// rar
	case 0x1f:
		c.rar()

	// ana r
	case 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa7:
		d := op & 0x7
		c.ana(c.R[d])
	// ana m
	case 0xa6:
		c.ana(c.Read8(c.hl()))
	// ani
	case 0xe6:
		c.ana(c.fetch8())

	// xra r
	case 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xaf:
		d := op & 0x7
		c.xra(c.R[d])
	// xra m
	case 0xae:
		c.xra(c.Read8(c.hl()))
	// xri
	case 0xee:
		c.xra(c.fetch8())

	// ora r
	case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb7:
		d := op & 0x7
		c.ora(c.R[d])
	// ora m
	case 0xb6:
		c.ora(c.Read8(c.hl()))
	// ori
	case 0xf6:
		c.ora(c.fetch8())

	// cmp r
	case 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbf:
		d := op & 0x7
		c.cmp(c.R[d])
	// cmp m
	case 0xbe:
		c.cmp(c.Read8(c.hl()))
	// cpi
	case 0xfe:
		c.cmp(c.fetch8())

	// jmp
	case 0xc3:
		c.jmp(c.fetch16())
	// jnz
	case 0xc2:
		c.condjmp(!c.Z)
	// jz
	case 0xca:
		c.condjmp(c.Z)
	// jnc
	case 0xd2:
		c.condjmp(!c.CY)
	// jc
	case 0xda:
		c.condjmp(c.CY)
	// jpo
	case 0xe2:
		c.condjmp(!c.P)
	// jpe
	case 0xea:
		c.condjmp(c.P)
	// jp
	case 0xf2:
		c.condjmp(!c.S)
	// jm
	case 0xfa:
		c.condjmp(c.S)

	// pchl
	case 0xe9:
		c.PC = c.hl()
	// call
	case 0xcd:
		c.call(c.fetch16())

	// cnz
	case 0xc4:
		c.condcall(!c.Z)
	// cz
	case 0xcc:
		c.condcall(c.Z)
	// cnc
	case 0xd4:
		c.condcall(!c.CY)
	// cc
	case 0xdc:
		c.condcall(c.CY)
	// cpo
	case 0xe4:
		c.condcall(!c.P)
	// cpe
	case 0xec:
		c.condcall(c.P)
	// cp
	case 0xf4:
		c.condcall(!c.S)
	// cm
	case 0xfc:
		c.condcall(c.S)

	// ret
	case 0xc9, 0xd9:
		c.ret()
	// rnz
	case 0xc0:
		c.condret(!c.Z)
	// rz
	case 0xc8:
		c.condret(c.Z)
	// rnc
	case 0xd0:
		c.condret(!c.CY)
	// rc
	case 0xd8:
		c.condret(c.CY)
	// rpo
	case 0xe0:
		c.condret(!c.P)
	// rpe
	case 0xe8:
		c.condret(c.P)
	// rp
	case 0xf0:
		c.condret(!c.S)
	// rm
	case 0xf8:
		c.condret(c.S)

	// rst 0..7
	case 0xc7, 0xcf, 0xd7, 0xdf, 0xe7, 0xef, 0xf7, 0xff:
		a := uint16(op & 0x38)
		c.call(a)

	// push b
	case 0xc5:
		c.push16(c.bc())
	// push d
	case 0xd5:
		c.push16(c.de())
	// push h
	case 0xe5:
		c.push16(c.hl())
	// push psw
	case 0xf5:
		c.pushpsw()
	// pop b
	case 0xc1:
		c.setbc(c.pop16())
	// pop d
	case 0xd1:
		c.setde(c.pop16())
	// pop h
	case 0xe1:
		c.sethl(c.pop16())
	// pop psw
	case 0xf1:
		c.poppsw()

	// in
	case 0xdb:
		p := c.fetch8()
		if c.IO != nil {
			c.In(p)
		}
	// out
	case 0xd3:
		p := c.fetch8()
		if c.IO != nil {
			c.Out(p)
		}

	default:
		return fmt.Errorf("unknown opcode %#x", op)
	}
	return nil
}

func (c *CPU) Read8(a uint16) uint8 {
	return c.M[a]
}

func (c *CPU) Read16(a uint16) uint16 {
	lo := c.Read8(a)
	hi := c.Read8(a + 1)
	return uint16(lo) | uint16(hi)<<8
}

func (c *CPU) Write8(a uint16, v uint8) {
	c.M[a] = v
}

func (c *CPU) Write16(a, v uint16) {
	c.Write8(a, uint8(v))
	c.Write8(a+1, uint8(v>>8))
}

func (c *CPU) fetch8() uint8 {
	v := c.Read8(c.PC)
	c.PC++
	return v
}

func (c *CPU) fetch16() uint16 {
	lo := c.fetch8()
	hi := c.fetch8()
	return uint16(lo) | uint16(hi)<<8
}

func (c *CPU) bc() uint16 {
	return uint16(c.R[B])<<8 | uint16(c.R[C])
}

func (c *CPU) de() uint16 {
	return uint16(c.R[D])<<8 | uint16(c.R[E])
}

func (c *CPU) hl() uint16 {
	return uint16(c.R[H])<<8 | uint16(c.R[L])
}

func (c *CPU) setbc(v uint16) {
	c.R[B] = uint8(v >> 8)
	c.R[C] = uint8(v)
}

func (c *CPU) setde(v uint16) {
	c.R[D] = uint8(v >> 8)
	c.R[E] = uint8(v)
}

func (c *CPU) sethl(v uint16) {
	c.R[H] = uint8(v >> 8)
	c.R[L] = uint8(v)
}

func (c *CPU) xchg() {
	de, hl := c.de(), c.hl()
	c.setde(hl)
	c.sethl(de)
}

func (c *CPU) xthl() {
	ds, hl := c.Read16(c.SP), c.hl()
	c.Write16(c.SP, hl)
	c.sethl(ds)
}

func (c *CPU) add(d, v, cy uint8) {
	r := int16(c.R[d]) + int16(v) + int16(cy)
	c.Z = r&0xff == 0
	c.S = (r & 0x80) != 0
	c.CY = (r & 0x100) != 0
	c.HC = (c.R[d]^uint8(r)^v)&0x10 != 0
	c.P = parity(uint8(r))
	c.R[d] = uint8(r)
}

func (c *CPU) sub(d, v, cy uint8) {
	r := int16(c.R[d]) - int16(v) - int16(cy)
	c.Z = r&0xff == 0
	c.S = (r & 0x80) != 0
	c.CY = (r & 0x100) != 0
	c.HC = ^(c.R[d]^uint8(r)^v)&0x10 != 0
	c.P = parity(uint8(r))
	c.R[d] = uint8(r)
}

func (c *CPU) dad(v uint16) {
	r := uint32(c.hl()) + uint32(v)
	c.sethl(uint16(r))
	c.CY = r&0x10000 != 0
}

func (c *CPU) inr(v uint8) uint8 {
	r := v + 1
	c.HC = r&0xf == 0
	c.Z = r == 0
	c.S = r&0x80 != 0
	c.P = parity(r)
	return r
}

func (c *CPU) dcr(v uint8) uint8 {
	r := v - 1
	c.HC = !(r&0xf == 0xf)
	c.Z = r == 0
	c.S = r&0x80 != 0
	c.P = parity(r)
	return r
}

func (c *CPU) ana(v uint8) {
	r := c.R[A] & v
	c.CY = false
	c.HC = ((c.R[A] | v) & 0x08) != 0
	c.Z = r == 0
	c.S = r&0x80 != 0
	c.P = parity(r)
	c.R[A] = r
}

func (c *CPU) xra(v uint8) {
	c.R[A] ^= v
	c.CY = false
	c.HC = false
	c.Z = c.R[A] == 0
	c.S = c.R[A]&0x80 != 0
	c.P = parity(c.R[A])
}

func (c *CPU) ora(v uint8) {
	c.R[A] |= v
	c.CY = false
	c.HC = false
	c.Z = c.R[A] == 0
	c.S = c.R[A]&0x80 != 0
	c.P = parity(c.R[A])
}

func (c *CPU) cmp(v uint8) {
	r := int16(c.R[A]) - int16(v)
	c.CY = r&0x100 != 0
	c.HC = ^(c.R[A]^uint8(r)^v)&0x10 != 0
	c.Z = r&0xff == 0
	c.S = r&0x80 != 0
	c.P = parity(uint8(r))
}

func (c *CPU) jmp(a uint16) {
	c.PC = a
}

func (c *CPU) condjmp(cond bool) {
	a := c.fetch16()
	if cond {
		c.PC = a
	}
}

func (c *CPU) call(a uint16) {
	c.push16(c.PC)
	c.jmp(a)
}

func (c *CPU) condcall(cond bool) {
	a := c.fetch16()
	if cond {
		c.call(a)
		c.Cycles += 6
	}
}

func (c *CPU) ret() {
	c.PC = c.pop16()
}

func (c *CPU) condret(cond bool) {
	if cond {
		c.ret()
		c.Cycles += 6
	}
}

func (c *CPU) pushpsw() {
	psw := uint8(0x2)
	if c.S {
		psw |= 0x80
	}
	if c.Z {
		psw |= 0x40
	}
	if c.HC {
		psw |= 0x10
	}
	if c.P {
		psw |= 0x04
	}
	if c.CY {
		psw |= 0x01
	}

	c.push16(uint16(c.R[A])<<8 | uint16(psw))
}

func (c *CPU) poppsw() {
	af := c.pop16()
	c.R[A] = uint8(af >> 8)

	psw := af & 0xff
	c.S = psw&0x80 != 0
	c.Z = psw&0x40 != 0
	c.HC = psw&0x10 != 0
	c.P = psw&0x04 != 0
	c.CY = psw&0x01 != 0
}

func (c *CPU) rlc() {
	c.CY = c.R[A]&0x80 != 0
	c.R[A] <<= 1
	if c.CY {
		c.R[A] |= 0x01
	}
}

func (c *CPU) rrc() {
	c.CY = c.R[A]&0x1 != 0
	c.R[A] >>= 1
	if c.CY {
		c.R[A] |= 0x80
	}
}

func (c *CPU) ral() {
	cy := c.CY
	c.CY = c.R[A]&0x80 != 0
	c.R[A] <<= 1
	if cy {
		c.R[A] |= 0x01
	}
}

func (c *CPU) rar() {
	cy := c.CY
	c.CY = c.R[A]&0x1 != 0
	c.R[A] >>= 1
	if cy {
		c.R[A] |= 0x80
	}
}

func (c *CPU) daa() {
	cy := c.CY
	v := uint8(0)

	lsb := c.R[A] & 0xf
	msb := c.R[A] >> 4
	if c.HC || lsb > 9 {
		v += 0x6
	}
	if c.CY || msb > 9 || (msb >= 9 && lsb > 9) {
		v += 0x60
		cy = true
	}
	c.add(A, v, 0)
	c.P = parity(c.R[A])
	c.CY = cy
}

func (c *CPU) push16(v uint16) {
	c.SP -= 2
	c.Write16(c.SP, v)
}

func (c *CPU) pop16() uint16 {
	v := c.Read16(c.SP)
	c.SP += 2
	return v
}

func parity(v uint8) bool {
	p := true
	for v != 0 {
		p = !p
		v &= (v - 1)
	}
	return p
}

func truth(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (c *CPU) Interrupt(a uint16) {
	if c.IF {
		c.IF = false
		c.call(a)
		c.Cycles += 11
	}
}

func (c *CPU) Disasm(pc uint16) string {
	var flags = []byte("......")
	if c.Z {
		flags[0] = 'z'
	}
	if c.S {
		flags[1] = 's'
	}
	if c.P {
		flags[2] = 'p'
	}
	if c.IF {
		flags[3] = 'i'
	}
	if c.CY {
		flags[4] = 'c'
	}
	if c.HC {
		flags[5] = 'a'
	}

	b := new(bytes.Buffer)
	fmt.Fprintf(b, "af\tbc\tde\thl\tpc\tsp\tflags\tcycles\n")
	fmt.Fprintf(b, "%02X__\t%04X\t%04X\t%04X\t%04X\t%04X\t%s\t%d\n",
		c.R[A], c.bc(), c.de(), c.hl(), c.PC, c.SP, flags, c.Cycles)

	fmt.Fprintf(b, "%04X: ", pc)
	fmt.Fprintf(b, "%02X %02X %02X", c.Read8(pc), c.Read8(pc+1), c.Read8(pc+2))

	fmt.Fprintf(b, " - %s", opstr[c.Read8(pc)])

	fmt.Fprintf(b, "\n================================")
	fmt.Fprintf(b, "==============================\n")

	return b.String()
}
