package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/khengsaurus/ng-gql-todos/consts"
	"github.com/khengsaurus/ng-gql-todos/database"
	"github.com/khengsaurus/ng-gql-todos/graph/generated"
	"github.com/khengsaurus/ng-gql-todos/graph/model"
	"github.com/khengsaurus/ng-gql-todos/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, newUser model.NewUser) (*model.User, error) {
	fmt.Println("CreateUser called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	usersColl, err := mongoClient.GetCollection(consts.UsersCollection)
	if err != nil {
		return nil, err
	}

	existCount, _ := usersColl.CountDocuments(ctx, bson.M{"email": newUser.Email})
	if existCount >= 1 {
		return nil, errors.New("user with that email already exists")
	}

	result, err := usersColl.InsertOne(ctx, newUser)
	if err != nil {
		fmt.Printf("Failed to insert document into %s collection", consts.UsersCollection)
	}

	username := ""
	if newUser.Username == nil {
		username = newUser.Email
	} else {
		username = *newUser.Username
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &model.User{
			ID:       oid.Hex(),
			Username: username,
			Email:    &newUser.Email,
		}, err
	}

	return nil, errors.New("failed to create user")
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, userID string) (*bool, error) {
	fmt.Println("DeleteUser called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	usersColl, err := mongoClient.GetCollection(consts.UsersCollection)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: userId}}

	_, err = usersColl.DeleteOne(ctx, filter)
	if err != nil {
		v := false
		return &v, err
	}
	v := true
	database.RemoveKeyFromRedis(ctx, utils.GetUserTodosKey(userID))

	return &v, nil
}

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, newTodo model.NewTodo) (*model.Todo, error) {
	fmt.Println("CreateTodo called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	todosColl, err := mongoClient.GetCollection(consts.TodosCollection)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(newTodo.UserID)
	if err != nil {
		return nil, err
	}

	currTime := time.Now()
	result, err := todosColl.InsertOne(ctx, bson.D{
		{Key: "userId", Value: userId},
		{Key: "text", Value: newTodo.Text},
		{Key: "done", Value: false},
		{Key: "priority", Value: 2},
		{Key: "tag", Value: "white"},
		{Key: "createdAt", Value: currTime},
		{Key: "updatedAt", Value: currTime},
	})

	database.RemoveKeyFromRedis(ctx, utils.GetUserTodosKey(newTodo.UserID))
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &model.Todo{
			ID:       oid.Hex(),
			UserID:   newTodo.UserID,
			Text:     newTodo.Text,
			Priority: 2,
			Tag:      "white",
			Done:     false,
		}, err
	}

	return nil, errors.New("failed to create todo")
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, updateTodo model.UpdateTodo) (string, error) {
	fmt.Println("UpdateTodo called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return "", err
	}
	defer mongoClient.Disconnect(ctx)

	todosColl, err := mongoClient.GetCollection(consts.TodosCollection)
	if err != nil {
		return "", err
	}

	todoId, err := primitive.ObjectIDFromHex(updateTodo.ID)
	if err != nil {
		return "", err
	}

	userId, err := primitive.ObjectIDFromHex(updateTodo.UserID)
	if err != nil {
		return "", err
	}

	filter := bson.D{{Key: "_id", Value: todoId}}
	updateVals := bson.D{
		{Key: "userId", Value: userId},
		{Key: "updatedAt", Value: time.Now()},
	}
	if updateTodo.Text != nil {
		updateVals = append(updateVals, bson.E{Key: "text", Value: updateTodo.Text})
	}
	if updateTodo.Priority != nil {
		updateVals = append(updateVals, bson.E{Key: "priority", Value: updateTodo.Priority})
	}
	if updateTodo.Tag != nil {
		updateVals = append(updateVals, bson.E{Key: "tag", Value: updateTodo.Tag})
	}
	if updateTodo.Done != nil {
		updateVals = append(updateVals, bson.E{Key: "done", Value: updateTodo.Done})
	}
	update := bson.D{{
		Key:   "$set",
		Value: updateVals,
	}}

	_, err = todosColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}
	database.RemoveKeyFromRedis(ctx, utils.GetUserTodosKey(updateTodo.UserID))

	return updateTodo.ID, nil
}

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, userID string, todoID string) (string, error) {
	fmt.Println("DeleteTodo called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return "", err
	}
	defer mongoClient.Disconnect(ctx)

	todosColl, err := mongoClient.GetCollection(consts.TodosCollection)
	if err != nil {
		return "", err
	}

	todoId, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return "", err
	}
	filter := bson.D{{Key: "_id", Value: todoId}}

	_, err = todosColl.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}
	database.RemoveKeyFromRedis(ctx, utils.GetUserTodosKey(userID))

	return todoID, nil
}

