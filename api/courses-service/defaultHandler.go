package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (svc *Service) CreateCourse(c *fiber.Ctx) error {
	var course Course

	if err := c.BodyParser(&course); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Invalid RequestBody"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	course.CreateAt = time.Now()

	result, err := svc.Db.Database("lms").Collection("courses").InsertOne(ctx, course)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while Inserting"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course Inserted Successfully", "result": result})

}

func (svc *Service) GetCourses(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := svc.Db.Database("lms").Collection("courses")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while Fetching"})
	}
	defer cursor.Close(ctx)

	var courses []bson.M

	if err := cursor.All(ctx, &courses); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": courses})
}

func (svc *Service) GetOneCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	objId, err := ConvertToObjectId(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Object Id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := svc.Db.Database(GetEnv("mongo.db")).Collection("courses")

	filter := bson.M{
		"_id": objId,
	}
	var course Course
	if err := coll.FindOne(ctx, filter).Decode(&course); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while Fetching"})
	}
	return c.JSON(course)
}

func (svc *Service) UpdateCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	objId, err := ConvertToObjectId(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Object Id"})
	}
	var updatedData Course
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request Body"})
	}

	filter := bson.M{
		"_id": objId,
	}

	update := bson.M{
		"$set": updatedData,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := svc.Db.Database(GetEnv("mongo.db")).Collection("courses")

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while updating"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No Document Found with given id"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "course updated successfully"})
}

func (svc *Service) DeleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	objId, err := ConvertToObjectId(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Object Id"})
	}

	filter := bson.M{
		"_id": objId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := svc.Db.Database(GetEnv("mongo.db")).Collection("courses").DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while deleting"})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No Document Found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted Successfully", "result": result})
}
