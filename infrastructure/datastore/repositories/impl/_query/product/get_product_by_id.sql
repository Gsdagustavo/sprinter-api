SELECT id,
       name,
       description,
       price,
       stock,
       image_url
FROM products
WHERE id = ? AND status_code = 0