// CreateBoard is the resolver for the createBoard field.
func (r *mutationResolver) CreateBoard(ctx context.Context, newBoard model.NewBoard) (*model.Board, error) {
	fmt.Println("CreateBoard called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	boardsColl, err := mongoClient.GetCollection(consts.BoardsCollection)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(newBoard.UserID)
	if err != nil {
		return nil, err
	}

	todos := make([]*primitive.ObjectID, 0)
	for _, s := range newBoard.TodoIds {
		todoId, err := primitive.ObjectIDFromHex(*s)
		if err == nil {
			todos = append(todos, &todoId)
		}
	}

	currTime := time.Now()
	result, err := boardsColl.InsertOne(ctx, bson.D{
		{Key: "userId", Value: userId},
		{Key: "name", Value: newBoard.Name},
		{Key: "todos", Value: todos},
		{Key: "createdAt", Value: currTime},
		{Key: "updatedAt", Value: currTime},
	})

	database.RemoveKeyFromRedis(ctx, utils.GetUserBoardsKey(newBoard.UserID))
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &model.Board{
			ID:      oid.Hex(),
			Name:    newBoard.Name,
			UserID:  newBoard.UserID,
			TodoIds: newBoard.TodoIds,
		}, err
	}

	return nil, errors.New("failed to create board")
}

// UpdateBoard is the resolver for the updateBoard field.
func (r *mutationResolver) UpdateBoard(ctx context.Context, updateBoard model.UpdateBoard) (string, error) {
	fmt.Println("UpdateBoard called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return "", err
	}
	defer mongoClient.Disconnect(ctx)

	boardsColl, err := mongoClient.GetCollection(consts.BoardsCollection)
	if err != nil {
		return "", err
	}

	boardId, err := primitive.ObjectIDFromHex(updateBoard.ID)
	if err != nil {
		return "", err
	}

	filter := bson.D{{Key: "_id", Value: boardId}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "todos", Value: updateBoard.Todos},
			{Key: "name", Value: updateBoard.Name},
			{Key: "updatedAt", Value: time.Now()},
		},
	}}

	_, err = boardsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}
	database.RemoveKeyFromRedis(ctx, utils.GetUserTodosKey(updateBoard.UserID))

	return updateBoard.ID, nil
}

