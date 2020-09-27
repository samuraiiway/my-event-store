package service

import (
	"fmt"

	"github.com/samuraiiway/my-event-store/model"
	"github.com/samuraiiway/my-event-store/repository"
)

var (
	// <namespace, <aggregated_id, task_struct>>
	_ROOT_TASK = map[string]map[string]*AggregatedTask{}
)

func GetRootTasks() map[string]map[string]*AggregatedTask {
	return _ROOT_TASK
}

func NewAggregatedTask(namespace string, request *model.AggregatedCreateRequest) {
	task := AggregatedTask{
		ID:           request.AggregatedID,
		GroupByKeyID: request.GroupByKeyID,
	}

	for _, f := range request.AggregatedFunction {
		aggFunc := AggregatedFunction{
			PropertyName: f.PropertyName,
			FieldName:    f.FieldName,
		}

		if function := getFunctionImpl(f.Function); function != nil {
			aggFunc.Function = function
			task.Functions = append(task.Functions, aggFunc)
		}
	}

	registerTask(namespace, task)
}

func getFunctionImpl(funcName string) Function {
	if funcName == "sum" {
		return SumFunction{}
	} else if funcName == "min" {
		return MinFunction{}
	} else if funcName == "max" {
		return MaxFunction{}
	} else if funcName == "count" {
		return CountFunction{}
	} else if funcName == "last" {
		return LastFunction{}
	}

	return nil
}

func registerTask(namespace string, task AggregatedTask) {
	namespaceTask, ok := _ROOT_TASK[namespace]

	if !ok {
		namespaceTask = map[string]*AggregatedTask{}
		_ROOT_TASK[namespace] = namespaceTask
	}

	namespaceTask[task.ID] = &task
}

func DoTask(namespace string, data map[string]interface{}) {
	if namespaceTask, ok := _ROOT_TASK[namespace]; ok {
		for aggID, aggTask := range namespaceTask {
			groupKey := getGroupKey(aggTask.GroupByKeyID, data)

			valueMap := repository.GetValueMap(namespace, aggID, groupKey)

			for _, task := range aggTask.Functions {
				value := task.Function.Apply(valueMap[task.FieldName], data[task.PropertyName])
				valueMap[task.FieldName] = value
			}

			fmt.Printf("%v: %v: %v\n", aggID, groupKey, valueMap)
			NotifyListener(namespace, aggID, groupKey, valueMap)
		}
	}
}

func getGroupKey(groupKey []string, data map[string]interface{}) string {
	keys := ""

	for i, key := range groupKey {
		if i == len(groupKey)-1 {
			keys += fmt.Sprintf("%v", data[key])
		} else {
			keys += fmt.Sprintf("%v:", data[key])
		}
	}

	return keys
}
