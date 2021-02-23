# State

## Parameters and base types
`Parameters` define the rules according to which votes are run. There can only
be one active parameter set at any given time. If governance wants to change a
parameter set, either to modify a value or add/remove a parameter field, a new
parameter set has to be created, and the previous one rendered inactive.

+++ https://github.com/cosmic.bet/ledger/blob/master/proto/cosmicbet/wta/v1beta1/params.proto#L10-L40

## Ticket
A single draw ticket is represented using the `Ticket` object. This contains a unique random generated id, the address of the ticket owner and the timestamp of the block in which the ticket has been created.

+++ https://github.com/cosmic.bet/ledger/blob/master/proto/cosmicbet/wta/v1beta1/models.proto#L10-L19

Tickets are created only when handling a `MsgBuyTickets` message. In order to generate a ticket id that's both unique and deterministic, the following process is used: 

```
hash = sha_256(block_hash + tx_hash + index) 
id = hex(hash[8:])
```

Each ticket is stored inside the state as 

```
TicketsStorePrefix + ticket_id | Ticket
```

## Draw
A single draw is represented inside the store using different keys. Particularly, its end time is stored using the `CurrentDrawEndTimeStoreKey` key, and the current prize is the balance of the module account having name `PrizeCollectorName`.

```
CurrentDrawEndTimeStoreKey | time.Time
```

## Historical draws
Once the winner for the current draw is extracted, the draw data and the winning ticket are both saved as a `HistoricalDrawData` object.

+++ https://github.com/cosmic.bet/ledger/blob/master/proto/cosmicbet/wta/v1beta1/models.proto#L36-L40

It is possible that a `HistoricalDrawData` does not have any winning ticket associated to it, if that draw was not entered by anyone. 

Historical draws data are stored using the following mapping: 

```
HistoricalDrawsStoreKey + Draw end time | HistoricalDrawData
```