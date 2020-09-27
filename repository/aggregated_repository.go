package repository

var (
	// <namespace, <aggregated_id, <group_by_key_id, <field_name, value>>>>
	_ROOT_AGGREGATED_VALUE = map[string]map[string]map[string]map[string]interface{}{}
)

func getAggregatedMap(namespace string) map[string]map[string]map[string]interface{} {
	aggMap, ok := _ROOT_AGGREGATED_VALUE[namespace]

	if !ok {
		aggMap = map[string]map[string]map[string]interface{}{}
		_ROOT_AGGREGATED_VALUE[namespace] = aggMap
	}

	return aggMap
}

func getGroupMap(namespace string, aggID string) map[string]map[string]interface{} {
	aggMap := getAggregatedMap(namespace)

	groupMap, ok := aggMap[aggID]

	if !ok {
		groupMap = map[string]map[string]interface{}{}
		aggMap[aggID] = groupMap
	}

	return groupMap
}

func GetValueMap(namespace string, aggID string, group_key string) map[string]interface{} {
	groupMap := getGroupMap(namespace, aggID)

	valueMap, ok := groupMap[group_key]

	if !ok {
		valueMap = map[string]interface{}{}
		groupMap[group_key] = valueMap
	}

	return valueMap
}
