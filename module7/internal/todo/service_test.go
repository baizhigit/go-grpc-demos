package todo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	store_mock "github.com/baizhigit/go-grpc-demos/module7/internal/mocks"
	todostore "github.com/baizhigit/go-grpc-demos/module7/internal/store"
	"github.com/baizhigit/go-grpc-demos/module7/internal/todo"
	"github.com/baizhigit/go-grpc-demos/module7/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewService(t *testing.T) {
	t.Run("returns error if todo store is nil", func(t *testing.T) {
		service, err := todo.NewService(nil)

		require.Error(t, err)
		assert.Nil(t, service)
	})

	t.Run("successfully initialises todo grpc service", func(t *testing.T) {
		mockStore := &store_mock.MockTaskStore{}

		service, err := todo.NewService(mockStore)

		require.NoError(t, err)
		require.NotNil(t, service)
	})
}

func newTestService(t *testing.T) (proto.TodoServiceServer, *store_mock.MockTaskStore) {
	store := new(store_mock.MockTaskStore)

	service, err := todo.NewService(store)
	require.NoError(t, err)

	t.Cleanup(func() {
		store.AssertExpectations(t)
	})

	return service, store
}

func TestService_AddTask(t *testing.T) {

	t.Run("returns INVALID_ARGUMENT status code when task is empty", func(t *testing.T) {
		service, _ := newTestService(t)

		res, err := service.AddTask(context.Background(), &proto.AddTaskRequest{Task: ""})

		require.Error(t, err)
		require.Nil(t, res)

		statusErr, ok := status.FromError(err)
		require.True(t, ok)

		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
		assert.Equal(t, "task cannot be empty", statusErr.Message())
	})

	t.Run("returns INTERNAL status code when an error is returned from the store", func(t *testing.T) {
		service, todoStore := newTestService(t)

		const task = "wake up"
		testErr := errors.New("some error")

		todoStore.On("AddTask", task).
			Return("", testErr).
			Once()

		res, err := service.AddTask(context.Background(), &proto.AddTaskRequest{Task: task})
		require.Error(t, err)
		require.Nil(t, res)

		statusErr, ok := status.FromError(err)
		require.True(t, ok)

		assert.Equal(t, codes.Internal, statusErr.Code())
		assert.Equal(t, fmt.Sprintf("failed to add task: %v", testErr), statusErr.Message())
	})

	t.Run("returns task ID when task is stored successfully", func(t *testing.T) {
		service, todoStore := newTestService(t)

		const (
			task   = "wake up"
			taskID = "some task id"
		)

		todoStore.On("AddTask", task).
			Return(taskID, nil).
			Once()

		res, err := service.AddTask(context.Background(), &proto.AddTaskRequest{Task: task})
		require.NoError(t, err)
		require.NotNil(t, res)

		assert.Equal(t, taskID, res.GetId())
	})
}

func TestService_CompleteTask(t *testing.T) {
	t.Run("returns NOT_FOUND status code when task is not found", func(t *testing.T) {
		service, todoStore := newTestService(t)

		const taskID = "some task id"

		todoStore.On("CompleteTask", taskID).
			Return(todostore.ErrTaskNotFound).
			Once()

		res, err := service.CompleteTask(context.Background(), &proto.CompleteTaskRequest{Id: taskID})
		require.Error(t, err)
		require.Nil(t, res)

		statusErr, ok := status.FromError(err)
		require.True(t, ok)

		assert.Equal(t, codes.NotFound, statusErr.Code())
		assert.Equal(t, "task not found", statusErr.Message())
	})

	t.Run("returns INTERNAL status code when an error is returned from the store", func(t *testing.T) {
		service, todoStore := newTestService(t)

		const taskID = "some task id"
		testErr := errors.New("some error")

		todoStore.On("CompleteTask", taskID).
			Return(testErr).
			Once()

		res, err := service.CompleteTask(context.Background(), &proto.CompleteTaskRequest{Id: taskID})
		require.Error(t, err)
		require.Nil(t, res)

		statusErr, ok := status.FromError(err)
		require.True(t, ok)

		assert.Equal(t, codes.Internal, statusErr.Code())
		assert.Equal(t, fmt.Sprintf("failed to complete task: %v", testErr), statusErr.Message())
	})

	t.Run("returns successful response when a task is completed", func(t *testing.T) {
		service, todoStore := newTestService(t)

		const taskID = "some task id"

		todoStore.On("CompleteTask", taskID).
			Return(nil).
			Once()

		res, err := service.CompleteTask(context.Background(), &proto.CompleteTaskRequest{Id: taskID})
		require.NoError(t, err)
		require.NotNil(t, res)
	})
}
