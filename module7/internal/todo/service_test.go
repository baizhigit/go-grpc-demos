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
	"github.com/stretchr/testify/mock"
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

func requireStatus(t *testing.T, err error, code codes.Code, msg string) {
	st, ok := status.FromError(err)
	require.True(t, ok)

	assert.Equal(t, code, st.Code())
	assert.Equal(t, msg, st.Message())
}

func TestService_AddTask(t *testing.T) {

	t.Run("returns INVALID_ARGUMENT status code when task is empty", func(t *testing.T) {
		service, store := newTestService(t)

		res, err := service.AddTask(context.Background(), &proto.AddTaskRequest{Task: ""})

		require.Error(t, err)
		require.Nil(t, res)

		store.AssertNotCalled(t, "AddTask", mock.Anything)

		requireStatus(t, err, codes.InvalidArgument, "task cannot be empty")
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

		requireStatus(t, err, codes.Internal, fmt.Sprintf("failed to add task: %v", testErr))
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

		requireStatus(t, err, codes.NotFound, "task not found")
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

		requireStatus(t, err, codes.Internal, fmt.Sprintf("failed to complete task: %v", testErr))
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

func TestService_ListTasks(t *testing.T) {
	t.Run("returns INTERNAL status code when an error is returned from store", func(t *testing.T) {
		service, todoStore := newTestService(t)

		testErr := errors.New("some error")

		todoStore.On("ListTasks").
			Return(nil, testErr).
			Once()

		res, err := service.ListTasks(context.Background(), &proto.ListTasksRequest{})
		require.Error(t, err)
		require.Nil(t, res)

		requireStatus(t, err, codes.Internal, fmt.Sprintf("failed to list tasks: %v", testErr))
	})

	t.Run("returns a list of tasks retrieved from store", func(t *testing.T) {
		service, todoStore := newTestService(t)

		tasks := []todostore.Task{
			{ID: "1", Task: "wake up"},
			{ID: "2", Task: "walk the dog"},
			{ID: "3", Task: "have breakfast"},
		}

		todoStore.On("ListTasks").
			Return(tasks, nil).
			Once()

		res, err := service.ListTasks(context.Background(), &proto.ListTasksRequest{})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.GetTasks(), 3)

		assert.Equal(t, "1", res.Tasks[0].Id)
		assert.Equal(t, "wake up", res.Tasks[0].Task)
		assert.Equal(t, "2", res.Tasks[1].Id)
		assert.Equal(t, "walk the dog", res.Tasks[1].Task)
		assert.Equal(t, "3", res.Tasks[2].Id)
		assert.Equal(t, "have breakfast", res.Tasks[2].Task)
	})
}
