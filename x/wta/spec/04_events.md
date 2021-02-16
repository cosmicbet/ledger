# Events
The Winner-Takes-All module emits the following events:

## BeginBlocker

| Type              | Attribute Key   | Attribute Value  |
| ----------------- | --------------- | ---------------- |
| winner_drawn [0]  | winner_address  | {WinnerAddress}             |
| winner_drawn [0]  | won_amount      | {WonAmount}                 |
| new_draw     [1]  | draw_closing    | {NewDrawClosingTimestamp}   |

- [0] Event only emitted when a winner is drawn
- [1] Event only emitted when the current draw is closed 

## Handlers

### MsgBuyTickets

| Type                | Attribute Key       | Attribute Value |
| ------------------- | ------------------- | --------------- |
| buy_ticket [0]      | ticket_id           | {TicketID}            |
| buy_ticket [0]      | ticket_buyer        | {BuyerAddress}        |
| buy_ticket [0]      | ticket_timestamp    | {PurchaseTimestamp}   |
| prize_increase      | prize_amount        | {TotalPrizeAmount}    |
| message             | module              | wat                   |
| message             | action              | buy_tickets           |
| message             | sender              | {senderAddress}       |

- [0] Event emitted for each ticket bought