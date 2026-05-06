SELECT
    id,
    name,
    description,
    price,
    stock,
    image_url
FROM products
WHERE status_code = 0
