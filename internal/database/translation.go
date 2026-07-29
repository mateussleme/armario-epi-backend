package database

import (
	"context"

	"github.com/TopSisErp/epi-backend/internal/viaonda"
	"github.com/lib/pq"
)

func TagsToProducts(ctx context.Context, tags []viaonda.TagEntry) (map[string]string, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tagIds := []string{}
	for _, tag := range tags {
		tagIds = append(tagIds, tag.Id)
	}

	rows, err := conn.QueryContext(ctx, `
		select
			tagId,
			produto
		from epitag
		where
			epitag.tagId = ANY($1)
	`, pq.Array(tagIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := map[string]string{}
	for rows.Next() {
		tagId := ""
		product := ""

		err := rows.Scan(&tagId, &product)
		if err != nil {
			return nil, err
		}

		products[tagId] = product
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
