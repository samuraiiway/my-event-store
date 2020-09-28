# my-event-store

## Features
- Send Domain Event
- Create Aggregation Schema
- Stream listener for aggregated event

## Example
### Create Aggregation Schema
POST http://localhost:8000/aggregation/transaction
```json
{
    "aggregated_id": "payer_status_summary",
    "group_by_key_id": ["payer_id", "status"],
    "aggregated_function": [
        {
            "property_name": "amount",
            "field_name": "sum_amount",
            "function": "sum"
        },
        {
            "property_name": "amount",
            "field_name": "min_amount",
            "function": "min"
        },
        {
            "property_name": "amount",
            "field_name": "max_amount",
            "function": "max"
        },
        {
            "property_name": "amount",
            "field_name": "count_transactions",
            "function": "count"
        }
    ]
}
```
### Send Event
POST http://localhost:8000/event/transaction
```json
{
    "transaction_id": "20200927-0000012345",
    "service_id": "transfer",
    "payer_id": "123450",
    "payee_id": "98765",
    "amount": 100,
    "status": "success"
}
```
### Stream Listener
GET http://localhost:8000/streaming/aggregation/transaction/payer_status_summary
 ```json
data: {"group_key":"123450:success","data":{"count_transactions":2,"max_amount":100,"min_amount":0,"sum_amount":200}}

data: {"group_key":"123450:success","data":{"count_transactions":3,"max_amount":100,"min_amount":0,"sum_amount":300}}
 ```