// DeleteBoard is the resolver for the deleteBoard field.
func (r *mutationResolver) DeleteBoard(ctx context.Context, userID string, boardID string) (string, error) {
	fmt.Println("DeleteBoard called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return "", err
	}
	defer mongoClient.Disconnect(ctx)

	boardsColl, err := mongoClient.GetCollection(consts.BoardsCollection)
	if err != nil {
		return "", err
	}

	boardId, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return "", err
	}
	filter := bson.D{{Key: "_id", Value: boardId}}

	_, err = boardsColl.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}
	database.RemoveKeyFromRedis(ctx, utils.GetUserTodosKey(userID))

	return boardID, nil
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, email string) (*model.User, error) {
	fmt.Println("GetUser called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	usersColl, err := mongoClient.GetCollection(consts.UsersCollection)
	if err != nil {
		return nil, err
	}

	result := usersColl.FindOne(ctx, bson.M{"email": email})
	var user model.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context) ([]*model.User, error) {
	fmt.Println("GetUsers called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	defer mongoClient.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	usersColl, err := mongoClient.GetCollection(consts.UsersCollection)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(10)
	cur, err := usersColl.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var users []*model.User
	for cur.Next(ctx) {
		var user model.User
		err := cur.Decode(&user)
		if err != nil {
			fmt.Println("Failed to decode user document")
		} else {
			users = append(users, &user)
		}
	}

	return users, nil
}

// GetTodo is the resolver for the getTodo field.
func (r *queryResolver) GetTodo(ctx context.Context, todoID string) (*model.Todo, error) {
	fmt.Println("GetTodo called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	todosColl, err := mongoClient.GetCollection(consts.TodosCollection)
	if err != nil {
		return nil, err
	}

	todoId, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return nil, err
	}

	result := todosColl.FindOne(ctx, bson.M{"_id": todoId})
	var todo model.Todo
	if err := result.Decode(&todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

// GetTodos is the resolver for the getTodos field.
func (r *queryResolver) GetTodos(ctx context.Context, userID string, fresh bool) ([]*model.Todo, error) {
	fmt.Println("GetTodos called")
	redisClient, redisClientErr := database.GetRedisClient(ctx)
	if !fresh && redisClient != nil {
		cachedTodos, _ := redisClient.GetTodos(ctx, userID)
		if cachedTodos != nil {
			fmt.Println("Retrieved todos from redis cache")
			return cachedTodos, nil
		}
	}

	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	todosColl, err := mongoClient.GetCollection(consts.TodosCollection)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	filter := bson.D{{Key: "userId", Value: userId}}
	cur, err := todosColl.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var todos []*model.Todo
	for cur.Next(ctx) {
		var todo model.Todo
		err := cur.Decode(&todo)
		if err != nil {
			fmt.Println("Failed to decode todo document:")
			fmt.Println(fmt.Printf("%v", err))
		} else {
			todos = append(todos, &todo)
		}
	}

	if redisClient != nil && redisClientErr == nil {
		redisClient.SetTodos(ctx, userID, todos)
	}

	return todos, nil
}

// GetBoard is the resolver for the getBoard field.
// FIXME: returning [] for todoIds
func (r *queryResolver) GetBoard(ctx context.Context, boardID string) (*model.Board, error) {
	fmt.Println("GetBoard called")
	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	boardsColl, err := mongoClient.GetCollection(consts.BoardsCollection)
	if err != nil {
		return nil, err
	}

	todoId, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return nil, err
	}

	result := boardsColl.FindOne(ctx, bson.M{"_id": todoId})
	var board model.Board
	if err := result.Decode(&board); err != nil {
		return nil, err
	}
	return &board, nil
}

// GetBoards is the resolver for the getBoards field.
// FIXME: returning [] for todoIds
func (r *queryResolver) GetBoards(ctx context.Context, userID string, fresh bool) ([]*model.Board, error) {
	fmt.Println("GetBoards called")
	redisClient, redisClientErr := database.GetRedisClient(ctx)
	if !fresh && redisClient != nil {
		cachedBoards, _ := redisClient.GetBoards(ctx, userID)
		if cachedBoards != nil {
			fmt.Println("Retrieved boards from redis cache")
			return cachedBoards, nil
		}
	}

	mongoClient, err := database.GetMongoClient(ctx, true)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Disconnect(ctx)

	boardsColl, err := mongoClient.GetCollection(consts.BoardsCollection)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// findOptions := options.Find()
	// findOptions.SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	filter := bson.D{{Key: "userId", Value: userId}}
	cur, err := boardsColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// result := todosColl.FindOne(ctx, bson.M{"_id": todoId})
	// var todo model.Todo
	// if err := result.Decode(&todo); err != nil {
	// 	return nil, err
	// }
	// return &todo, nil

	boards := make([]*model.Board, 0)
	for cur.Next(ctx) {
		var board model.Board
		err := cur.Decode(&board)
		if err != nil {
			fmt.Println("Failed to decode board document:")
			fmt.Println(fmt.Printf("%v", err))
		} else {
			boards = append(boards, &board)
		}
	}

	if redisClient != nil && redisClientErr == nil {
		redisClient.SetBoards(ctx, userID, boards)
	}

	return boards, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
