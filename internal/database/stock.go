package database

import "context"

func GetStock(ctx context.Context) (map[string]int, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			id,
			quantidade
		from epiproduto
		order by id
	`)
	if err != nil {
		return nil, err
	}

	productCount := map[string]int{}
	defer rows.Close()
	for rows.Next() {
		id := ""
		quantity := 0

		err := rows.Scan(&id, &quantity)
		if err != nil {
			return nil, err
		}

		productCount[id] = quantity
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productCount, nil
}
