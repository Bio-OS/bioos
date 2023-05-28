//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"math"
	"strings"

	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

const (
	defaultSize = 10
	defaultPage = 1
)

// Pagination ...
type Pagination struct {
	Size   int `validate:"gte=0,lte=100"`
	Page   int `validate:"gte=1"`
	Orders []Order
}

// Order ...
type Order struct {
	Field     string
	Ascending bool
}

// NewPagination returns Pagination of size and page.
func NewPagination(size int, page int) *Pagination {
	if size == 0 {
		return &Pagination{Size: defaultSize, Page: defaultPage}
	}
	return &Pagination{Size: size, Page: page}
}

// SetOrderBy Set order by.
func (q *Pagination) SetOrderBy(orderByQuery string) error {
	if len(orderByQuery) == 0 {
		return nil
	}
	orders := strings.Split(orderByQuery, consts.QuerySliceDelimiter)
	q.Orders = make([]Order, 0, len(orders))
	for _, order := range orders {
		o, err := newOrder(order)
		if err != nil {
			return err
		}
		if o != nil {
			q.Orders = append(q.Orders, *o)
		}
	}
	return nil
}

// GetOffset Get offset.
func (q *Pagination) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit Get limit.
func (q *Pagination) GetLimit() int {
	return q.Size
}

// GetOrderBy Get orderBy string.
func (q *Pagination) GetOrderBy() string {
	orders := make([]string, 0, len(q.Orders))
	for i := range q.Orders {
		orders = append(orders, orderString(q.Orders[i]))
	}
	return strings.Join(orders, consts.QuerySliceDelimiter)
}

// GetPage Get OrderBy.
func (q *Pagination) GetPage() int {
	return q.Page
}

// GetSize Get OrderBy.
func (q *Pagination) GetSize() int {
	return q.Size
}

// GetQueryString get query string.
func (q *Pagination) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

// GetTotalPages Get total pages int.
func (q *Pagination) GetTotalPages(totalCount int) int {
	d := float64(totalCount) / float64(q.GetSize())
	return int(math.Ceil(d))
}

// GetHasMore Get has more.
func (q *Pagination) GetHasMore(totalCount int) bool {
	return q.GetPage() < totalCount/q.GetSize()
}

func newOrder(order string) (*Order, error) {
	// length will be 0, 1, 2
	orderInfo := strings.SplitN(order, consts.OrderDelimiter, 2)

	if len(orderInfo) == 0 {
		return nil, nil
	}

	o := &Order{Field: orderInfo[0], Ascending: true}
	if len(orderInfo) == 1 {
		return o, nil
	}

	switch orderInfo[1] {
	case consts.ASCOrdering:
		o.Ascending = true
	case consts.DESCOrdering:
		o.Ascending = false
	default:
		return nil, apperrors.NewInvalidError("order")
	}
	return o, nil
}

func orderString(order Order) string {
	if order.Ascending {
		return strings.Join([]string{order.Field, consts.ASCOrdering}, consts.OrderDelimiter)
	}
	return strings.Join([]string{order.Field, consts.DESCOrdering}, consts.OrderDelimiter)
}
