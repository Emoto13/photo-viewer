package image_store

const (
	addImageToDatabase = `WITH rows AS (
						  INSERT INTO images(username, path)
						  VALUES ($1, $2)
						  RETURNING id AS new_image_id)
						  INSERT INTO posts(username, image_id, name)
						  SELECT $1, new_image_id, $3 FROM rows;`
)

/*
`WITH rows AS (
	INSERT INTO images(username, path)
	VALUES ('Emoto13', 'https://photo-viewer.s3.eu-west-1.amazonaws.com/22699e233e40705502be537894d5cf5c5a8b0383c758c697ce642cd17b70c7b4.png')
	RETURNING id AS new_image_id)
	INSERT INTO posts(username, image_id, name)
	SELECT 'Emoto13', new_image_id, 'Winner' FROM rows;`
)*/
