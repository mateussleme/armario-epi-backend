package database

import (
	"context"
	"slices"

	"github.com/TopSisErp/epi-backend/internal/viaonda"
	"github.com/lib/pq"
)

func TagsToProducts(ctx context.Context, tags []viaonda.TagEntry) (map[string]string, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tagEpcs := []string{}
	for _, tag := range tags {
		if !slices.Contains(tagEpcs, tag.Epc) {
			tagEpcs = append(tagEpcs, tag.Epc)
		}
	}

	rows, err := conn.QueryContext(ctx, `
		select
			tagEpc,
			produto
		from epitag
		where
			epitag.tagEpc = ANY($1)
	`, pq.Array(tagEpcs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := map[string]string{}
	for rows.Next() {
		tagEpc := ""
		product := ""

		err := rows.Scan(&tagEpc, &product)
		if err != nil {
			return nil, err
		}

		products[tagEpc] = product
	}

	return products, nil
}

func ProductsToCount(products map[string]string) map[string]int {
	productCount := map[string]int{}
	for _, product := range products {
		if _, ok := productCount[product]; !ok {
			productCount[product] = 0
		}

		productCount[product] = productCount[product] + 1
	}

	return productCount
}
