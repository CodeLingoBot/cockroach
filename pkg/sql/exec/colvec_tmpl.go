// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

// {{/*
// +build execgen_template
//
// This file is the execgen template for colvec.eg.go. It's formatted in a
// special way, so it's both valid Go and a valid text/template input. This
// permits editing this file with editor support.
//
// */}}

package exec

import (
	"fmt"

	"github.com/cockroachdb/apd"
	"github.com/cockroachdb/cockroach/pkg/sql/exec/types"
)

// {{/*

// Dummy import to pull in "apd" package.
var _ apd.Decimal

// _TYPES_T is the template type variable for types.T. It will be replaced by
// types.Foo for each type Foo in the types.T type.
const _TYPES_T = types.Unhandled

// */}}

func (m *memColumn) Append(vec ColVec, colType types.T, toLength uint64, fromLength uint16) {
	switch colType {
	// {{range .}}
	case _TYPES_T:
		m.col = append(m._TemplateType()[:toLength], vec._TemplateType()[:fromLength]...)
		// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", colType))
	}

	if fromLength > 0 {
		m.nulls = append(m.nulls, make([]int64, (fromLength-1)>>6+1)...)

		if vec.HasNulls() {
			for i := uint16(0); i < fromLength; i++ {
				if vec.NullAt(i) {
					m.SetNull64(toLength + uint64(i))
				}
			}
		}
	}
}

func (m *memColumn) AppendWithSel(
	vec ColVec, sel []uint16, batchSize uint16, colType types.T, toLength uint64,
) {
	switch colType {
	// {{range .}}
	case _TYPES_T:
		toCol := append(m._TemplateType()[:toLength], make([]_GOTYPE, batchSize)...)
		fromCol := vec._TemplateType()

		for i := uint16(0); i < batchSize; i++ {
			toCol[uint64(i)+toLength] = fromCol[sel[i]]
		}

		m.col = toCol
		// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", colType))
	}

	if batchSize > 0 {
		m.nulls = append(m.nulls, make([]int64, (batchSize-1)>>6+1)...)
		for i := uint16(0); i < batchSize; i++ {
			if vec.NullAt(sel[i]) {
				m.SetNull64(toLength + uint64(i))
			}
		}
	}
}

func (m *memColumn) Copy(src ColVec, srcStartIdx, srcEndIdx uint64, typ types.T) {
	switch typ {
	// {{range .}}
	case _TYPES_T:
		copy(m._TemplateType(), src._TemplateType()[srcStartIdx:srcEndIdx])
		// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", typ))
	}
}

func (m *memColumn) CopyWithSelInt64(vec ColVec, sel []uint64, nSel uint16, colType types.T) {
	m.UnsetNulls()

	// todo (changangela): handle the case when nSel > ColBatchSize
	switch colType {
	// {{range .}}
	case _TYPES_T:
		toCol := m._TemplateType()
		fromCol := vec._TemplateType()

		if vec.HasNulls() {
			for i := uint16(0); i < nSel; i++ {
				if vec.NullAt64(sel[i]) {
					m.SetNull(i)
				} else {
					toCol[i] = fromCol[sel[i]]
				}
			}
		} else {
			for i := uint16(0); i < nSel; i++ {
				toCol[i] = fromCol[sel[i]]
			}
		}
		// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", colType))
	}
}

func (m *memColumn) CopyWithSelInt16(vec ColVec, sel []uint16, nSel uint16, colType types.T) {
	m.UnsetNulls()

	switch colType {
	// {{range .}}
	case _TYPES_T:
		toCol := m._TemplateType()
		fromCol := vec._TemplateType()

		if vec.HasNulls() {
			for i := uint16(0); i < nSel; i++ {
				if vec.NullAt(sel[i]) {
					m.SetNull(i)
				} else {
					toCol[i] = fromCol[sel[i]]
				}
			}
		} else {
			for i := uint16(0); i < nSel; i++ {
				toCol[i] = fromCol[sel[i]]
			}
		}
		// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", colType))
	}
}

func (m *memColumn) CopyWithSelAndNilsInt64(
	vec ColVec, sel []uint64, nSel uint16, nils []bool, colType types.T,
) {
	m.UnsetNulls()

	switch colType {
	// {{range .}}
	case _TYPES_T:
		toCol := m._TemplateType()
		fromCol := vec._TemplateType()

		if vec.HasNulls() {
			// TODO(jordan): copy the null arrays in batch.
			for i := uint16(0); i < nSel; i++ {
				if nils[i] {
					m.SetNull(i)
				} else {
					if vec.NullAt64(sel[i]) {
						m.SetNull(i)
					} else {
						toCol[i] = fromCol[sel[i]]
					}
				}
			}
		} else {
			for i := uint16(0); i < nSel; i++ {
				if nils[i] {
					m.SetNull(i)
				} else {
					toCol[i] = fromCol[sel[i]]
				}
			}
		}
	// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", colType))
	}
}

func (m *memColumn) PrettyValueAt(colIdx uint16, colType types.T) string {
	if m.NullAt(colIdx) {
		return "NULL"
	}
	switch colType {
	// {{range .}}
	case _TYPES_T:
		return fmt.Sprintf("%v", m._TemplateType()[colIdx])
	// {{end}}
	default:
		panic(fmt.Sprintf("unhandled type %d", colType))
	}
}
