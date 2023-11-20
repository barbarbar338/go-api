package server

import (
	"go-api/pkg/entities"
	"go-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdatePostBody struct {
	Title       string `json:"title" validate:"required,alpha,min=2,max=30"`
	Description string `json:"description" validate:"required,min=20,max=200"`
}

type UpdatePostResponse struct {
	Post *entities.Post `json:"post"`
}

func (s *Server) editPost(c *fiber.Ctx) error {
	id := c.Params("id")

	body := new(UpdatePostBody)
	if err := c.BodyParser(body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	post := new(entities.Post)

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "title", Value: body.Title},
				{Key: "description", Value: body.Description},
			},
		},
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res := s.posts.FindOneAndUpdate(c.Context(), bson.D{{Key: "_id", Value: oid}}, update)

	err = res.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewHttpError(c, errors.NOT_FOUND, err.Error())
		}
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res.Decode(post)

	return c.JSON(UpdatePostResponse{post})
}
