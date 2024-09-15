CREATE TABLE IF NOT EXISTS sale_item_additions (
    sale_id UUID NOT NULL,
    item_id INTEGER NOT NULL,
    addition_id UUID NOT NULL,
    PRIMARY KEY (sale_id, item_id, addition_id),
    FOREIGN KEY (sale_id, item_id) REFERENCES sale_items(sale_id, item_id),
    FOREIGN KEY (addition_id) REFERENCES additions(id)
);
