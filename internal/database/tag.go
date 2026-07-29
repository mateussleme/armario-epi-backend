package database

import (
	"context"
	"slices"

	"github.com/TopSisErp/epi-backend/internal/viaonda"
	"github.com/lib/pq"
)

func UnknownTags(ctx context.Context, tags []viaonda.TagEntry) ([]viaonda.TagEntry, error) {
	products, err := TagsToProducts(ctx, tags)
	if err != nil {
		return nil, err
	}

	newTags := []viaonda.TagEntry{}
	for _, tag := range tags {
		if _, ok := products[tag.Epc]; !ok {
			newTags = append(newTags, tag)
		}
	}

	return newTags, nil
}

func InventoryWithTags(ctx context.Context, tags []viaonda.TagEntry, product string) error {
	newTags := []string{}
	for _, tag := range tags {
		if !slices.Contains(newTags, tag.Epc) {
			newTags = append(newTags, tag.Epc)
		}
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

	_, err = tx.ExecContext(ctx, `
		insert into epitag (produto, tagEpc)
		select 
			$1, 
			t.tagEpc
		from unnest($2::text[]) as t(tagEpc);
	`, product, pq.Array(newTags))
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		update epiproduto
		set
			quantidade = epiproduto.quantidade + $1
		where
			epiproduto.id = $2
	`, len(newTags), product)
	if err != nil {
		return err
	}

	return tx.Commit()
}
