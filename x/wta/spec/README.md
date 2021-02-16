# `x/wta`

## Abstract 
This document specifies the WTA module of Cosmic Casino.

This module implements a "Winner-Takes-All" style of lottery, in which a single random winning ticket is extracted during each draw.

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
    - [Parameters and base types](02_state.md#parameters-and-base-types)
    - [Ticket](02_state.md#ticket)
    - [Draw](02_state.md#draw)
3. **[Messages](03_messages.md)**
    - [Buy tickets](03_messages.md#buy-tickets)
4. **[Events](04_events.md)**
    - [EndBlocker](04_events.md#beginblocker)
    - [Handlers](04_events.md#handlers)
6. **[Parameters](05_params.md)**