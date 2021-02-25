<!--
order: 6
-->

# Parameters

The Winner-Takes-All module contains the following parameters:

| Key           | Type   | Example                                                                                      |
|---------------|--------|----------------------------------------------------------------------------------------------|
| DistributionParams    | object    | {"prize_percentage":"0.98","burn_percentage":"0.01","fee_percentage":"0.01"} [0]  |
| DrawParams            | object    | {"duration":"60s"} [1]                                                            |
| TicketParams          | object    | {"price":{"denom":"stake","amount":"1000000"}" [2]                                |

* [0] `prize_percentage`, `burn_percentage` `fee_percentage` must be positive, and their sum cannot exceed 1.00
* [1] `duration` must be positive and not lower than 1 minute
* [2] `amount` must be greater than 0