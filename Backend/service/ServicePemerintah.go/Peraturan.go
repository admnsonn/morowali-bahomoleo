package servicepemerintahgo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func Peraturan(c *gin.Context) {
	type Visi_desa_kontainer struct {
		Peraturan string `json:"peraturan"`
	}

	id := c.Param("id")

	ctx := context.Background()
	tx, err := DBConnect.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(err.Error())
	}
	defer tx.Rollback(ctx)

	query := `
		select peraturan from dev.desa where id_desa = $1
	`

	row, err := tx.Query(ctx, query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		err = tx.Commit(ctx)
		if err != nil {
			panic(err.Error())
		}
		return
	}

	defer row.Close()

	var Tampung []Visi_desa_kontainer

	for row.Next() {
		var ambil Visi_desa_kontainer

		err := row.Scan(
			&ambil.Peraturan,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
			err = tx.Commit(ctx)
			if err != nil {
				panic(err.Error())
			}
			return
		}

		Tampung = append(Tampung, ambil)
	}

	if len(Tampung) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data Ada!",
			"data":    Tampung,
		})
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Data Tidak Ada!",
			"data":    []Visi_desa_kontainer{},
		})
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
}
