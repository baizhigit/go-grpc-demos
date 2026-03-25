package todo

import (
	"context"
	"sync"

	"github.com/baizhigit/go-grpc-demos/module2/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	proto.UnimplementedTodoServiceServer
	tasks map[string]string
	mu    sync.RWMutex
}

func NewService() *service {
	return &service{
		tasks: make(map[string]string),
	}
}

func (s *service) AddTask(ctx context.Context, request *proto.AddTaskRequest) (*proto.AddTaskResponse, error) {
	// validate input
	if request.GetTask() == "" {
		return nil, status.Error(codes.InvalidArgument, "task cannot be empty")
	}

	// generate ID for task
	id := uuid.New().String()

	// add task to store
	s.mu.Lock()
	s.tasks[id] = request.GetTask()
	s.mu.Unlock()

	// return generated ID
	return &proto.AddTaskResponse{
		Id: id,
	}, nil
}

func (s *service) CompleteTask(ctx context.Context, request *proto.CompleteTaskRequest) (*proto.CompleteTaskResponse, error) {
	// check if task exists
	s.mu.RLock()
	if _, ok := s.tasks[request.GetId()]; !ok {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	s.mu.RUnlock()

	// remove task from store
	s.mu.Lock()
	delete(s.tasks, request.GetId())
	s.mu.Unlock()

	// return response
	return &proto.CompleteTaskResponse{}, nil
}

func (s *service) ListTasks(ctx context.Context, request *proto.ListTasksRequest) (*proto.ListTasksResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// initialise a slice of tasks
	tasks := make([]*proto.Task, 0, len(s.tasks))

	// iterate through tasks in our store
	for id, task := range s.tasks {
		tasks = append(tasks, &proto.Task{
			Id:   id,
			Task: task,
		})
	}

	// return list of tasks
	return &proto.ListTasksResponse{
		Tasks: tasks,
	}, nil
}
