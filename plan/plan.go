// Copyright 2014 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"github.com/pingcap/tidb/context"
	"github.com/pingcap/tidb/expression"
	"github.com/pingcap/tidb/field"
	"github.com/pingcap/tidb/table"
	"github.com/pingcap/tidb/util/format"
)

// RowIterFunc is the callback for iterating records.
type RowIterFunc func(id interface{}, data []interface{}) (more bool, err error)

// Plan is the interface of query execution plan.
type Plan interface {
	// Explain the plan.
	Explain(w format.Formatter)
	// GetFields returns the result field list for a plan.
	GetFields() []*field.ResultField
	// Filter try to use index plan to reduce the result set.
	// If index can be used, a new index plan is returned, 'filtered' is true.
	// If no index can be used, the original plan is returned and 'filtered' return false.
	Filter(ctx context.Context, expr expression.Expression) (p Plan, filtered bool, err error)

	// Next returns the next row that contains data and row keys, nil row means there is no more to return.
	// Aggregation plan will fetch all the data at the first call.
	Next(ctx context.Context) (row *Row, err error)

	// Close closes the underlying iterator of the plan.
	// If you call Next after Close, it will start iteration from the beginning.
	// If the plan is not returned as Recordset.Plan, it must be Closed to prevent resource leak.
	Close() error
}

// Planner is implemented by any structure that has a Plan method.
type Planner interface {
	// Plan function returns Plan.
	Plan(ctx context.Context) (Plan, error)
}

// Row represents a record row.
type Row struct {
	// Data is the output record data for current Plan.
	Data []interface{}
	// FromData is the first origin record data, generated by From.
	FromData []interface{}

	RowKeys []*RowKeyEntry
}

// RowKeyEntry is designed for Delete statement in multi-table mode,
// we should know which table this row comes from.
type RowKeyEntry struct {
	// The table which this row come from.
	Tbl table.Table
	// Row handle.
	Key string
}
