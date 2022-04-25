package notes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Title string
}

type UpdateRequest struct {
	Content string
}

func getCreateHandler(ctx context.Context, ns NotesStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createRequest CreateRequest
		c.BindJSON(&createRequest)

		note := Note{
			Username: "shnoo", // TODO: add middleware to get this
			Title:    createRequest.Title,
			Content:  "",
		}

		noteResponse, err := ns.Create(ctx, note)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to create note")
			return
		}

		c.JSON(http.StatusOK, noteResponse)
	}
}

func getListHandler(ctx context.Context, ns NotesStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := "shnoo" // TODO: add middleware to get this

		notesResponseList, err := ns.List(ctx, username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to list notes")
			return
		}

		c.JSON(http.StatusOK, notesResponseList)
	}
}

func getReadHandler(ctx context.Context, ns NotesStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		content, err := ns.Read(ctx, id)

		if err == ErrNoteDoesNotExist {
			c.JSON(http.StatusBadRequest, "Note does not exist")
			return
		}

		if err == ErrInvalidId {
			c.JSON(http.StatusBadRequest, "Invalid id")
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to delete note")
			return
		}

		c.JSON(http.StatusOK, content)
	}
}

func getUpdateHandler(ctx context.Context, ns NotesStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var updateRequest UpdateRequest
		c.BindJSON(&updateRequest)

		err := ns.Update(ctx, id, updateRequest.Content)

		if err == ErrNoteDoesNotExist {
			c.JSON(http.StatusBadRequest, "Note does not exist") // TODO: DRY these if statements
			return
		}

		if err == ErrInvalidId {
			c.JSON(http.StatusBadRequest, "Invalid id")
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to delete note")
			return
		}

		c.JSON(http.StatusOK, "Updated note")
	}
}

func getDeleteHandler(ctx context.Context, ns NotesStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := ns.Delete(ctx, id) // TODO add username

		if err == ErrNoteDoesNotExist {
			c.JSON(http.StatusBadRequest, "Note does not exist")
			return
		}

		if err == ErrInvalidId {
			c.JSON(http.StatusBadRequest, "Invalid id")
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to delete note")
			return
		}

		c.JSON(http.StatusOK, "Deleted note")
	}
}

func CreateRoutes(router *gin.RouterGroup, ctx context.Context, ns NotesStore) {
	router.POST("/create", getCreateHandler(ctx, ns))
	router.GET("/list", getListHandler(ctx, ns))
	router.GET("/read/:id", getReadHandler(ctx, ns))
	router.POST("/update/:id", getUpdateHandler(ctx, ns))
	router.POST("/delete/:id", getDeleteHandler(ctx, ns))
}