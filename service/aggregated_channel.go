package service

import (
	"fmt"

	"github.com/samuraiiway/my-event-store/model"
)

var (
	// <namespace, channel<data>>
	_ROOT_AGGREGATED_CHANNEL = map[string]chan map[string]interface{}{}
	// <namespace, <aggregated_id, <client_id, channel<aggregated_stream_response>>>>
	_ROOT_AGGREGATED_LISTENER = map[string]map[string]map[string]chan model.AggregatedStreamResponse{}
)

func GetRootAggregatedChannel() map[string]chan map[string]interface{} {
	return _ROOT_AGGREGATED_CHANNEL
}

func GetRootAggregatedListener() map[string]map[string]map[string]chan model.AggregatedStreamResponse {
	return _ROOT_AGGREGATED_LISTENER
}

func NewAggregatedChannel(namespace string) {
	aggChan, ok := _ROOT_AGGREGATED_CHANNEL[namespace]

	if !ok {
		aggChan = make(chan map[string]interface{}, 1000)
		_ROOT_AGGREGATED_CHANNEL[namespace] = aggChan

		go func(aggChan chan map[string]interface{}, namespace string) {
			for {
				data, ok := <-aggChan

				if !ok {
					return
				}

				DoTask(namespace, data)
			}
		}(aggChan, namespace)
	}
}

func SendTask(namespace string, data map[string]interface{}) {
	aggChan, ok := _ROOT_AGGREGATED_CHANNEL[namespace]

	if ok {
		aggChan <- data
	}
}

func getNamespaceListener(namespace string) map[string]map[string]chan model.AggregatedStreamResponse {
	namespaceListener, ok := _ROOT_AGGREGATED_LISTENER[namespace]

	if !ok {
		namespaceListener = map[string]map[string]chan model.AggregatedStreamResponse{}
		_ROOT_AGGREGATED_LISTENER[namespace] = namespaceListener
	}

	return namespaceListener
}

func getAggregatedListener(namespace string, aggID string) map[string]chan model.AggregatedStreamResponse {
	namespaceListener := getNamespaceListener(namespace)
	aggregatedListener, ok := namespaceListener[aggID]

	if !ok {
		aggregatedListener = map[string]chan model.AggregatedStreamResponse{}
		namespaceListener[aggID] = aggregatedListener
	}

	return aggregatedListener
}

func RegisterClientListener(namespace string, aggID string, clientID string) chan model.AggregatedStreamResponse {
	aggregatedListener := getAggregatedListener(namespace, aggID)
	channel, ok := aggregatedListener[clientID]
	if !ok {
		channel = make(chan model.AggregatedStreamResponse, 1000)
		aggregatedListener[clientID] = channel
	}
	return channel
}

func DeregisterClientListener(namespace string, aggID string, clientID string) {
	aggregatedListener := getAggregatedListener(namespace, aggID)
	channel, _ := aggregatedListener[clientID]

	if channel != nil {
		close(channel)
		delete(aggregatedListener, clientID)
	}
}

func NotifyListener(namespace string, aggID string, groupKey string, data map[string]interface{}) {
	aggregatedListener := getAggregatedListener(namespace, aggID)
	response := model.AggregatedStreamResponse{
		GroupKey: groupKey,
		Data:     data,
	}
	for _, channel := range aggregatedListener {
		fmt.Println("notify listener")
		channel <- response
	}
}
