package storage

import (
	"context"
	"github.com/Calevin/go_palantir/ent"
	"github.com/Calevin/go_palantir/ent/file"
	"log"
)

func SaveFile(fileName string, client *ent.Client, ctx context.Context) *ent.File {
	// Buscamos si ya existe un File con ese nombre, si no, lo creamos.
	f, err := client.File.Query().Where(file.NameEQ(fileName)).Only(ctx)
	if err != nil {
		// Si no se encuentra, se crea uno.
		f, err = client.File.
			Create().
			SetName(fileName).
			Save(ctx)
		if err != nil {
			log.Printf("Error creando File para %s: %v", fileName, err)
			return nil
		}
	}
	return f
}

func SaveTokens(tokens []*ent.Token, client *ent.Client, f *ent.File, ctx context.Context, path string) error {
	// Preparamos un slice de creadores para insertar los tokens en bloque.
	var bulkCreates []*ent.TokenCreate
	for _, t := range tokens {
		// Cada "t" es una instancia de ent.Token
		tc := client.Token.Create().
			SetLine(t.Line).
			SetOrder(t.Order).
			SetToken(t.Token).
			SetFile(f) // Aquí se establece la FK hacia la entidad File.
		bulkCreates = append(bulkCreates, tc)
	}

	// Si hay tokens para insertar, se realiza la inserción bulk.
	if len(bulkCreates) > 0 {
		_, err := client.Token.CreateBulk(bulkCreates...).Save(ctx)
		if err != nil {
			log.Printf("Error insertando tokens para el archivo %s: %v", path, err)
			return err
		} else {
			log.Printf("Se insertaron %d tokens del archivo %s", len(bulkCreates), path)
		}
	}

	return nil
}
