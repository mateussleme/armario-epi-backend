package database

import (
	"context"

	"github.com/TopSisErp/epi-backend/internal/viaonda"
	"github.com/lib/pq"
)

func UpdateWithTags(ctx context.Context, tags []viaonda.TagEntry) error {
	products, err := TagsToProducts(ctx, tags)
	if err != nil {
		return err
	}

	stock, err := GetStock(ctx)
	if err != nil {
		return err
	}
	counts := ProductsToCount(products)

	for product := range stock {
		if _, ok := counts[product]; !ok {
			counts[product] = 0
		}
	}

	tagList := []string{}
	for _, tag := range tags {
		tagList = append(tagList, tag.Id)
	}

	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for product, qty := range counts {
		_, err := tx.ExecContext(ctx, `
			update epiproduto
			set
				quantidade = $2
			where
				epiproduto.id = $1
		`, product, qty)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			delete from epitag
			where
				epitag.produto = $1
				and epitag.tagId != ANY($2)
		`, product, pq.Array(tagList))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
