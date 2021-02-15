# `x/wta`

## Abstract 
This document specifies the WTA (Winner-Take-All) module of Cosmic.bet

The WTA module implements a "Winner-Take-All" style of lottery, in which a single random winning ticket is extracted during each draw.  

## Concepts

### Draws
Over time, different draws will be held. Particularly, there will be a draw every `n` blocks created. Initially, the number of blocks will be set to `17280` which represents a time spam of approximately 24 hour considering a block time of 5 seconds. The draw will be held after the `n`th block is created.

In order to take part to a draw, users will have to buy one or more tickets before the winner is drawn. Once the draw is held, the bought tickets will not be considered valid anymore, and users will have to buy new ones if the want to take part to the next draw.

For each draw, there will only be a single winning ticket drawn. The prize of each draw will be 98% of the earnings from ticket sales. 

The remaining 2% of income generated thought the sale of tickets will be used as follows: 
- 1% sent to the community pool, to fund for future developments
- 1% burnt

### Tickets
Each ticket will be represented by an object containing 
- a unique id 
- the creation time 
- the owner of the ticket 

In order to obtain a ticket, a user will have to pay using the chain token `FCHS`. A single ticket will have an initial cost of `10 FCHS`.

A single user is allowed to buy as many tickets as they can afford, no limitations set.